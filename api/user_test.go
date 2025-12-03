package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yilinyo/project_bank/db/mock"
	db "github.com/yilinyo/project_bank/db/sqlc"
	util2 "github.com/yilinyo/project_bank/util"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	if !util2.CheckPassword(e.password, arg.HashedPassword) {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreatUserAPI(t *testing.T) {
	user, password := randomUser(t)
	hashPassword, err := util2.HashPassword(password)
	require.NoError(t, err)
	print(hashPassword)

	testCases := []struct {
		name        string
		body        gin.H
		buildStub   func(store *mockdb.MockStore)
		checkStatus func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
					//HashedPassword: hashPassword,
				}
				//满足调用方法 传入参数匹配次数匹配则会返回一个mock值
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkStatus: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/users")

			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkStatus(t, recorder)

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

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
