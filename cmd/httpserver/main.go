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

	// WaitGroup waits ran goroutines
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		logger.Info("done")
	}()

	// run server v1 in separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		// set server config
		config := httpserverv1.Config{Addr: ":7373"}
		// init new server
		srv := httpserverv1.New(config, logger)
		// start server (blocing operation)
		if err := srv.Start(); err != nil {
			slog.Error("httpserver.v1 start", uerr.Error(err))
		}
	}()

	// run server v2 in separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		// server's config
		config := httpserverv2.Config{Addr: ":7374"}
		// init and start server
		srv := httpserverv2.New(config, logger)
		if err := srv.Start(); err != nil {
			slog.Error("httpserver.v2 start", uerr.Error(err))
		}
	}()
}
