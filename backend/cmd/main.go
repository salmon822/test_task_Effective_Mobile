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
)

func main() {
	logging, err := logger.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	defer logging.Sync()

	logging.Infof("Starting application")

	var cfgPath string

	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()

	cfg, err := config.Init(cfgPath)
	if err != nil {
		logging.Fatalf("Failed to initialize config: %v", err)
	}

	logging.Infof("Configuration initialized successfully")

	ctx := context.Background()

	pool, err := db.NewPostgresClient(ctx, cfg.Postgres.PgSource())
	if err != nil {
		logging.Fatalf("Failed to initialize database: %v", err)
	}
	defer pool.DB.Close()

	logging.Infof("Database connection established")

	repo, err := repository.NewRepository(cfg, pool.DB, logging)
	if err != nil {
		logging.Panicf("Failed to initialize repository: %v", err)
	}

	logging.Infof("Repository initialized successfully")

	service, err := service.NewService(ctx, repo, logging)
	if err != nil {
		logging.Panicf("Failed to initialize service: %v", err)
	}

	logging.Infof("Services initialized successfully")

	migrateCommand := exec.Command("goose", "-dir", "migrations", "-allow-missing", "postgres", cfg.Postgres.PgSource(), "up")
	migrateCommand.Stdout = os.Stdout
	migrateCommand.Stderr = os.Stderr

	if err := migrations.Migrate(cfg.Postgres.PgSource()); err != nil {
		logging.Fatalf("Failed to apply migrations: %v", err)
	}

	logging.Infof("Migrations applied successfully")

	router := handler.NewHandler(
		service.Songs,
		cfg.Handler,
		logging)

	srv := server.NewServer(cfg.Server, router.Init())

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatalf("Failed to start server: %v", err)
		}
	}()

	logging.Infof("Server listening on port: %d", cfg.Server.Port)

	// graceful shutdown here
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logging.Infof("Shutdown signal received, shutting down server...")

	err = srv.Shutdown(context.Background())
	if err != nil {
		logging.Errorf("Failed to gracefully shutdown server: %v", err)
	}

	logging.Infof("Server shutdown complete")
	logging.Infof("Application finished")
}
