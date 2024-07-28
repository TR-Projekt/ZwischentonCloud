package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"

	"github.com/TR-Projekt/zwischentoncloud/server/database"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
)

func IsAuthorizedToUseHandler(claims *token.UserClaims, userObjectIDs []int, r *http.Request) error {
	if claims.UserRole != token.ADMIN {
		objectID, err := ObjectID(r)
		if err != nil {
			return err
		}
		if !slices.Contains(userObjectIDs, objectID) {
			return errors.New("user is not authorized to use handler")
		}
	}
	return nil
}

func RegisterZwischentonForUser(userID string, zwischentonID string, db *sql.DB) error {
	return registerEntityForUser(userID, database.Zwischenton, zwischentonID, db)
}

func registerEntityForUser(userID string, entity database.Entity, entityID string, db *sql.DB) error {

	_, err := database.SetEntityForUser(entity, db, entityID, userID)
	if err != nil {
		return err
	}
	return nil
}
