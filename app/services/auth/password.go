package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// password: Password123!
	DefaultHash string = "$argon2id$v=19$m=15000,t=2,p=2$ZW9QUTMwSlc5UmVrVTJtQQ$kRLJ4I+I0a/qLndBwL11ZKL+Vfw0nOpeoN09h0EbrS8"
)

var (
	ErrInvalidHashFormat   error = fmt.Errorf("hash doesn't follow PHC format")
	ErrIncompatibleVersion error = fmt.Errorf("incompatible algorithm version")
)

type ArgonConfig struct {
	SaltLength  uint32
	Iterations  uint32
	Memory      uint32
	KeyLength   uint32
	Parallelism uint8
}

func (c *ArgonConfig) formatPHC(hash, salt []byte) string {

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, c.Memory, c.Iterations, c.Parallelism, b64Salt, b64Hash)
}

var DefaultConf ArgonConfig = ArgonConfig{
	SaltLength:  16,
	Iterations:  3,
	Memory:      64 * 1024,
	KeyLength:   32,
	Parallelism: 2,
}

func generateSalt(size uint32) ([]byte, error) {
	buff := make([]byte, size)
	_, err := rand.Read(buff)
	return buff, err
}

func HashPassword(password string, c *ArgonConfig) (encodedHash string, err error) {
	salt, err := generateSalt(c.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	return c.formatPHC(hash, salt), nil
}

func ValidatePassword(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	c, salt, hash, err := getParametersFromHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func getParametersFromHash(encodedHash string) (c *ArgonConfig, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHashFormat
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	c = &ArgonConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Iterations, &c.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	c.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	c.KeyLength = uint32(len(hash))

	return c, salt, hash, nil
}
