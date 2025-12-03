package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {

	password, err := HashPassword("hello world")
	require.NoError(t, err)
	fmt.Println(password)
	require.True(t, CheckPassword("hello world", password))
	require.False(t, CheckPassword("hello world!", password))

}
