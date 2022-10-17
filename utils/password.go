package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

// package utils

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func HashPassword(password string) (string, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 	if err != nil {
// 		return "", fmt.Errorf("could not hash password %w", err)
// 	}
// 	return string(hashedPassword), nil
// }

// func VerifyPassword(hashed_password string, candidate_password string) error {

// 	return bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(candidate_password))
// }
