package integration_tests

import (
	"context"
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
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	pgClient *db.PostgresClient

	cfg    *config.Config
	logger logger.Logger

	handler handler.Handler
	srv     *server.Server
}

func (s *TestSuite) SetupSuite() {
	var err error

	s.cfg, err = config.Init("../configs/local")
	s.Require().NoError(err, "Failed to initialize config")

	s.logger, err = logger.NewLogger()
	s.Require().NoError(err, "Failed to initialize logger")

	s.pgClient, err = db.NewPostgresClient(context.Background(), s.cfg.Postgres.PgSource())
	s.Require().NoError(err, "Failed to initialize Postgres client")

	migrateCommand := exec.Command("goose", "-dir", "migrations", "-allow-missing", "postgres", s.cfg.Postgres.PgSource(), "up")
	migrateCommand.Stdout = os.Stdout
	migrateCommand.Stderr = os.Stderr
	err = migrateCommand.Run()
	s.Require().NoError(err, "Failed to apply migrations")

	repo, err := repository.NewRepository(s.cfg, s.pgClient.DB, s.logger)
	s.Require().NoError(err, "Failed to initialize repository")

	services, err := service.NewService(context.Background(), repo, s.logger)
	s.Require().NoError(err, "Failed to initialize services")

	s.handler = handler.NewHandler(services.Songs, s.cfg.Handler, s.logger)
	s.Require().NotNil(s.handler, "Handler is nil")

	s.srv = server.NewServer(s.cfg.Server, s.handler.Init())
	s.Require().NotNil(s.srv, "Server is nil")
}

func (s *TestSuite) TearDownSuite() {
	if s.srv != nil {
		err := s.srv.Shutdown(context.Background())
		s.Require().NoError(err)
	}

	if s.pgClient != nil {
		s.pgClient.DB.Close()
	}
}

func (s *TestSuite) TearDownTest() {
	ctx := context.Background()

	query := `
		DELETE FROM songs;
		ALTER SEQUENCE songs_id_seq RESTART WITH 1;
	`
	_, err := s.pgClient.DB.ExecContext(ctx, query)
	s.Require().NoError(err)
}

func (s *TestSuite) RunTestServer() {
	if s.srv == nil {
		s.logger.Fatalf("Server is not initialized")
		return
	}

	s.logger.Infof("Starting test server")

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	s.logger.Infof("Shutdown signal received, shutting down server...")
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		s.logger.Errorf("Failed to gracefully shutdown server: %v", err)
	}
}
