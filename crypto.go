package main

import "golang.org/x/crypto/bcrypt"

// Encrypt password and create salt this shouldn't
// be called directly except inside Create
func generatePassword(clearPass string) string {
	// Hashing the password with the default cost of 10
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(clearPass), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// Check that the passed in password is valid
func validPassword(clearPass, cryptText string) error {
	err := bcrypt.CompareHashAndPassword([]byte(cryptText), []byte(clearPass))
	return err
}
