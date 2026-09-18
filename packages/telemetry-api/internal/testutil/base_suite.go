package testutil

import (
	"context"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"resty.dev/v3"

	"omoidememo.com/telemetry-api/internal/config"
	"omoidememo.com/telemetry-api/internal/httpapi"
	pgstore "omoidememo.com/telemetry-api/internal/store/postgres"
	"omoidememo.com/telemetry-api/internal/telemetry"
)

const (
	testPublicAPIKey     = "test-public-api-key"
	testAdminAPIKey      = "test-admin-api-key"
	postgresImage        = "postgres:16-alpine"
	containerBootTimeout = 60 * time.Second
	containerStopTimeout = 30 * time.Second
)

type BaseSuite struct {
	suite.Suite
	PublicClient *resty.Client
	AdminClient  *resty.Client
	ServerURL    string
	Store        telemetry.Store
	DB           *gorm.DB

	server    *httptest.Server
	container *postgres.PostgresContainer
}

func (s *BaseSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), containerBootTimeout)
	defer cancel()

	container, err := postgres.Run(ctx, postgresImage, postgres.BasicWaitStrategies())
	s.Require().NoError(err, "start postgres container")
	s.container = container

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err, "get postgres connection string")

	db, err := gorm.Open(gormdriver.Open(connStr), &gorm.Config{})
	s.Require().NoError(err, "open gorm connection")
	s.Require().NoError(pgstore.Migrate(db), "run auto-migration")
	s.DB = db
	s.Store = telemetry.NewStore(db)

	cfg := config.Config{
		PublicAPIKey:  testPublicAPIKey,
		AdminAPIKey:   testAdminAPIKey,
		AllowedOrigin: "*",
	}
	handler, _ := httpapi.NewRouter(cfg, s.Store)
	s.server = httptest.NewServer(handler)
	s.ServerURL = s.server.URL

	s.PublicClient = resty.New().
		SetBaseURL(s.server.URL).
		SetHeader("X-Api-Public-Key", testPublicAPIKey)

	s.AdminClient = resty.New().
		SetBaseURL(s.server.URL).
		SetHeader("X-Admin-Api-Key", testAdminAPIKey)
}

func (s *BaseSuite) Reset() {
	s.Require().NoError(s.DB.Where("1 = 1").Delete(&telemetry.Event{}).Error)
}

func (s *BaseSuite) SetupTest() {
	s.Reset()
}

func (s *BaseSuite) TearDownSuite() {
	s.server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), containerStopTimeout)
	defer cancel()
	s.Require().NoError(s.container.Terminate(ctx), "terminate postgres container")
}
