package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"packs-api/api"
	"packs-api/internal/config"
	"packs-api/internal/utils"
)

const addr = ":8001"

func main() {
	cfg, err := config.NewConfig(addr)
	if err != nil {
		log.Fatalln("could not setup config, ", err)
	}

	logger := utils.NewLogger("dev", "packs-api")
	s := api.NewServer(cfg, logger)

	go func() {
		err = s.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Log.Info("shutting down gracefully...")

	// Context for server shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctxShutDown); err != nil {
		s.Log.Fatalf("server forced to shutdown: %v", err)
	}

	s.Log.Info("server exiting...")
}
