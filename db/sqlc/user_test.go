package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	util2 "github.com/yilinyo/project_bank/util"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util2.HashPassword(util2.RandomStr(6))
	require.NoError(t, err)
	if err != nil {
		return User{}
	}
	arg := CreateUserParams{
		Username:       util2.RandomStr(8),
		HashedPassword: hashPassword,
		FullName:       util2.RandomStr(5),
		Email:          util2.RandomEmail(8),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {

	createRandomUser(t)

}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	newFullName := util2.RandomStr(10)
	arg := UpdateUserParams{
		Username: user.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	}
	u, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user.Username, u.Username)
	require.Equal(t, newFullName, u.FullName)
	require.Equal(t, user.Email, u.Email)
	//require.WithinDuration(t, user.CreatedAt, u.CreatedAt, time.Second)
}
