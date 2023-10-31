package helpers

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(passowrd []byte) string {
	hashedPassword, err := bcrypt.GenerateFromPassword(passowrd, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func PasswordCompare(pasword []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, pasword)
	return err
}
