package server

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	festivalspki "github.com/Festivals-App/festivals-pki"
	"github.com/TR-Projekt/zwischentoncloud/server/config"
	"github.com/TR-Projekt/zwischentoncloud/server/database"
	"github.com/TR-Projekt/zwischentoncloud/server/handler"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/TR-Projekt/zwischentoncloud/server/servertools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

// Server has router and db instances
type Server struct {
	Router        *chi.Mux
	zwischentonDB *sql.DB
	identityDB    *sql.DB
	Config        *config.Config
	TLSConfig     *tls.Config
	Auth          *token.AuthService
	Validator     *token.ValidationService
}

func NewServer(config *config.Config) *Server {
	server := &Server{}
	server.initialize(config)
	return server
}

// Initialize the server with predefined configuration
func (s *Server) initialize(config *config.Config) {

	s.Config = config
	s.Router = chi.NewRouter()

	s.setZwischentonDatabase()
	s.setIdentityDatabase()
	s.setIdentityService(s.identityDB)
	s.setTLSHandling()
	s.setMiddleware()
	s.setRoutes(config)
}

func (s *Server) setIdentityService(identityDB *sql.DB) {

	config := s.Config
	s.Auth = token.NewAuthService(config.AccessTokenPrivateKeyPath, config.AccessTokenPublicKeyPath, config.JwtExpiration, "de.zwischenton.cloud.issuer")
	apiKeys, err := database.GetAllAPIKeys(identityDB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load api keys from zwischenton identity database.")
	}
	s.Validator = token.NewValidationService(s.Auth, apiKeys)
}

func (s *Server) setZwischentonDatabase() {

	dbConfig := s.Config.ZwischentonDBConf

	dbURI := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		"zwischenton.api.writer",
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Charset,
	)

	db, err := sql.Open(dbConfig.Dialect, dbURI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open zwischenton database handle.")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to zwischenton database.")
	}

	db.SetConnMaxIdleTime(time.Minute * 1)
	db.SetConnMaxLifetime(time.Minute * 5)

	s.zwischentonDB = db
}

func (s *Server) setIdentityDatabase() {

	dbConfig := s.Config.IdentityDBConf

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		"zwischenton.identity.writer",
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Charset,
	)
	db, err := sql.Open(dbConfig.Dialect, dbURI)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open identity database handle.")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to identity database.")
	}

	db.SetConnMaxIdleTime(time.Minute * 1)
	db.SetConnMaxLifetime(time.Minute * 5)

	s.identityDB = db
}

func (s *Server) setTLSHandling() {
	tlsConfig := &tls.Config{
		ClientAuth:     tls.RequireAndVerifyClientCert,
		GetCertificate: festivalspki.LoadServerCertificateHandler(s.Config.TLSCert, s.Config.TLSKey, s.Config.TLSRootCert),
	}
	s.TLSConfig = tlsConfig
}

func (s *Server) setMiddleware() {

	// tell the router which middleware to use
	s.Router.Use(
		// used to log the request to the console
		servertools.Middleware(servertools.TraceLogger("/var/log/zwischenton/trace.log")),
		// tries to recover after panics (?)
		middleware.Recoverer,
	)
}

