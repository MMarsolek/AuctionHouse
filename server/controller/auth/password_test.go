package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateEncodedPasswordDoesNotError(t *testing.T) {
	_, err := GenerateEncodedPassword("password")
	require.NoError(t, err)
}

func TestGenerateEncodedPasswordDoesNotContainOriginalPassword(t *testing.T) {
	originalPassword := "hunter2"
	encodedPassword, err := GenerateEncodedPassword(originalPassword)
	require.NoError(t, err)
	require.NotContains(t, encodedPassword, originalPassword)
}

func TestGenerateEncodedPasswordErrorsOnEmptyPassword(t *testing.T) {
	_, err := GenerateEncodedPassword("")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPassword)

	_, err = GenerateEncodedPassword(" ")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPassword)
}

func TestComparePasswordAndHashReturnsTrueIfEqual(t *testing.T) {
	originalPassword := "hunter2"
	encodedPassword, err := GenerateEncodedPassword(originalPassword)
	require.NoError(t, err)

	correctPassword, err := ComparePasswordAndHash(originalPassword, encodedPassword)
	require.NoError(t, err)
	require.True(t, correctPassword, "the password was correct")
}

func TestComparePasswordAndHashErrorsOnEmptyPassword(t *testing.T) {
	_, err := ComparePasswordAndHash("", "stuff")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPassword)

	_, err = ComparePasswordAndHash(" ", "otherStuff")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPassword)
}

func TestDecodeHashReturnsInvalidHashWhenEncodedPasswordFormatIncorrect(t *testing.T) {
	_, err := ComparePasswordAndHash("whatever", "not the correct format")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidHash)
}

func TestDecodeHashReturnsIncompatibleVersionWhenEncodedPasswordVersionIncorrect(t *testing.T) {
	_, err := ComparePasswordAndHash("whatever", "0$1$v=8008135$1$2$3")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrIncompatibleVersion)
}
