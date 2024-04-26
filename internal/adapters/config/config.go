package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

type (
	Container struct {
		App  *App
		DB   *DB
		HTTP *HTTP
	}

	App struct {
		Name           string `koanf:"name"`
		Env            string `koanf:"env"`
		SecretKey      string `koanf:"secretKey"`
		OtpSecretKey   string `koanf:"otpSecretKey"`
		OtpLength      uint   `koanf:"otp_length"`
		QueryThreshold uint   `koanf:"query_threshold"`
	}

	// Database contains all the environment variables for the database
	DB struct {
		Connection string `koanf:"connection"`
		Host       string `koanf:"host"`
		Port       string `koanf:"port"`
		User       string `koanf:"user"`
		Password   string `koanf:"password"`
		Name       string `koanf:"name"`
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		URL            string   `koanf:"url"`
		Port           string   `koanf:"port"`
		AllowedOrigins []string `koanf:"allowed_origins"`
	}

	Redis struct {
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
	}
)

// New returns a new container
func New() (*Container, error) {
	// Load config file.
	if err := k.Load(file.Provider("config/config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}

	// Load secrets file
	if err := k.Load(file.Provider("secret/secret.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}

	var app App
	var db DB
	var http HTTP

	if err := k.UnmarshalWithConf("", &app, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true}); err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("", &http, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true}); err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("postgres", &db, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true}); err != nil {
		return nil, err
	}

	return &Container{
		App: &app, DB: &db, HTTP: &http,
	}, nil

}
