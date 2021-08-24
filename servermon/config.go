package main


import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/mitchellh/mapstructure"
)

var (
	viperCfg *viper.Viper
	config  =  &Config{}

)

type Config struct {
	Address string            `mapstructure:"address"`
}

func InitConfig()  {

	viperCfg = viper.New()
	viperCfg.SetConfigName("config") 
	viperCfg.SetConfigType("json")
	viperCfg.AddConfigPath("./")
	viperCfg.ReadInConfig()
	if err := viperCfg.Unmarshal(config); err != nil {
	   fmt.Errorf("error parsing config file: %s", err)
	}
 
}