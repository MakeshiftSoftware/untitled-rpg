package service

import "github.com/gorilla/mux"

// Service describes a collection of related http handlers.
type Service interface {
	Register(router *mux.Router)
}
