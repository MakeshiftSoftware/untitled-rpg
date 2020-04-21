package service

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func respondErr(w http.ResponseWriter, err *httpError) {
	if err.code == http.StatusInternalServerError {
		log.Error().Err(err.internalError).Send()
	}
	http.Error(w, err.message, err.code)
}
