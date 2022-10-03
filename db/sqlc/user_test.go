package db

import (
	"context"
	"simplebank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomUser -> creates an user for testing finctionality

func createRandomUser(t *testing.T) User {

	hashPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:     utils.RandomOwner(),
		HashPassword: hashPassword,
		FullName:     utils.RandomOwner(),
		Email:        utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashPassword, user.HashPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordCahngedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestCreateUser -> unit testing the creating user functionality
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestGetUser -> unit testing the get user functionality
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordCahngedAt, user2.PasswordCahngedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
