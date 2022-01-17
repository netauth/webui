package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/webui/pkg/http"
)

func main() {
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
