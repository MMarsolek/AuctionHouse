package auth

import (
	"crypto/rand"
	"time"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/pkg/errors"
)

// VerificationResults contains specific payload data from the JWT if the token is valid.
type VerificationResults struct {
	Username   string
	Permission model.PermissionLevel
}

type userTokenPayload struct {
	jwt.Payload
	Username   string                `json:"username"`
	Permission model.PermissionLevel `json:"permission"`
}

var secretHash jwt.Algorithm

func init() {
	secret := make([]byte, 16)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}

	secretHash = jwt.NewHS256(secret)
}

// NewToken creates a new token for the specified user.
func NewToken(user *model.User) ([]byte, error) {
	now := time.Now().UTC()
	payload := userTokenPayload{
		Payload: jwt.Payload{
			Issuer:         "AuctionHouse",
			Subject:        "user",
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour * 8)),
			IssuedAt:       jwt.NumericDate(now),
		},
		Username:   user.Username,
		Permission: user.Permission,
	}

	token, err := jwt.Sign(payload, secretHash)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign token")
	}

	return token, nil
}

// VerifyToken validates the token. If the token is not valid then an error is returned.
func VerifyToken(token []byte) (*VerificationResults, error) {
	var payload userTokenPayload
	_, err := jwt.Verify(token, secretHash, &payload)
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify token")
	}

	return &VerificationResults{
		Username:   payload.Username,
		Permission: payload.Permission,
	}, nil
}
