package domain

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

// Account represents a user account.
type Account struct {
	Meta
	Email    string `json:"email,omitempty" db:"email" valid:"email"`
	Password string `json:"password,omitempty" db:"password" valid:"matches(12345)"`
}

// MarshalJSON is a custom Account json marshaller that omits the password field.
func (account Account) MarshalJSON() ([]byte, error) {
	type AccountAlias Account
	return json.Marshal(&struct {
		AccountAlias
		Password string `json:"password,omitempty"`
	}{
		AccountAlias: (AccountAlias)(account),
	})
}

// HashPassword hashes the account password.
func (account *Account) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(account.Password), 14)
	if err != nil {
		return err
	}
	account.Password = string(bytes)
	return nil
}

// CheckPassword validates a password against an account's hashed password.
func (account Account) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	return err == nil
}

// NormalizeEmail normalizes the account email.
func (account *Account) NormalizeEmail() error {
	normalizedEmail, err := govalidator.NormalizeEmail(account.Email)
	if err != nil {
		return err
	}
	account.Email = normalizedEmail
	return nil
}
