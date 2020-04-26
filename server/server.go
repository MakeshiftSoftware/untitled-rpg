package server

import (
	"net/http"
	"strconv"
	"time"
	"untitled_rpg/logger"
	"untitled_rpg/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server represents the http server that handles requests.
//
type Server struct {
	logger         logger.Logger           // logger provides logging.
	srv            *http.Server            // srv is the underlying http server.
	accountService *service.AccountService // accountService is the service that handles account related requests.
	authService    *service.AuthService    // authService is the service that handles login requests.
}

// NewServer initializes and returns a new server.
func NewServer(logger logger.Logger, port int, accountService *service.AccountService, authService *service.AuthService) *Server {
	s := &Server{logger: logger}

	handler := s.setupServices(accountService, authService)

	// TODO: timeout values from config object as well as port
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
		s.logger.Fatal().Err(err).Send()
	}
}

// Stop stops the server.
func (s *Server) Stop() {

}

// setupServices initializes the server's http handler and then attaches core
// middleware functions and services.
func (s *Server) setupServices(services ...service.Service) http.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for _, s := range services {
		s.Register(router)
	}

	CORSHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "User-Agent"})
	CORSOrigins := handlers.AllowedOrigins([]string{"*"})
	CORSMethods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete})
	handler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router)
	handler = handlers.CORS(CORSHeaders, CORSOrigins, CORSMethods)(handler)
	handler = handlers.CompressHandler(handler)
	return handler
}
