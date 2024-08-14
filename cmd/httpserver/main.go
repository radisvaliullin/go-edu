package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/radisvaliullin/go-edu/internal/httpserverv1"
	"github.com/radisvaliullin/go-edu/internal/httpserverv2"
	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("http server run")

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		logger.Info("done")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		config := httpserverv1.Config{Addr: ":7373"}
		srv := httpserverv1.New(config, logger)
		if err := srv.Start(); err != nil {
			slog.Error("httpserver.v1 start", uerr.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		config := httpserverv2.Config{Addr: ":7374"}
		srv := httpserverv2.New(config, logger)
		if err := srv.Start(); err != nil {
			slog.Error("httpserver.v2 start", uerr.Error(err))
		}
	}()
}
