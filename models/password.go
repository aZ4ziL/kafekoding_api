package models

import "golang.org/x/crypto/bcrypt"

// EncryptionPassword is function for encrypt the password.
func EncryptionPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// DecryptionPassword is function for decrypt the password.
func DecryptionPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
