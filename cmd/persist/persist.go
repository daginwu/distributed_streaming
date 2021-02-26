package main

import (
	"log"
	"strings"

	"distributed_streaming/cmd/persist/app/sub"
	keyvalue "distributed_streaming/pkg/badger"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {

	// From the environment
	viper.SetEnvPrefix("DISTRIBUTED_STREAMING")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No configuration file was loaded")
	}
}

func main() {

	app := fx.New(
		keyvalue.Modual,
		sub.Modual,
	)

	app.Run()
}
