package model

import (
	"strings"

	"github.com/pkg/errors"
)

type PermissionLevel string

const (
	PermissionLevelBidder PermissionLevel = "Bidder"
	PermissionLevelAdmin  PermissionLevel = "Admin"
)

var permissionLevelMapping = map[string]PermissionLevel{
	strings.ToLower(string(PermissionLevelBidder)): PermissionLevelBidder,
	strings.ToLower(string(PermissionLevelAdmin)):  PermissionLevelAdmin,
}

func (pl PermissionLevel) MarshalText() ([]byte, error) {
	return []byte(pl), nil
}

func (pl *PermissionLevel) UnmarshalText(raw []byte) error {
	level, err := PermissionLevelFromString(string(raw))
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal permission level")
	}
	*pl = level
	return nil
}

func PermissionLevelFromString(text string) (PermissionLevel, error) {
	loweredText := strings.ToLower(text)
	if pl, ok := permissionLevelMapping[loweredText]; ok {
		return pl, nil
	}

	return PermissionLevel(""), errors.Errorf("%s is not a valid PermissionLevel", text)
}
