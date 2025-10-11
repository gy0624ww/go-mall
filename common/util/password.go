package util

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func BcryptPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 11) // 第二个参数是cost， 值越大，加密越慢，安全性越高
	return string(bytes), err
}

func BcryptCompare(passwordHash, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
	return err == nil
}

func BcryptGetCost(passwordHash string) int {
	cost, _ := bcrypt.Cost([]byte(passwordHash))
	return cost
}

func PasswordComplexityVerify(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
