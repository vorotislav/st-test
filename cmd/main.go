package main

import (
	"context"
	stdlog "log"
	"os"
	"sync"
	"time"

	"st-test/cmd/util"
	"st-test/internal/http"
	"st-test/internal/logger"
	"st-test/internal/repo"
	"st-test/internal/settings"
	"st-test/internal/signals"
	"st-test/internal/storage"

	"go.uber.org/zap"
)

const (
	serviceShutdownTimeout = 1 * time.Second
)

func main() {
	configFile := util.ParseFlags()

	sets, err := settings.NewSettings(configFile)
	if err != nil {
		stdlog.Fatal(err)
	}

	log, err := logger.New(sets.Log.Level, sets.Log.Format, "stdout", sets.Log.Verbose)
	if err != nil {
		stdlog.Fatal(err)
	}

	nlog := log.Named("main")
	nlog.Debug("Server starting...")
	nlog.Debug(util.Version())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oss := signals.NewOSSignals(ctx)

	oss.Subscribe(func(sig os.Signal) {
		nlog.Info("Stopping by OS Signal...",
			zap.String("signal", sig.String()))

		cancel()
	})

	ls, err := repo.NewRepo(sets.Storage)
	if err != nil {
		stdlog.Fatal(err)
	}

	store := storage.NewStore(log, ls)

	httpService := http.NewService(log, &sets.API, store)

	serviceErrCh := make(chan error, 1)

	var wg sync.WaitGroup

	wg.Add(1)

	go func(errCh chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		defer close(errCh)

		if err := httpService.Run(); err != nil {
			errCh <- err
		}
	}(serviceErrCh, &wg)

	select {
	case err := <-serviceErrCh:
		if err != nil {
			nlog.Error("service error", zap.Error(err))
			cancel()
		}
	case <-ctx.Done():
		nlog.Info("Server stopping...")

		ctxShutdown, ctxCancelShutdown := context.WithTimeout(context.Background(), serviceShutdownTimeout)

		if err := httpService.Stop(ctxShutdown); err != nil {
			nlog.Error("cannot stop server", zap.Error(err))
		}

		store.Stop()
		ls.Close()

		ctxCancelShutdown()
	}

	wg.Wait()
}
