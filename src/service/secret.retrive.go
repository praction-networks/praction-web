package service

import (
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

func GetJWTSECRET() string {
	JWTConfig, err := config.JWTEnvGet()

	if err != nil {
		logger.Fatal("Unable to get JWT SECRET Config, Existing from application")
	}

	return JWTConfig.JwtSecret
}

func GetStrictAuthKey() string {

	JWTConfig, err := config.JWTEnvGet()

	if err != nil {
		logger.Fatal("Unable to get JWT SECRET Config, Existing from application")
	}

	return JWTConfig.AuthKey

}
