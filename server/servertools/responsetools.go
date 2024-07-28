package servertools

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func RespondCode(w http.ResponseWriter, code int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
}

func RespondString(w http.ResponseWriter, code int, message string) {

	response := []byte(message)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {

	//TODO String comparison is not very elegant!
	if fmt.Sprint(payload) == "[]" {
		payload = []interface{}{}
	}

	resultMap := map[string]interface{}{"data": payload}
	response, err := json.Marshal(resultMap)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func RespondError(w http.ResponseWriter, code int, message string) {

	resultMap := map[string]interface{}{"error": message}
	response, err := json.Marshal(resultMap)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func UnauthorizedResponse(w http.ResponseWriter) {

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}
