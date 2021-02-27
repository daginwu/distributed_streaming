package main

import (
	"distributed_streaming/cmd/persist/app/api"
	keyvalue "distributed_streaming/pkg/badger"
	"distributed_streaming/pkg/grpc"
	"log"
	"strings"

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
		api.Modual,
		grpc.Modual,
	)

	app.Run()
}
