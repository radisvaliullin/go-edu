package main

import (
	"log/slog"
	"os"

	httpserverv1 "github.com/radisvaliullin/go-edu/internal/httpseverv1"
	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("http server run")

	config := httpserverv1.Config{Addr: ":7373"}
	srv := httpserverv1.New(config, logger)
	if err := srv.Start(); err != nil {
		slog.Error("httpserver.v1 start", uerr.Error(err))
	}
}
