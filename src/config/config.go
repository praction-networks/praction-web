package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	EnvironmentEnv EnvironmentConfig `mapstructure:"ENVIRONMENT"`
	LoggerEnv      LoggerConfig      `mapstructure:"LOGGER"`
	MongoDBEnv     MongoConfig       `mapstructure:"MONGODB"`
	ChiEnv         ChiConfig         `mapstructure:"SERVER"`
	JWTEnv         JWTConfig         `mapstructure:"AUTH"`
	MSG91Env       MSG91Config       `mapstructure:"MSG91"`
	PostalEnv      POSTALConfig      `mapstructure:"POSTAL"`
	CloudinaryEnv  CLOUDINARYConfig  `mapstructure:"CLOUDINARY"`
}

type CLOUDINARYConfig struct {
	CloudName string `mapstructure:"CLOUDNAME"`
	ApiKey    string `mapstructure:"APIKEY"`
	ApiSecret string `mapstructure:"APISECRET"`
}

type POSTALConfig struct {
	ServerURL string `mapstructure:"SERVERURL"`
	ApiKey    string `mapstructure:"APIKEY"`
	EmailFrom string `mapstructure:"EMAILFORM"`
}

type MSG91Config struct {
	AuthKey    string `mapstructure:"AUTHKEY"`
	TemplateID string `mapstructure:"TEMPLATEID"`
}

type JWTConfig struct {
	JwtSecret string `mapstructure:"SECRET"`
	AuthKey   string `mapstructure:"AUTH_KEY"`
}

type LoggerConfig struct {
	LogLevel       string `mapstructure:"LOGLEVEL"`
	Console        bool   `mapstructure:"CONSOLE"`
	ElkEnabled     bool   `mapstructure:"ELKENABLED"`
	ElkHost        string `mapstructure:"ELKHOST"`
	ElkPort        string `mapstructure:"ELKPORT"`
	ElkSearchIndex string `mapstructure:"ELKSEARCHINDEX"`
}

type MongoConfig struct {
	Host       string `mapstructure:"HOST"`
	Port       int    `mapstructure:"PORT"`
	TLSEnabled bool   `mapstructure:"TLSENABLED"`
	DBName     string `mapstructure:"DATABASE"`
	DBUser     string `mapstructure:"USERNAME"`
	DBPassword string `mapstructure:"PASSWORD"`
}

type ChiConfig struct {
	Port string `mapstructure:"PORT"`
}

type EnvironmentConfig struct {
	Env     string `mapstructure:"ENV"`
	App     string `mapstructure:"APP"`
	Version string `mapstructure:"VERSION"`
}

func EnvGet() (EnvConfig, error) {
	// viper.SetConfigFile("src/config/environment.yaml")
	viper.SetConfigFile("/src/config/environment.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return EnvConfig{}, errors.New("error reading config file: " + err.Error())
	}

	var envConfig EnvConfig
	if err := viper.Unmarshal(&envConfig); err != nil {
		return EnvConfig{}, errors.New("unable to decode into struct: " + err.Error())
	}

	return envConfig, nil
}

func LoggerEnvGet() (LoggerConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get Logger Env Config Error: %v", err)
		return LoggerConfig{}, errors.New(msg)
	}

	if !envConfig.LoggerEnv.Console && !envConfig.LoggerEnv.ElkEnabled {
		return LoggerConfig{}, errors.New("at least one logging option (Console or ELK) must be enabled")
	}

	return envConfig.LoggerEnv, nil
}

func MongoEnvGet() (MongoConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get MongoDB Env Config Error: %v", err)
		return MongoConfig{}, errors.New(msg)
	}

	return envConfig.MongoDBEnv, nil
}

func ChiEnvGet() (ChiConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get Fiber Env Config Error: %v", err)
		return ChiConfig{}, errors.New(msg)
	}

	return envConfig.ChiEnv, nil
}

func EnvironmentEnvGet() (EnvironmentConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get Environment Env Config Error: %v", err)
		return EnvironmentConfig{}, errors.New(msg)
	}

	return envConfig.EnvironmentEnv, nil
}

func JWTEnvGet() (JWTConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get JWT Env Config Error: %v", err)
		return JWTConfig{}, errors.New(msg)
	}

	return envConfig.JWTEnv, nil
}

func MSG91EnvGet() (MSG91Config, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get MSG91 Env Config Error: %v", err)
		return MSG91Config{}, errors.New(msg)
	}

	return envConfig.MSG91Env, nil
}

func POSTALEnvGet() (POSTALConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get POSTAL Env Config Error: %v", err)
		return POSTALConfig{}, errors.New(msg)
	}

	return envConfig.PostalEnv, nil
}

func CLOUDINARYEnvGet() (CLOUDINARYConfig, error) {
	envConfig, err := EnvGet()
	if err != nil {
		msg := fmt.Sprintf("Unable to get CLOUDINARY Env Config Error: %v", err)
		return CLOUDINARYConfig{}, errors.New(msg)
	}

	return envConfig.CloudinaryEnv, nil
}
