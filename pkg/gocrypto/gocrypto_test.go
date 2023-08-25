package gocrypto

import (
	"testing"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestHashAndCompare(t *testing.T) {
	for index := 0; index < 5; index++ {
		randomID, err := uuid.NewV4()
		assert.NoError(t, err, "failed to generate random uuid")

		hash, err := Hash(
			randomID.String(),
			HashParsms{
				Memory:      64 * 1024,
				Iterations:  2,
				Parallelism: 2,
				SaltLength:  16,
				KeyLentgh:   32,
			},
		)
		assert.NoError(t, err, "failed to hash random uuid")

		err = Compare(randomID.String(), hash)
		assert.NoError(t, err, "failed to compare uuid and hash")
	}
}

func TestCompare(t *testing.T) {
	type testInput struct {
		password    string
		encodedHash string
	}
	type test struct {
		name          string
		input         testInput
		expectedError error
	}
	tests := []test{
		{
			name: "should works fine",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$65536$2$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: nil,
		},
		{
			name: "should fail after receiving an invalid hash",
			input: testInput{
				password:    "",
				encodedHash: "",
			},
			expectedError: ErrInvalidEncodedHash,
		},
		{
			name: "should fail after receiving an hash with an invalid algo",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "x$19$65536$2$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrUnsupportedHashAlgorithm,
		},
		{
			name: "should fail after receiving an hash with an invalid algo version",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19x$65536$2$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrInvalidAlgorithmVersion,
		},
		{
			name: "should fail after try to parse memory",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$x$2$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrFailedToParseData,
		},
		{
			name: "should fail after try to parse iterations",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$65536$x$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrFailedToParseData,
		},
		{
			name: "should fail after try to parse parallelism",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$65536$2$x$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrFailedToParseData,
		},
		{
			name: "should fail after try to parse salt",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$65536$2$2$+$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrFailedToParseData,
		},
		{
			name: "should fail after try to parse hash",
			input: testInput{
				password:    "1e0a2ce0-0c25-4da1-9219-8424142ef62d",
				encodedHash: "argon2id$19$65536$2$2$6LtCNxqma7JMX1vvHZEByA$+",
			},
			expectedError: ErrFailedToParseData,
		},
		{
			name: "should fail when password and hash mismatch",
			input: testInput{
				password:    "one punch",
				encodedHash: "argon2id$19$65536$2$2$6LtCNxqma7JMX1vvHZEByA$EUAlsaCBOwwflrVBjlJXsjFTjhQ1XA0bSQnW6tMTC1M",
			},
			expectedError: ErrHashMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Compare(tt.input.password, tt.input.encodedHash)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
