package service

import "github.com/gorilla/mux"

// Service defines an interface representing a collection of related http handlers.
type Service interface {
	Register(router *mux.Router)
}
