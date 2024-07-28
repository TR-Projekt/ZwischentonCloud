package token

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type Validation interface {
	ValidateAccessToken(token string) (string, error)
}

type ValidationService struct {
	Key     *rsa.PublicKey
	APIKeys *[]string
}

func NewValidationService(auth *AuthService, APIKeys []string) *ValidationService {
	return &ValidationService{Key: auth.ValidationKey, APIKeys: &APIKeys}
}

// ValidateAccessToken parses and validates the given access token
// returns the custom claim present in the token payload
func (validator *ValidationService) ValidateAccessToken(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error().Msg("Unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth token")
		}
		return validator.Key, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to parse claims")
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid || claims.UserID == "" {
		return nil, errors.New("invalid token: authentication failed")
	}
	return claims, nil
}
