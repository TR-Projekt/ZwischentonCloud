package handler

import (
	"net/http"
	"strconv"
	"time"

	servertools "github.com/Festivals-App/festivals-server-tools"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/TR-Projekt/zwischentoncloud/server/model"
	"github.com/rs/zerolog/log"
)

func GetSituations(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	situations, err := GetObjects(enc.ZwischentonDB, "situation", nil, r.URL.Query())
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch situations")
		servertools.RespondError(w, http.StatusBadRequest, "failed to fetch situations")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, situations)
}

func GetSituation(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	situations, err := GetObject(enc.ZwischentonDB, r, "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch situation")
		servertools.RespondError(w, http.StatusBadRequest, "failed to fetch situation")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, situations)
}

func CreateSituation(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.CREATOR && claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to create a situation.")
		servertools.UnauthorizedResponse(w)
		return
	}

	situations, err := Create(enc.ZwischentonDB, r, "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to create situation")
		servertools.RespondError(w, http.StatusBadRequest, "failed to create situation")
		return
	}

	if len(situations) != 1 {
		log.Error().Err(err).Msg("failed to retrieve situation after creation")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = RegisterSituationForUser(claims.UserID, strconv.Itoa(situations[0].(model.Situation).ID), enc.IdentityDB)
	if err != nil {
		// try again a little bit later
		time.Sleep(2 * time.Second)
		err = RegisterSituationForUser(claims.UserID, strconv.Itoa(situations[0].(model.Situation).ID), enc.IdentityDB)
		if err != nil {
			log.Error().Err(err).Msg("failed to register zwischenton for user after creation")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}

	servertools.RespondJSON(w, http.StatusOK, situations)
}

func UpdateSituation(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToUseHandler(claims, claims.UserSituations, r)
	if err != nil {
		log.Error().Msg("User not authorized to use UpdateSituation on the given situation")
		servertools.UnauthorizedResponse(w)
		return
	}

	situations, err := Update(enc.ZwischentonDB, r, "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to update situation")
		servertools.RespondError(w, http.StatusBadRequest, "failed to update situation")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, situations)
}

func DeleteSituation(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToUseHandler(claims, claims.UserSituations, r)
	if err != nil {
		log.Error().Msg("User not authorized to use DeleteSituation on the given situation")
		servertools.UnauthorizedResponse(w)
		return
	}

	err = Delete(enc.ZwischentonDB, r, "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to delete situation")
		servertools.RespondError(w, http.StatusBadRequest, "failed to delete situation")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, nil)
}
