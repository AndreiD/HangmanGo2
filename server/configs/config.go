package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load loads the config
func Load() *viper.Viper {
	// Configs
	Config, err := readConfig("server_config", map[string]interface{}{
		"port":        9090,
		"hostname":    "localhost",
		"environment": "debug",
	})
	if err != nil {
		log.Errorf("Error when reading config: %v\n", err)
	}
	return Config
}

//read the config file, helper function
func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath("./configs")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
