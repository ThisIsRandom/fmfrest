package internal

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 10)

	return string(bytes), err
}

func HashCompare(pw, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
		return false, err
	}

	return true, nil
}
