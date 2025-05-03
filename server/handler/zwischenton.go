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

func GetZwischentons(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	zwischentons, err := GetObjects(enc.ZwischentonDB, "zwischenton", nil, r.URL.Query())
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch zwischentons")
		servertools.RespondError(w, http.StatusBadRequest, "failed to fetch zwischentons")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, zwischentons)
}

func GetZwischenton(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	zwischentons, err := GetObject(enc.ZwischentonDB, r, "zwischenton")
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to fetch zwischenton")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, zwischentons)
}

func GetZwischentonSituation(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	situations, err := GetAssociation(enc.ZwischentonDB, r, "zwischenton", "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch zwischenton situation")
		servertools.RespondError(w, http.StatusBadRequest, "failed to fetch zwischenton situation")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, situations)
}

func SetSituationForZwischenton(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToAssociateEntities(claims, claims.UserZwischentons, claims.UserSituations, r)
	if err != nil {
		log.Error().Msg("User not authorized to use SetSituationForZwischenton  on the given zwischenton")
		servertools.UnauthorizedResponse(w)
		return
	}

	err = SetAssociation(enc.ZwischentonDB, r, "zwischenton", "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to set situation for zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to set situation for zwischenton")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, []any{})
}

func RemoveSituationForZwischenton(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToAssociateEntities(claims, claims.UserZwischentons, claims.UserSituations, r)
	if err != nil {
		log.Error().Msg("User not authorized to use RemoveSituationForZwischenton on the given zwischenton")
		servertools.UnauthorizedResponse(w)
		return
	}

	err = RemoveAssociation(enc.ZwischentonDB, r, "zwischenton", "situation")
	if err != nil {
		log.Error().Err(err).Msg("failed to remove situation from zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to remove situation from zwischenton")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, []any{})
}

func CreateZwischenton(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.CREATOR && claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to create a zwischenton.")
		servertools.UnauthorizedResponse(w)
		return
	}

	zwischentons, err := Create(enc.ZwischentonDB, r, "zwischenton")
	if err != nil {
		log.Error().Err(err).Msg("failed to create zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to create zwischenton")
		return
	}

	if len(zwischentons) != 1 {
		log.Error().Err(err).Msg("failed to retrieve zwischenton after creation")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = RegisterZwischentonForUser(claims.UserID, strconv.Itoa(zwischentons[0].(model.Zwischenton).ID), enc.IdentityDB)
	if err != nil {
		// try again a little bit later
		time.Sleep(2 * time.Second)
		err = RegisterZwischentonForUser(claims.UserID, strconv.Itoa(zwischentons[0].(model.Zwischenton).ID), enc.IdentityDB)
		if err != nil {
			log.Error().Err(err).Msg("failed to register zwischenton for user after creation")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}
	servertools.RespondJSON(w, http.StatusOK, zwischentons)
}

func UpdateZwischenton(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToUseHandler(claims, claims.UserZwischentons, r)
	if err != nil {
		log.Error().Msg("User not authorized to use UpdateZwischenton on the given zwischenton")
		servertools.UnauthorizedResponse(w)
		return
	}

	zwischentons, err := Update(enc.ZwischentonDB, r, "zwischenton")
	if err != nil {
		log.Error().Err(err).Msg("failed to update zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to update zwischenton")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, zwischentons)
}

func DeleteZwischenton(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	err := IsAuthorizedToUseHandler(claims, claims.UserZwischentons, r)
	if err != nil {
		log.Error().Msg("User not authorized to use DeleteZwischenton on the given zwischenton")
		servertools.UnauthorizedResponse(w)
		return
	}

	err = Delete(enc.ZwischentonDB, r, "zwischenton")
	if err != nil {
		log.Error().Err(err).Msg("failed to delete zwischenton")
		servertools.RespondError(w, http.StatusBadRequest, "failed to delete zwischenton")
		return
	}
	servertools.RespondJSON(w, http.StatusOK, nil)
}
