package gapi

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yilinyo/project_bank/db/mock"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/pb"
	util2 "github.com/yilinyo/project_bank/util"
	"github.com/yilinyo/project_bank/worker"
	mockwk "github.com/yilinyo/project_bank/worker/mock"

	"reflect"
	"testing"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}
	if !util2.CheckPassword(e.password, arg.HashedPassword) {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}
	err := arg.AfterCreate(e.user)
	if err != nil {
		return false
	}
	return true

}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreatUserAPI(t *testing.T) {
	user, password := randomUser(t)
	hashPassword, err := util2.HashPassword(password)
	require.NoError(t, err)
	print(hashPassword)

	testCases := []struct {
		name        string
		body        *pb.CreateUserRequest
		buildStub   func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkStatus func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "ok",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStub: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				//满足调用方法 传入参数匹配次数匹配则会返回一个mock值
				store.EXPECT().
					//这里hashpassword 是在CreateUserTx方法里 生成的区 所以这里提前定义是不知道的，要自己写matcher
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).Return(nil)

			},
			checkStatus: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createUser := res.GetUser()
				require.Equal(t, user.Username, createUser.GetUsername())
				require.Equal(t, user.FullName, createUser.GetFullName())
				require.Equal(t, user.Email, createUser.GetEmail())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			mockStoreCtrl := gomock.NewController(t)
			defer mockStoreCtrl.Finish()
			store := mockdb.NewMockStore(mockStoreCtrl)
			mockDistributorCtrl := gomock.NewController(t)
			defer mockDistributorCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(mockDistributorCtrl)
			tc.buildStub(store, taskDistributor)
			server := newTestServer(t, store, taskDistributor)
			res, e := server.CreateUser(context.Background(), tc.body)
			tc.checkStatus(t, res, e)

		})

	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util2.RandomStr(8)
	hashedPassword, err := util2.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util2.RandomStr(6),
		HashedPassword: hashedPassword,
		FullName:       util2.RandomStr(8),
		Email:          util2.RandomEmail(8),
	}
	return
}

//func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
//	data, err := io.ReadAll(body)
//	require.NoError(t, err)
//
//	var gotUser db.User
//	err = json.Unmarshal(data, &gotUser)
//
//	require.NoError(t, err)
//	require.Equal(t, user.Username, gotUser.Username)
//	require.Equal(t, user.FullName, gotUser.FullName)
//	require.Equal(t, user.Email, gotUser.Email)
//	require.Empty(t, gotUser.HashedPassword)
//}
