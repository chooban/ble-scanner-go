package main

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type MQTTClientConfig struct {
	Broker   string `mapstructure:"broker"`
	ClientId string `mapstructure:"clientID"`
	Channel  string `mapstructure:"channel"`
}

type Configurations struct {
	Mqtt MQTTClientConfig `mapstructure:"mqtt"`
}

func ReadConfiguration() (c Configurations) {
	flag.String("c", "./config.yml", "config file")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println("Failed to read flags")
		fmt.Println(err)
	}

	filepath := viper.GetString("c")

	viper.SetConfigFile(filepath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}

	return c
}
