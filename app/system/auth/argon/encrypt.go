package argon

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateFromPassword(password string) (encodedHash string, err error) {
	salt, err := generateRandomBytes(argonParams.saltLength)

	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonParams.iterations, argonParams.memory, argonParams.parallelism, argonParams.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argonParams.memory, argonParams.iterations, argonParams.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}
