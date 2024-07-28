package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func GenerateAccessToken(user *token.User, db *sql.DB, auth *token.AuthService) (string, error) {

	userID := fmt.Sprint(user.ID)
	userRole := user.Role
	userZwischentons, err := GetEntitiesForUser(Zwischenton, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch zwischentons for user.")
		return "", errors.New("could not generate access token. please try again later")
	}

	claims := token.UserClaims{
		UserID:           userID,
		UserRole:         userRole,
		UserZwischentons: userZwischentons,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenLifetime)),
			Issuer:    auth.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(auth.SigningKey)
}

func RegenerateAccessToken(user *token.User, oldClaims *token.UserClaims, db *sql.DB, auth *token.AuthService) (string, error) {

	userID := fmt.Sprint(user.ID)
	userRole := user.Role
	userZwischentons, err := GetEntitiesForUser(Zwischenton, db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch zwischentons for user.")
		return "", errors.New("could not generate access token. please try again later")
	}

	claims := token.UserClaims{
		UserID:           userID,
		UserRole:         userRole,
		UserZwischentons: userZwischentons,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: oldClaims.ExpiresAt,
			Issuer:    auth.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(auth.SigningKey)
}
