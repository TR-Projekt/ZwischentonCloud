package handler

import (
	"encoding/json"
	"io"
	"net/http"

	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/TR-Projekt/zwischentoncloud/server/database"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func Signup(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var signupVars map[string]string
	err = json.Unmarshal(body, &signupVars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	email := signupVars["email"]
	password := signupVars["password"]

	if validEmail(email) && validPassword(password) {

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate password hash from provided password.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = database.CreateUserWithEmailAndPasswordHash(enc.IdentityDB, email, string(passwordHash))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user with given email and password.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		servertools.RespondCode(w, http.StatusCreated)
		return
	}

	// If the Authentication header is not present, is invalid, or the username or password is wrong
	servertools.UnauthorizedResponse(w)
}

func Login(enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	// Extract the username and password from the request
	// Authorization header. If no Authentication header is present
	// or the header value is invalid, then the 'ok' return value
	// will be false.
	email, password, ok := r.BasicAuth()

	if ok {

		// retrieve user for the given username
		requestedUser, err := database.GetUserByEmail(enc.IdentityDB, email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch user.")
			// do i need to mitigate timing attacks on email guessing?
			servertools.UnauthorizedResponse(w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(requestedUser.PasswordHash), []byte(password))
		// If the password is correct return the authentication jwt token
		if err == nil {
			token, err := database.GenerateAccessToken(requestedUser, enc.IdentityDB, enc.Auth)
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate access token for user.")
				servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			servertools.RespondString(w, http.StatusOK, token)
			return
		} else {
			log.Error().Err(err).Msg("The password provided was wrong.")
		}
	}

	// If the Authentication header is not present, is invalid, or the username or password is wrong
	servertools.UnauthorizedResponse(w)
}

func Refresh(claims *token.UserClaims, enc *HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request) {

	requestedUser, err := database.GetUserByID(enc.IdentityDB, claims.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user.")
		servertools.UnauthorizedResponse(w)
		return
	}

	token, err := database.RegenerateAccessToken(requestedUser, claims, enc.IdentityDB, enc.Auth)
	if err != nil {
		log.Error().Err(err).Msg("Failed to regenerate access token for user.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondString(w, http.StatusOK, token)
}

func validEmail(email string) bool {

	if len(email) < 6 {
		return false
	}

	return true
}

func validPassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	return true
}
