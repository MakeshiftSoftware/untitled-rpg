package service

import (
	"encoding/json"
	"net/http"
	"untitled_rpg/domain"
	"untitled_rpg/store"
	"untitled_rpg/token"

	"github.com/gorilla/mux"
)

// AuthService is a collection of authentication related http handlers.
type AuthService struct {
	store         *store.AccountStore // store is the account store used to access and save account data.
	tokenProvider *token.Provider     // tokenProvider is used to generate a new auth token following a successful login.
}

// NewAuthService initializes and returns a new auth service.
func NewAuthService(store *store.AccountStore, tokenProvider *token.Provider) *AuthService {
	return &AuthService{
		store:         store,
		tokenProvider: tokenProvider,
	}
}

// Register registers all service routes with the provided router.
func (s *AuthService) Register(router *mux.Router) {
	router.HandleFunc("/authenticate", s.authenticate).Methods(http.MethodPost)
}

// authenticate is an http handler that validates an email and password combination and
// returns an auth token that can be used to interact with the server.
func (s *AuthService) authenticate(w http.ResponseWriter, r *http.Request) {
	var checkAccount domain.Account

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&checkAccount); err != nil {
		respondErr(w, newBadRequestError("Invalid request body"))
		return
	}

	if err := checkAccount.NormalizeEmail(); err != nil {
		respondErr(w, newInternalServerError(err))
		return
	}

	account, err := s.store.GetAccount(checkAccount.Email)
	if err != nil {
		if err == store.ErrAccountNotFound {
			respondErr(w, newNotFoundError(err.Error()))
		} else {
			respondErr(w, newInternalServerError(err))
		}
		return
	}

	if match := account.CheckPassword(checkAccount.Password); !match {
		respondErr(w, newUnauthorizedError())
		return
	}

	token, err := s.tokenProvider.IssueToken(account)
	if err != nil {
		respondErr(w, newUnauthorizedError())
		return
	}

	if _, err := w.Write([]byte(token)); err != nil {
		respondErr(w, newInternalServerError(err))
	}
}
