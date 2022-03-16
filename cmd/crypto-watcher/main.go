package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	watch_matches "github.com/rafadias/crypto-watcher/internal/application/usecase/watchmatches"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/rafadias/crypto-watcher/internal/infra/providers/exchange/coinbase"
)

var build = "develop"

func main() {
	log := log.New(os.Stdout, "CRYPTO WATCHER SERVICE: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("main: error", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// =========================================================================
	// Configuration

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		return err
	}
	cfg := config.New()
	expvar.NewString("build").Set(build)

	// =========================================================================
	// Start Debug service

	log.Println("main: Starting debugging service")

	// Just a go func to enable debugging
	// Start the service listening for debug requests.
	go func() {
		log.Println("main: Debug listening on port: ", cfg.DebugPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.DebugPort), http.DefaultServeMux); err != nil {
			log.Printf("main: Debug closed: %v", err)
		}
	}()

	// =========================================================================
	// Setup App

	log.Println("starting crypto watcher", build)
	defer log.Println("crypto watcher ended")

	coinbaseExchange, err := coinbase.New(exchange.Config{
		BaseURL:       cfg.Exchange.BaseURL,
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	}, log, true)
	if err != nil {
		log.Println("Err trying to create a coinbase instance", err)
		return err
	}

	watchMatchesUseCase := watch_matches.New(log, coinbaseExchange)
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		watchMatchesUseCase.Execute(cfg.Exchange.WindowSize, ctx)
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	// =========================================================================
	// Handle shutdown
	cancelFunc()                               // Signal cancellation to context.Context
	watchMatchesUseCase.Wait()                 // Block here until are workers are done
	log.Println(watchMatchesUseCase.GetVWAP()) // Get the last VWAP

	return nil
}
