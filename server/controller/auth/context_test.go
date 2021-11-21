package auth

import (
	"context"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/stretchr/testify/require"
)

func TestWithUsernameStoresString(t *testing.T) {
	expectedUsername := "user"
	ctx := context.Background()
	ctx = WithUsername(ctx, expectedUsername)
	username := ctx.Value(usernameKey{})
	require.IsType(t, "", username)
	require.EqualValues(t, expectedUsername, username)
}

func TestExtractUsernameGetsStoredString(t *testing.T) {
	expectedUsername := "user"
	ctx := context.Background()
	ctx = WithUsername(ctx, expectedUsername)
	username := ExtractUsername(ctx)
	require.EqualValues(t, expectedUsername, username)
}

func TestWithPermissionStoresPermissionLevel(t *testing.T) {
	expectedPermission := model.PermissionLevelBidder
	ctx := context.Background()
	ctx = WithPermission(ctx, expectedPermission)
	permission := ctx.Value(permissionKey{})
	require.IsType(t, model.PermissionLevel(""), permission)
	require.EqualValues(t, expectedPermission, permission)
}

func TestExtractPermissionGetsStoredString(t *testing.T) {
	expectedPermission := model.PermissionLevelBidder
	ctx := context.Background()
	ctx = WithPermission(ctx, expectedPermission)
	permission := ExtractPermission(ctx)
	require.EqualValues(t, expectedPermission, permission)
}
