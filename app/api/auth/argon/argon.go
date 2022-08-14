package argon

func PasswordToHash(password string) (string, error) {
	return generateFromPassword(password)
}

func PasswordIsMatch(password, hash string) (bool, error) {
	return comparePasswordAndHash(password, hash)
}
