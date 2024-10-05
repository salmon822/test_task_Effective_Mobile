package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/db"
	"github.com/salmon822/test_task/internal/server"
	"golang.org/x/exp/slog"
)

func main() {
	var cfgPath string

	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()

	cfg, err := config.Init(cfgPath)
	if err != nil {
		log.Fatalf("config init err: %v", err)
	}

	ctx := context.Background()

	mux := http.NewServeMux()

	pool, err := db.NewPostgresClient(ctx, cfg.Postgres.PgSource())
	if err != nil {
		log.Fatalf("db init err: %v", err)
	}
	defer pool.Close()

	_ = pool

	srv := server.NewServer(cfg.Server, mux)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server init err: %v", err)
		}
	}()

	slog.Info("server listening", "port", cfg.Server.Port)

	// graceful shutdown here
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	slog.Info("kill signal received, shutting down application....")
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}

	slog.Info("application finished")
}
