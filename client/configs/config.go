package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/spf13/viper"
)

// Load the config file
func Load() *viper.Viper {
	// Configs
	Config, err := tools.ReadConfig("client_config", map[string]interface{}{
		"port":        1234,
		"hostname":    "localhost",
		"environment": "debug",
		"auth": map[string]interface{}{
			"username": "user",
			"password": "password",
			"clientID": "12345",
		},
	})
	if err != nil {
		log.Errorf("Error when reading config: %v\n", err)
	}
	return Config
}
