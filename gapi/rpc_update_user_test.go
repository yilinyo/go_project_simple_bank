package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yilinyo/project_bank/db/mock"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/pb"
	"github.com/yilinyo/project_bank/token"
	util2 "github.com/yilinyo/project_bank/util"
	"google.golang.org/grpc/metadata"
)

//type eqCreateUserTxParamsMatcher struct {
//	arg      db.CreateUserTxParams
//	password string
//	user     db.User
//}
//
//func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
//	arg, ok := x.(db.CreateUserTxParams)
//	if !ok {
//		return false
//	}
//	if !util2.CheckPassword(e.password, arg.HashedPassword) {
//		return false
//	}
//
//	e.arg.HashedPassword = arg.HashedPassword
//	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
//		return false
//	}
//	err := arg.AfterCreate(e.user)
//	if err != nil {
//		return false
//	}
//	return true
//
//}
//
//func (e eqCreateUserTxParamsMatcher) String() string {
//	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
//}
//
//func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
//	return eqCreateUserTxParamsMatcher{arg, password, user}
//}

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	//hashPassword, err := util2.HashPassword(password)
	//require.NoError(t, err)
	newName := util2.RandomStr(5)
	newEmail := util2.RandomEmail(8)

	testCases := []struct {
		name        string
		body        *pb.UpdateUserRequest
		buildStub   func(store *mockdb.MockStore)
		setupAuth   func(t *testing.T, tokenMaker token.Maker) context.Context
		checkStatus func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "ok",
			body: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				ctx := context.Background()
				accessToken, _, err := tokenMaker.CreateToken(user.Username, "admin", time.Hour, token.TokenTypeAccessToken)
				bearerToken := fmt.Sprintf("%s %s", authorizationType, accessToken)
				require.NoError(t, err)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}
				return metadata.NewIncomingContext(ctx, md)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}
				//满足调用方法 传入参数匹配次数匹配则会返回一个mock值
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.User{
						Username:          user.Username,
						HashedPassword:    user.HashedPassword,
						FullName:          newName,
						Email:             newEmail,
						CreatedAt:         user.CreatedAt,
						IsEmailVerified:   user.IsEmailVerified,
						PasswordChangedAt: user.PasswordChangedAt,
					}, nil)

			},
			checkStatus: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updateUser := res.GetUser()
				require.Equal(t, user.Username, updateUser.GetUsername())
				require.Equal(t, newName, updateUser.GetFullName())
				require.Equal(t, newEmail, updateUser.GetEmail())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			mockStoreCtrl := gomock.NewController(t)
			defer mockStoreCtrl.Finish()
			store := mockdb.NewMockStore(mockStoreCtrl)
			server := newTestServer(t, store, nil)
			tc.buildStub(store)
			ctx := tc.setupAuth(t, server.tokenMaker)
			res, e := server.UpdateUser(ctx, tc.body)
			tc.checkStatus(t, res, e)

		})

	}

}
