package ws

import (
	"strings"

	"github.com/pkg/errors"
)

type SocketCommand string

const (
	SocketCommandUnknown  SocketCommand = "Unknown"
	SocketCommandPlaceBid SocketCommand = "PlaceBid"
)

var socketCommandMapping = map[string]SocketCommand{
	strings.ToLower(string(SocketCommandUnknown)):  SocketCommandUnknown,
	strings.ToLower(string(SocketCommandPlaceBid)): SocketCommandPlaceBid,
}

func (sc SocketCommand) MarshalText() ([]byte, error) {
	return []byte(sc), nil
}

func (sc *SocketCommand) UnmarshalText(raw []byte) error {
	level, err := SocketCommandFromString(string(raw))
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal permission level")
	}
	*sc = level
	return nil
}

func SocketCommandFromString(text string) (SocketCommand, error) {
	loweredText := strings.ToLower(text)
	if pl, ok := socketCommandMapping[loweredText]; ok {
		return pl, nil
	}

	return SocketCommand(""), errors.Errorf("%s is not a valid SocketCommand", text)
}
