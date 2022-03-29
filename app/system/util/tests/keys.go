package tests

import (
	"fmt"
	"os"
)

const (
	TEST_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwRlWyS34x0jmmC4ojIi/77jjf/Y0
+WZw36Kl+D8dGEbgxzASO9giv3XvUuEedE34Nj5nTor426R/QIDVH5HkYw==
-----END PUBLIC KEY-----`

	TEST_PRIVATE_KEY = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHonss0eFux7aVl0ABXre45wkJ5SvRBbJVlVE1E70htaoAoGCCqGSM49
AwEHoUQDQgAEwRlWyS34x0jmmC4ojIi/77jjf/Y0+WZw36Kl+D8dGEbgxzASO9gi
v3XvUuEedE34Nj5nTor426R/QIDVH5HkYw==
-----END EC PRIVATE KEY-----`
)

func LoadTestKeys() {
	os.Setenv("PRIVATE_KEY", TEST_PRIVATE_KEY)
	os.Setenv("PUBLIC_KEY", TEST_PRIVATE_KEY)

	fmt.Print(os.Getenv("PRIVATE_KEY"))
} 