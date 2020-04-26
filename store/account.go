package store

import (
	"database/sql"
	"errors"
	"untitled_rpg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

var (
	// ErrAccountNotFound is returned when no account is found.
	ErrAccountNotFound = errors.New("Account not found")
	// ErrAccountExists is returned when an account already exists with the provided email address.
	ErrAccountExists = errors.New("Account already exists")
)

// AccountStore provides functions for retrieving and saving account data.
type AccountStore struct {
	db *sqlx.DB
}

// NewAccountStore initializes and returns a new account store with the provided db handle.
func NewAccountStore(db *sqlx.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

// CreateAccount saves a new account to storage.
func (s *AccountStore) CreateAccount(account domain.Account) error {
	query := `INSERT INTO accounts (email, password) VALUES ($1, $2)`

	_, err := s.db.Exec(query, account.Email, account.Password)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			return ErrAccountExists
		}
		return err
	}

	return nil
}

// GetAccount retrieves an account from storage by email.
func (s *AccountStore) GetAccount(email string) (domain.Account, error) {
	query := `SELECT id, password FROM accounts WHERE email = $1`
	var account domain.Account

	if err := s.db.Get(&account, query, email); err != nil {
		if err == sql.ErrNoRows {
			return account, ErrAccountNotFound
		}
		return account, err
	}

	return account, nil
}
