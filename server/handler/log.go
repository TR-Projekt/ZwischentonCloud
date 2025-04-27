package handler

import (
	"errors"
	"net/http"
	"os"

	servertools "github.com/Festivals-App/festivals-server-tools"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/rs/zerolog/log"
)

func GetLog(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server info.")
		servertools.UnauthorizedResponse(w)
		return
	}
	l, err := Log("/var/log/festivals-server/info.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get info log.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func GetTraceLog(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server info.")
		servertools.UnauthorizedResponse(w)
		return
	}
	l, err := Log("/var/log/festivals-server/trace.log.")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trace log")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func Log(location string) (string, error) {

	l, err := os.ReadFile(location)
	if err != nil {
		return "", errors.New("Failed to read log file at: '" + location + "' with error: " + err.Error())
	}
	return string(l), nil
}
