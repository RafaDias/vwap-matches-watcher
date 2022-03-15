package main

import (
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
	log.Println("starting crypto watcher", build)
	defer log.Println("crypto watcher ended")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()

	log.Println("main: Starting debugging support")
	go func() {
		log.Println("main: Debug listening on port: ", cfg.DebugPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.DebugPort), http.DefaultServeMux); err != nil {
			log.Printf("main: Debug closed: %v", err)
		}
	}()

	expvar.NewString("build").Set(build)

	go func() {
		log.Println("initializing service")
		coinbaseExchange, err := coinbase.New(exchange.Config{
			BaseURL:       cfg.Exchange.BaseURL,
			Channels:      cfg.Exchange.Channels,
			Subscriptions: cfg.Exchange.Subscriptions,
		})
		if err != nil {
			log.Println("Err trying to create a coinbase instance", err)
		}

		watchMatchesUseCase := watch_matches.New(log, coinbaseExchange)
		watchMatchesUseCase.Execute(cfg.Exchange.WindowSize)
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	return nil
}
