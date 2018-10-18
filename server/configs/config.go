package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/spf13/viper"
)

func Load() *viper.Viper {
	// Configs
	Config, err := tools.ReadConfig("server_config", map[string]interface{}{
		"port":        9090,
		"hostname":    "localhost",
		"environment": "debug",
	})
	if err != nil {
		log.Errorf("Error when reading config: %v\n", err)
	}
	return Config
}
