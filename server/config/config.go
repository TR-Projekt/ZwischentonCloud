package config

import (
	"github.com/rs/zerolog/log"

	"github.com/pelletier/go-toml"

	servertools "github.com/Festivals-App/festivals-server-tools"
)

type Config struct {
	ServiceBindHost           string
	ServicePort               int
	TLSRootCert               string
	TLSCert                   string
	TLSKey                    string
	JwtExpiration             int
	AccessTokenPrivateKeyPath string
	AccessTokenPublicKeyPath  string
	InfoLog                   string
	TraceLog                  string
	IdentityDBConf            *DBConfig
	ZwischentonDBConf         *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal().Msg("server initialize: could not read config file at '" + cfgFile + "'. Error: " + err.Error())
	}

	serviceBindHost := content.Get("service.bind-host").(string)
	servicePort := content.Get("service.port").(int64)

	tlsrootcert := content.Get("tls.zwischenton-root-ca").(string)
	tlscert := content.Get("tls.cert").(string)
	tlskey := content.Get("tls.key").(string)

	jwtExpiration := content.Get("jwt.expiration").(int64)
	accessTokenPrivateKeyPath := content.Get("jwt.accessprivatekeypath").(string)
	accessTokenPublicKeyPath := content.Get("jwt.accesspublickeypath").(string)

	infoLogPath := content.Get("log.info").(string)
	traceLogPath := content.Get("log.trace").(string)

	dbPassword := content.Get("database.password").(string)

	tlsrootcert = servertools.ExpandTilde(tlsrootcert)
	tlscert = servertools.ExpandTilde(tlscert)
	tlskey = servertools.ExpandTilde(tlskey)
	accessTokenPublicKeyPath = servertools.ExpandTilde(accessTokenPublicKeyPath)
	accessTokenPrivateKeyPath = servertools.ExpandTilde(accessTokenPrivateKeyPath)
	infoLogPath = servertools.ExpandTilde(infoLogPath)
	traceLogPath = servertools.ExpandTilde(traceLogPath)

	return &Config{
		ServiceBindHost:           serviceBindHost,
		ServicePort:               int(servicePort),
		TLSRootCert:               tlsrootcert,
		TLSCert:                   tlscert,
		TLSKey:                    tlskey,
		JwtExpiration:             int(jwtExpiration),
		AccessTokenPublicKeyPath:  accessTokenPublicKeyPath,
		AccessTokenPrivateKeyPath: accessTokenPrivateKeyPath,
		InfoLog:                   infoLogPath,
		TraceLog:                  traceLogPath,
		IdentityDBConf: &DBConfig{
			Dialect:  "mysql",
			Host:     "localhost",
			Port:     int(3306),
			Username: "zwischenton.identity.writer",
			Password: dbPassword,
			Name:     "zwischenton_identity_database",
			Charset:  "utf8",
		},
		ZwischentonDBConf: &DBConfig{
			Dialect:  "mysql",
			Host:     "localhost",
			Port:     int(3306),
			Username: "zwischenton.api.writer",
			Password: dbPassword,
			Name:     "zwischenton_cloud_database",
			Charset:  "utf8",
		},
	}
}
