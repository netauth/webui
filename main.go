package main

import (
	"os"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/webui/pkg/http"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath("/etc/netauth/")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	viper.Set("client.ServiceName", "webui")

	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "webui",
		Level: hclog.LevelFromString("DEBUG"),
	})

	srv, err := http.New(appLogger)
	if err != nil {
		appLogger.Error("Error during webserver init", "error", err)
		os.Exit(1)
	}

	srv.Serve(":1323")
}
