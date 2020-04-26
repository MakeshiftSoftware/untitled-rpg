package service

import (
	"encoding/json"
	"net/http"
	"untitled_rpg/domain"
	"untitled_rpg/store"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// AccountService is a collection of account related http handlers.
type AccountService struct {
	store *store.AccountStore // store is the account store used to access and save account data.
}

// NewAccountService initializes and returns a new account service.
func NewAccountService(store *store.AccountStore) *AccountService {
	return &AccountService{
		store: store,
	}
}

// Register registers all service routes with the provided router.
func (s *AccountService) Register(router *mux.Router) {
	router.HandleFunc("/accounts", s.createAccount).Methods(http.MethodPost)
}

// createAccount is an http handler that creates a new account.
func (s *AccountService) createAccount(w http.ResponseWriter, r *http.Request) {
	var account domain.Account

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondErr(w, newBadRequestError("Invalid request body"))
		return
	}

	_, err := govalidator.ValidateStruct(account)
	if err != nil {
		respondErr(w, newValidationError(err))
		return
	}

	if err := account.NormalizeEmail(); err != nil {
		respondErr(w, newInternalServerError(err))
		return
	}

	if err := account.HashPassword(); err != nil {
		respondErr(w, newInternalServerError(err))
		return
	}

	if err := s.store.CreateAccount(account); err != nil {
		if err == store.ErrAccountExists {
			respondErr(w, newConflictError(err.Error()))
		} else {
			respondErr(w, newInternalServerError(err))
		}
	}
}
