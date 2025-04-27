package handler

import (
	"net/http"

	servertools "github.com/Festivals-App/festivals-server-tools"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/TR-Projekt/zwischentoncloud/server/status"
	"github.com/rs/zerolog/log"
)

func MakeUpdate(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to update the service.")
		servertools.UnauthorizedResponse(w)
		return
	}
	newVersion, err := servertools.RunUpdate(status.ServerVersion, "Festivals-App", "festivals-server", "/usr/local/festivals-server/update.sh")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondString(w, http.StatusAccepted, newVersion)
}
