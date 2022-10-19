package common

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(val string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashed(hashed string, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return false, err
	}
	return true, nil
}
