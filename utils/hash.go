package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func ComparePasswords(reqPass, dbPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(reqPass))
	return err == nil
}
