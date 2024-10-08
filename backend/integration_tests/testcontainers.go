package integration_tests

import (
	"context"
	"net/http"
	"os"
	"os/exec"

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

	httpHandler http.Handler
	srv         *server.Server
}

func (s *TestSuite) SetupSuite() {
	var err error

	s.cfg, err = config.Init("../configs/local.json")
	s.Require().NoError(err, "Failed to initialize config")

	s.logger, err = logger.NewLogger()
	s.Require().NoError(err, "Failed to initialize logger")

	s.pgClient, err = db.NewPostgresClient(context.Background(), s.cfg.PostgresTestConfig.PgTestSource())
	s.Require().NoError(err, "Failed to initialize Postgres client")

	migrateCommand := exec.Command("goose", "-dir", "../migrations", "-allow-missing", "postgres", s.cfg.PostgresTestConfig.PgTestSource(), "up")
	migrateCommand.Stdout = os.Stdout
	migrateCommand.Stderr = os.Stderr
	err = migrateCommand.Run()
	s.Require().NoError(err, "Failed to apply migrations")

	repo, err := repository.NewRepository(s.cfg, s.pgClient.DB, s.logger)
	s.Require().NoError(err, "Failed to initialize repository")

	services, err := service.NewService(context.Background(), repo, s.logger)
	s.Require().NoError(err, "Failed to initialize services")

	h := handler.NewHandler(services.Songs, s.cfg.Handler, s.logger)
	s.httpHandler = h.Init()

	s.srv = server.NewServer(s.cfg.Server, s.httpHandler)
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
}