// setRouters sets the all required routers
func (s *Server) setRoutes(config *config.Config) {

	s.Router.Get("/version", s.handleRequest(handler.GetVersion))
	s.Router.Get("/info", s.handleRequest(handler.GetInfo))
	s.Router.Get("/health", s.handleRequest(handler.GetHealth))

	s.Router.Post("/update", s.handleRequest(handler.MakeUpdate))
	s.Router.Get("/log", s.handleRequest(handler.GetLog))
	s.Router.Get("/log/trace", s.handleRequest(handler.GetTraceLog))

	s.Router.Post("/users/signup", s.handleAPIRequest(handler.Signup))
	s.Router.Get("/users/login", s.handleAPIRequest(handler.Login))
	s.Router.Get("/users/refresh", s.handleRequest(handler.Refresh))

	s.Router.Get("/zwischentons", s.handleAPIRequest(handler.GetZwischentons))
	s.Router.Get("/zwischentons/{objectID}", s.handleAPIRequest(handler.GetZwischenton))
	s.Router.Post("/zwischentons", s.handleRequest(handler.CreateZwischenton))

	/*
		s.Router.Get("/festivals", s.handleAPIRequest(handler.GetFestivals))
		s.Router.Get("/festivals/{objectID}", s.handleAPIRequest(handler.GetFestival))
		s.Router.Get("/festivals/{objectID}/events", s.handleAPIRequest(handler.GetFestivalEvents))
		s.Router.Get("/festivals/{objectID}/image", s.handleAPIRequest(handler.GetFestivalImage))
		s.Router.Get("/festivals/{objectID}/links", s.handleAPIRequest(handler.GetFestivalLinks))
		s.Router.Get("/festivals/{objectID}/place", s.handleAPIRequest(handler.GetFestivalPlace))
		s.Router.Get("/festivals/{objectID}/tags", s.handleAPIRequest(handler.GetFestivalTags))

		s.Router.Post("/festivals", s.handleRequest(handler.CreateFestival))
		s.Router.Patch("/festivals/{objectID}", s.handleRequest(handler.UpdateFestival))
		s.Router.Delete("/festivals/{objectID}", s.handleRequest(handler.DeleteFestival))
		s.Router.Post("/festivals/{objectID}/events/{resourceID}", s.handleRequest(handler.SetEventForFestival))
		s.Router.Post("/festivals/{objectID}/image/{resourceID}", s.handleRequest(handler.SetImageForFestival))
		s.Router.Post("/festivals/{objectID}/links/{resourceID}", s.handleRequest(handler.SetLinkForFestival))
		s.Router.Post("/festivals/{objectID}/place/{resourceID}", s.handleRequest(handler.SetPlaceForFestival))
		s.Router.Post("/festivals/{objectID}/tags/{resourceID}", s.handleRequest(handler.SetTagForFestival))
		s.Router.Delete("/festivals/{objectID}/image/{resourceID}", s.handleRequest(handler.RemoveImageForFestival))
		s.Router.Delete("/festivals/{objectID}/links/{resourceID}", s.handleRequest(handler.RemoveLinkForFestival))
		s.Router.Delete("/festivals/{objectID}/place/{resourceID}", s.handleRequest(handler.RemovePlaceForFestival))
		s.Router.Delete("/festivals/{objectID}/tags/{resourceID}", s.handleRequest(handler.RemoveTagForFestival))
	*/
}

func (s *Server) Run(conf *config.Config) {

	server := http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,

		Addr:      conf.ServiceBindHost + ":" + strconv.Itoa(conf.ServicePort),
		Handler:   s.Router,
		TLSConfig: s.TLSConfig,
	}

	//server.SetKeepAlivesEnabled(false)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal().Err(err).Str("type", "server").Msg("Failed to run server")
	}
}

type APIKeyAuthenticatedHandlerFunction func(enc *handler.HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request)

func (s *Server) handleAPIRequest(requestHandler APIKeyAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey := token.GetAPIToken(r)
		if !slices.Contains((*s.Validator.APIKeys), apikey) {
			claims := token.GetValidClaims(r, s.Validator)
			if claims == nil {
				servertools.UnauthorizedResponse(w)
				return
			}
		}
		enc := &handler.HandlerEncapsulator3000{ZwischentonDB: s.zwischentonDB, IdentityDB: s.identityDB, Auth: s.Auth, Validator: s.Validator}
		requestHandler(enc, w, r)
	})
}

type JWTAuthenticatedHandlerFunction func(claims *token.UserClaims, enc *handler.HandlerEncapsulator3000, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequest(requestHandler JWTAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims := token.GetValidClaims(r, s.Validator)
		if claims == nil {
			servertools.UnauthorizedResponse(w)
			return
		}
		enc := &handler.HandlerEncapsulator3000{ZwischentonDB: s.zwischentonDB, IdentityDB: s.identityDB, Auth: s.Auth, Validator: s.Validator}
		requestHandler(claims, enc, w, r)
	})
}
