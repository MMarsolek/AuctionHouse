package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/stretchr/testify/require"
)

func TestNewTokenReturnsSignedTokenForUser(t *testing.T) {
	expectedUsername := "hunter"
	user := &model.User{
		Username: expectedUsername,
	}

	rawToken, err := NewToken(user)
	require.NoError(t, err)

	rawPayloadBytes := bytes.Split(rawToken, []byte("."))[1]
	rawPayload, err := base64.StdEncoding.DecodeString(string(rawPayloadBytes))
	require.NoError(t, err)
	payload := make(map[string]interface{})
	require.NoError(t, json.Unmarshal(rawPayload, &payload))

	rawUsername := payload["username"]
	require.IsType(t, "", rawUsername)
	username := rawUsername.(string)
	require.EqualValues(t, expectedUsername, username)
}

func TestVerifyTokenReturnsVerifiedTokenWithPayload(t *testing.T) {
	user := &model.User{
		Username:   "hunter",
		Permission: model.PermissionLevelAdmin,
	}

	rawToken, err := NewToken(user)
	require.NoError(t, err)

	result, err := VerifyToken(rawToken)
	require.NoError(t, err)

	require.EqualValues(t, user.Username, result.Username)
	require.EqualValues(t, user.Permission, result.Permission)
}
