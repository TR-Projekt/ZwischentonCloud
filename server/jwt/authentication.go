package token

// rename package to prevent collisons with package go/token

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Authentication interface lists the methods that our authentication service should implement
type Authentication interface {
	GenerateAccessToken(user *User) (string, error)
}

type AuthService struct {
	SigningKey        *rsa.PrivateKey
	ValidationKey     *rsa.PublicKey
	ValidationKeyFile string
	TokenLifetime     time.Duration
	Issuer            string
}

func NewAuthService(privatekey string, publickey string, tokenLifetime int, issuer string) *AuthService {

	var signKey *rsa.PrivateKey = nil
	signBytes, err := os.ReadFile(privatekey)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read privat auth key.")
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse private auth key.")
	}

	var verifyKey *rsa.PublicKey = nil
	verifyBytes, err := os.ReadFile(publickey)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read public auth key.")
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse public auth key.")
	}

	return &AuthService{SigningKey: signKey, ValidationKey: verifyKey, ValidationKeyFile: publickey, TokenLifetime: time.Minute * time.Duration(tokenLifetime), Issuer: issuer}
}
