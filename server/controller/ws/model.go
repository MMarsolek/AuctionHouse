package ws

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type commandMessage struct {
	Command SocketCommand `json:"command"`
	Payload interface{}   `json:"-"`

	RawPayload json.RawMessage `json:"payload"`
}

func (cm *commandMessage) UnmarshalText(raw []byte) error {
	type commandMessageProxy commandMessage
	var proxy commandMessageProxy
	err := json.Unmarshal(raw, &proxy)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal commandMessage")
	}

	cm.Command = proxy.Command
	if proxy.Command == SocketCommandPlaceBid {
		var payload commandMessagePlaceBid
		err = json.Unmarshal(proxy.RawPayload, &payload)
		if err != nil {
			return errors.Wrapf(err, "unable to unmarshal payload to '%s'", proxy.Command)
		}

		cm.Payload = &payload
	} else {
		return errors.Errorf("unrecognized command '%s'", proxy.Command)
	}

	return nil
}

func (cm commandMessage) MarshalText() ([]byte, error) {
	type commandMessageProxy commandMessage
	var (
		proxy commandMessageProxy
		err   error
	)

	proxy.Command = cm.Command
	proxy.RawPayload, err = json.Marshal(cm.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal payload")
	}

	return json.Marshal(proxy)
}

type commandMessagePlaceBid struct {
	ItemName  string `json:"itemName"`
	BidAmount int    `json:"bidAmount"`
}

type responseMessage struct {
	StatusCode int           `json:"statusCode"`
	Command    SocketCommand `json:"command,omitempty"`
	Message    string        `json:"message,omitempty"`
	Data       interface{}   `json:"data,omitempty"`
}
