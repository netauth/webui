package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/netauth/webui/pkg/http"

	_ "github.com/netauth/netauth/pkg/token/jwt"
	_ "github.com/netauth/netauth/pkg/token/keyprovider/fs"
)

func init() {
	viper.SetDefault("token.duration", time.Minute*5)
	viper.SetDefault("token.keyprovider", "fs")
	viper.SetDefault("token.backend", "jwt-rsa")
}

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
