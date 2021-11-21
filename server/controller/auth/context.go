package auth

import (
	"context"

	"github.com/MMarsolek/AuctionHouse/model"
)

type (
	usernameKey   struct{}
	permissionKey struct{}
)

// WithUsername stores a username in the context.
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey{}, username)
}

// ExtractUsername retrieves the username from the context. Returns an empty string if the context doesn't have a username.
func ExtractUsername(ctx context.Context) string {
	rawUsername := ctx.Value(usernameKey{})
	if rawUsername == nil {
		return ""
	}

	result, ok := rawUsername.(string)
	if !ok {
		return ""
	}

	return result
}

// WithPermission stores a user's permission in the context.
func WithPermission(ctx context.Context, permission model.PermissionLevel) context.Context {
	return context.WithValue(ctx, permissionKey{}, permission)
}

// ExtractPermission retrieves the user's permission from the context. Returns an empty string if the context doesn't
// have a permission.
func ExtractPermission(ctx context.Context) model.PermissionLevel {
	rawUsername := ctx.Value(permissionKey{})
	if rawUsername == nil {
		return ""
	}

	result, ok := rawUsername.(model.PermissionLevel)
	if !ok {
		return ""
	}

	return result
}
