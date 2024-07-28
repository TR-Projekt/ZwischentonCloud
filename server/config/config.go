package config

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/pelletier/go-toml"
)

type Config struct {
	ServiceBindAddress        string
	ServiceBindHost           string
	ServicePort               int
	ServiceKey                string
	TLSRootCert               string
	TLSCert                   string
	TLSKey                    string
	JwtExpiration             int
	AccessTokenPrivateKeyPath string
	AccessTokenPublicKeyPath  string
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

func DefaultConfig() *Config {

	// first we try to parse the config at the global configuration path
	if fileExists("/etc/zwischentoncloud.conf") {
		config := ParseConfig("/etc/zwischentoncloud.conf")
		if config != nil {
			return config
		}
	}

	// if there is no global configuration check the current folder for the template config file
	// this is mostly so the application will run in development environment
	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Msg("server initialize: could not read default config file.")
	}
	path = path + "/config_template.toml"
	return ParseConfig(path)
}

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal().Msg("server initialize: could not read config file at '" + cfgFile + "'. Error: " + err.Error())
	}

	serviceBindAdress := content.Get("service.bind-address").(string)
	serviceBindHost := content.Get("service.bind-host").(string)
	servicePort := content.Get("service.port").(int64)
	serviceKey := content.Get("service.key").(string)

	tlsrootcert := content.Get("tls.zwischenton-root-ca").(string)
	tlscert := content.Get("tls.cert").(string)
	tlskey := content.Get("tls.key").(string)

	jwtExpiration := content.Get("jwt.expiration").(int64)
	accessTokenPrivateKeyPath := content.Get("jwt.accessprivatekeypath").(string)
	accessTokenPublicKeyPath := content.Get("jwt.accesspublickeypath").(string)

	dbPassword := content.Get("database.password").(string)

	return &Config{
		ServiceBindAddress:        serviceBindAdress,
		ServiceBindHost:           serviceBindHost,
		ServicePort:               int(servicePort),
		ServiceKey:                serviceKey,
		TLSRootCert:               tlsrootcert,
		TLSCert:                   tlscert,
		TLSKey:                    tlskey,
		JwtExpiration:             int(jwtExpiration),
		AccessTokenPublicKeyPath:  accessTokenPublicKeyPath,
		AccessTokenPrivateKeyPath: accessTokenPrivateKeyPath,
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

func CheckForArguments() {

	if len(os.Args) == 2 {
		if os.Args[1] == "--debug" {

			os.Setenv("DEBUG", "true")
			log.Info().Msg("Running in debug mode")
		}
	}
}

func Debug() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return isPresent
}

func Production() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return !isPresent
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
// see: https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
