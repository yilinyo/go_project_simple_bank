package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yilinyo/project_bank/db/util"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomStr(6))
	require.NoError(t, err)
	if err != nil {
		return User{}
	}
	arg := CreateUserParams{
		Username:       util.RandomStr(8),
		HashedPassword: hashPassword,
		FullName:       util.RandomStr(5),
		Email:          util.RandomEmail(8),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
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
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}
