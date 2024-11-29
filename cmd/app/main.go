package main

import (
	"fmt"
	"github.com/gabriel-98/bingo-backend/internal/application"
	"github.com/gabriel-98/bingo-backend/internal/application/services"
	//"github.com/gabriel-98/bingo-backend/internal/application/types"
	"github.com/gabriel-98/bingo-backend/internal/config"
	"github.com/gabriel-98/bingo-backend/internal/domain"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/api/rest"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/providers"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/repositories/pgrepos"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	configFilepath = "config/config.yaml"
)

func main() {
	// Create and configure the initialization logger.
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp: true,
	})

	// Load configuration.
	log.Info("Loading configuration ...")
	cfg, err := config.LoadConfig(configFilepath)
	if err != nil {
		log.Error("Error loading configuration: %s", err)
		return
	}
	log.Info("Configuration loaded")

	// Open database.
	log.Info("Connecting to the database ...")
	dbConfig := cfg.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.SslMode,
		dbConfig.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Error opening database: %s", err)
		return
	}
	log.Info("Successful connection to the database")

	// Migrate models to the database.
	log.Info("Migrating models to the database ...")
	err = db.AutoMigrate(
		&entities.User{},
		&entities.RefreshToken{},
	)
	if err != nil {
		log.Error("Error migrating models to the database: %s", err)
		return
	}
	log.Info("Model migration to the database completed")

	// Initialization of repositories.
	log.Info("Initializing repositories ...")
	dao := InitRepositories()
	log.Info("Repositories initialized successfully")

	// Initialization of providers.
	log.Info("Initializing providers ...")
	providerGroup := InitProviders(cfg)
	log.Info("Providers initialized successfully")

	// Initialization of services.
	log.Info("Initializing services ...")
	serviceGroup := InitServices(dao, providerGroup)
	log.Info("Services initialized successfully")

	// Initialization of the rest server.
	log.Info("Initializing REST server ...")
	port := cfg.Server.Port
	restServer := rest.NewServer(port, serviceGroup, db)
	log.Info("REST server initialized successfully")
	
	// Run the REST server.
	log.Info("Running the REST server ...")
	restServer.Run()
	log.Info("REST server has been stopped")
}

func InitRepositories() *domain.DAO {
	return domain.NewDAO(
		pgrepos.NewUserRepo(),
		pgrepos.NewRefreshTokenRepo(),
	)
}

func InitProviders(cfg *config.Config) *application.ProviderGroup {
	return application.NewProviderGroup(
		providers.NewPasswordManager(cfg.PasswordHashing.Bcrypt.Cost),
		providers.NewAuthTokenManager(cfg.Auth),
	)
}

func InitServices(dao *domain.DAO, providerGroup *application.ProviderGroup) *application.ServiceGroup {
	passwordManager := providerGroup.PasswordManager()
	authTokenManager := providerGroup.AuthTokenManager()
	return application.NewServiceGroup(
		services.NewAuthService(dao, passwordManager, authTokenManager),
	)
}