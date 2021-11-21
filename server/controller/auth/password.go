package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = stderrors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = stderrors.New("incompatible version of argon2")
	ErrEmptyPassword       = stderrors.New("password is empty")
)

const saltLength = 16

// GenerateEncodedPassword uses the argon2 algorithm to generate a password.
func GenerateEncodedPassword(clearText string) (string, error) {
	if strings.TrimSpace(clearText) == "" {
		return "", errors.Wrap(ErrEmptyPassword, "supplied password is empty")
	}
	var (
		iterations  uint32 = 3
		memory      uint32 = 64 * 1024
		parallelism uint8  = 2
		keyLength   uint32 = 32
	)

	hash, salt, err := hashPassword(clearText, iterations, memory, parallelism, keyLength)
	if err != nil {
		return "", errors.Wrap(err, "unable to generate password hash")
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// ComparePasswordAndHash determines if the clearPassword can generate an identical hash to the encodedHash.
func ComparePasswordAndHash(clearPassword string, encodedHash string) (bool, error) {
	if strings.TrimSpace(clearPassword) == "" {
		return false, errors.Wrap(ErrEmptyPassword, "supplied password is empty")
	}

	hash, salt, memory, iterations, parallelism, keyLength, err := decodeHash(encodedHash)
	if err != nil {
		return false, errors.Wrap(err, "unable to decode hash")
	}

	otherHash := argon2.IDKey([]byte(clearPassword), salt, iterations, memory, parallelism, keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func hashPassword(clearText string, iterations uint32, memory uint32, parallelism uint8, keyLength uint32) ([]byte, []byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to generate random salt")
	}

	hash := argon2.IDKey([]byte(clearText), salt, iterations, memory, parallelism, keyLength)
	return hash, salt, nil
}

func decodeHash(encodedHash string) ([]byte, []byte, uint32, uint32, uint8, uint32, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, 0, 0, 0, 0, errors.Wrap(ErrInvalidHash, "encoded hash does not have the expected amount of parts")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, 0, 0, 0, 0, errors.Wrap(err, "could not read in argon2 version")
	}
	if version != argon2.Version {
		return nil, nil, 0, 0, 0, 0, errors.Wrapf(ErrIncompatibleVersion, "unexpected version '%v'", argon2.Version)
	}

	var (
		iterations  uint32
		memory      uint32
		parallelism uint8
	)
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return nil, nil, 0, 0, 0, 0, errors.Wrap(err, "unable to read in encoded hash properties")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, 0, 0, 0, 0, errors.Wrap(err, "unable to decode salt")
	}
	decodedSaltLength := uint32(len(salt))
	if decodedSaltLength != saltLength {
		return nil, nil, 0, 0, 0, 0, errors.Wrapf(ErrInvalidHash, "salt has length of %d when %d is expected", decodedSaltLength, saltLength)
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, 0, 0, 0, 0, errors.Wrap(err, "unable to decode hash")
	}
	keyLength := uint32(len(hash))

	return hash, salt, memory, iterations, parallelism, keyLength, nil
}
