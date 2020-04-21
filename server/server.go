package server

import (
	"net/http"
	"strconv"
	"time"
	"untitled_rpg/service"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Server is the server.
type Server struct {
	srv            *http.Server            // srv is the underlying http server.
	accountService *service.AccountService // accountService is the service that handles account related requests.
	authService    *service.AuthService    // authService is the service that handles login requests.
}

// NewServer creates a new server.
func NewServer(port int, accountService *service.AccountService, authService *service.AuthService) *Server {
	s := &Server{}

	handler := s.setupServices(accountService, authService)

	s.srv = &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           handler,
	}

	return s
}

// Start starts the server.
func (s *Server) Start() {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Send()
	}
}

// Stop stops the server.
func (s *Server) Stop() {

}

// setupServices initializes services by attaching each service to the
// root application http handler.
func (s *Server) setupServices(services ...service.Service) http.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for _, s := range services {
		s.Register(router)
	}
	return router
}
