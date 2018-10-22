package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load the config file
func Load() *viper.Viper {
	// Configs
	Config, err := readConfig("client_config", map[string]interface{}{
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

//read the config file, helper function
func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath("./configs")
	err := v.ReadInConfig()
	return v, err
}
