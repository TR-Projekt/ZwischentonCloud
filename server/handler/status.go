package handler

import (
	"net/http"

	servertools "github.com/Festivals-App/festivals-server-tools"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/TR-Projekt/zwischentoncloud/server/status"
	"github.com/rs/zerolog/log"
)

func GetVersion(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server version.")
		servertools.UnauthorizedResponse(w)
		return
	}
	servertools.RespondString(w, http.StatusOK, status.VersionString())
}

func GetInfo(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server info.")
		servertools.UnauthorizedResponse(w)
		return
	}
	servertools.RespondJSON(w, http.StatusOK, status.InfoString())
}

func GetHealth(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server health.")
		servertools.UnauthorizedResponse(w)
		return
	}
	servertools.RespondCode(w, status.HealthStatus())
}
