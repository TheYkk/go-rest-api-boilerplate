package config

import (
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	cfgReader *configReader
)

type (
	Configuration struct {
		MongoSettings MongoSettings
	}
	MongoSettings struct {
		DatabaseName string
		Uri          string
		Timeout      int
	}

	configReader struct {
		configFile string
		v          *viper.Viper
	}
)

func GetAllValues(configPath, configFile string) (configuration *Configuration, err error) {

	newConfigReader(configPath, configFile)
	if err = cfgReader.v.ReadInConfig(); err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "Config: Failed to read config file.")
	}

	err = cfgReader.v.Unmarshal(&configuration)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "Config: Failed to unmarshal yaml file to configuration object.")
	}
	return
}

func newConfigReader(configPath, configFile string) {

	vip := viper.GetViper()
	vip.SetConfigType("yaml")
	vip.SetConfigName(configFile)
	vip.AddConfigPath(configPath)
	cfgReader = &configReader{
		configFile: configFile,
		v:          vip,
	}
}
