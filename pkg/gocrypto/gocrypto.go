package gocrypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	requiredHashParts      = 7
	supportedHashAlgorithm = "argon2id"
)

var (
	ErrRandomBytes              = errors.New("can not generate random bytes")
	ErrHashMismatch             = errors.New("hash mismatch")
	ErrFailedToParseData        = errors.New("failed to parse data")
	ErrInvalidEncodedHash       = errors.New("invalid encoded hash")
	ErrInvalidAlgorithmVersion  = errors.New("invalid algorithm version")
	ErrUnsupportedHashAlgorithm = errors.New("unsupported hash algorithm")
)

type HashParsms struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLentgh   uint32
}

// Hash Creats a new hash based using the given params
func Hash(password string, parms HashParsms) (string, error) {
	if salt, err := randomBytes(parms.SaltLength); err != nil {
		return "", err
	} else {
		hash := argon2.IDKey(
			[]byte(password),
			salt,
			parms.Iterations,
			parms.Memory,
			parms.Parallelism,
			parms.KeyLentgh,
		)

		// Base64 encode the salt and hashed password.
		b64Salt := base64.RawStdEncoding.EncodeToString(salt)
		b64Hash := base64.RawStdEncoding.EncodeToString(hash)

		encodedHash := fmt.Sprintf(
			"argon2id$%d$%d$%d$%d$%s$%s",
			argon2.Version,
			parms.Memory,
			parms.Iterations,
			parms.Parallelism,
			b64Salt,
			b64Hash,
		)
		return encodedHash, nil
	}
}

// Compare Check if the password and hash mataches. If not returns an error
func Compare(password, encodedHash string) error {
	hashParts := strings.Split(encodedHash, "$")

	if len(hashParts) == 0 || len(hashParts) != requiredHashParts {
		return ErrInvalidEncodedHash
	}

	if hashParts[0] != supportedHashAlgorithm {
		return ErrUnsupportedHashAlgorithm
	}

	if fmt.Sprint(argon2.Version) != hashParts[1] {
		return ErrInvalidAlgorithmVersion
	}

	memory, err := strconv.ParseUint(hashParts[2], 10, 32)
	if err != nil {
		return fmt.Errorf("%s. %w", err, ErrFailedToParseData)
	}

	iterations, err := strconv.ParseUint(hashParts[3], 10, 32)
	if err != nil {
		return fmt.Errorf("%s. %w", err, ErrFailedToParseData)
	}

	parallelism, err := strconv.ParseUint(hashParts[4], 10, 32)
	if err != nil {
		return fmt.Errorf("%s. %w", err, ErrFailedToParseData)
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(hashParts[5])
	if err != nil {
		return fmt.Errorf("%s. %w", err, ErrFailedToParseData)
	}

	storedHash, err := base64.RawStdEncoding.Strict().DecodeString(hashParts[6])
	if err != nil {
		return fmt.Errorf("%s. %w", err, ErrFailedToParseData)
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(iterations), uint32(memory), uint8(parallelism), uint32(len(storedHash)))

	if subtle.ConstantTimeCompare(storedHash, hash) == 0 {
		return ErrHashMismatch
	}

	return nil
}

func randomBytes(size uint32) ([]byte, error) {
	bs := make([]byte, size)
	if _, err := rand.Read(bs); err != nil {
		return nil, fmt.Errorf("%s. %w", err, ErrRandomBytes)
	}
	return bs, nil
}
