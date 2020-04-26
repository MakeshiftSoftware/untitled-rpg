package token

import (
	"untitled_rpg/domain"

	"github.com/dgrijalva/jwt-go"
)

// Provider is a utility that handles issuing and verifying auth tokens.
// Most services will require a valid auth token to be able to interact with them.
type Provider struct {
	encryptionKey string // encryptionKey is the secret key used when issuing and verifying auth tokens.
}

// NewProvider initializes and returns a new token provider with the provided key.
func NewProvider(encryptionKey string) *Provider {
	return &Provider{
		encryptionKey: encryptionKey,
	}
}

// IssueToken creates an auth token to use when interacting with the service.
func (p *Provider) IssueToken(account domain.Account) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = account.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.encryptionKey))
}

// VerifyToken verifies an auth token.
func VerifyToken(token string) bool {
	return false
}
