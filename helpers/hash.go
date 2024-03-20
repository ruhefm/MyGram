package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashDong(item string) string {
	salt := 8
	items := []byte(item)
	hash, _ := bcrypt.GenerateFromPassword(items, salt)
	return string(hash)
}

func CompareDong(enkripsi, password []byte) bool {
	hash, pass := []byte(enkripsi), []byte(password)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil

}
