package interactor

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type IPasswordInteractor interface {
	Compare(data string, hash string) bool
	Hash(data string, saltSize int) string
}

type PasswordInteractor struct {
}

func (p *PasswordInteractor) genSha256(data string) string {
	var hasher = sha256.New()
	hasher.Write([]byte(data))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (p *PasswordInteractor) Compare(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p.genSha256(password)))
	return err == nil
}

func (p *PasswordInteractor) Hash(password string, saltSize int) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.genSha256(password)), saltSize)
	if err != nil {
		return ""
	}

	return string(hashedPassword)
}

func NewPasswordInteractor() IPasswordInteractor {
	return &PasswordInteractor{}
}
