package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/db"
	"github.com/salmon822/test_task/internal/handler"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/repository"
	"github.com/salmon822/test_task/internal/server"
	"github.com/salmon822/test_task/internal/service"
	"github.com/salmon822/test_task/migrations"
	"golang.org/x/exp/slog"
)

func main() {
	logging, err := logger.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	defer logging.Sync()

	var cfgPath string

	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()

	cfg, err := config.Init(cfgPath)
	if err != nil {
		log.Fatalf("config init err: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPostgresClient(ctx, cfg.Postgres.PgSource())
	if err != nil {
		log.Fatalf("db init err: %v", err)
	}
	defer pool.Close()

	repo, err := repository.NewRepository(cfg, pool.DB)
	if err != nil {
		logging.Panic(err)
	}

	service, err := service.NewService(ctx, repo)
	if err != nil {
		logging.Panic(err)
	}

	migrateCommand := exec.Command("goose", "-dir", "migrations", "-allow-missing", "postgres", cfg.Postgres.PgSource(), "up")

	migrateCommand.Stdout = os.Stdout
	migrateCommand.Stderr = os.Stderr

	if err := migrations.Migrate(cfg.Postgres.PgSource()); err != nil {
		log.Fatalf("migrations init err: %v", err)
	}

	router := handler.NewHandler(
		service.Songs,
		cfg.Handler,
		logging)

	srv := server.NewServer(cfg.Server, router.Init())

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
