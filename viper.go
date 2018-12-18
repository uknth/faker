package faker

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// NewConfig reads the config and loads the properties
func NewConfig() error {

	// config file name would be faker.yaml
	viper.SetConfigType("yaml")
	viper.SetConfigName("fake")

	// paths of config file
	viper.AddConfigPath("/etc/faker")
	viper.AddConfigPath("$HOME/.faker")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Error Reading Config")
	}

	// Watch for changes in config
	viper.WatchConfig()
	viper.OnConfigChange(
		func(event fsnotify.Event) {
			log.Println("Config file changed, restart server")
		},
	)

	return err
}
