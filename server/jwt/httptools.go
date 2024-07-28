package token

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetAPIToken(r *http.Request) string {
	return r.Header.Get("Api-Key")
}

func GetValidClaims(r *http.Request, validator *ValidationService) *UserClaims {

	tokenString := getBearerToken(r)
	if tokenString == "" {
		log.Error().Msg("No bearer token send with request.")
		return nil
	}
	token, err := validator.ValidateAccessToken(tokenString)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate access token for user.")
		return nil
	}
	return token
}

// Gets the bearer token from the Authorization header field of the given request. If there is no bearer token, GetBearerToken() returns "".
func getBearerToken(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}
