package main

import (
	"fmt"
	"github.com/gabriel-98/bingo-backend/internal/application"
	"github.com/gabriel-98/bingo-backend/internal/application/services"
	//"github.com/gabriel-98/bingo-backend/internal/application/types"
	"github.com/gabriel-98/bingo-backend/internal/domain"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/api/rest"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/providers"
	"github.com/gabriel-98/bingo-backend/internal/infrastructure/repositories/pgrepos"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Database
	dbpath := "host=localhost user=postgres password=password database=Bingo port=5432 sslmode=disable Timezone=America/Bogota"

	// Open db
	db, err := gorm.Open(postgres.Open(dbpath), &gorm.Config{})
	if err != nil {
		log.Printf("error opening db: %s", err)
		return
	}
	log.Printf("Successful connection to the database")

	fmt.Println(db)

	err = db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
	if err != nil {
		log.Printf("error automigrating db: %s", err)
		return
	}
	fmt.Println("Migrate completed successfully")

	// Init repositories.
	dao := InitRepositories()

	// Init providers.
	providerGroup := InitProviders()

	// Init services.
	serviceGroup := InitServices(dao, providerGroup)

	/*user := entities.User{ Username: "abcd", Password: "password" }
	ctx := context.WithValue(context.Background(), "QueryExecutor", db)
	u, err := userRepo.Create(ctx, user)
	fmt.Println(u)*/

	/*card := types.NewRandomCard()
	fmt.Println(card)*/

	// Create the rest server.
	restServer := rest.NewServer(9000, serviceGroup, db)
	
	// Start the rest server.
	restServer.Run()
}

func InitRepositories() *domain.DAO {
	dao := domain.NewDAO(
		pgrepos.NewUserRepo(),
		pgrepos.NewRefreshTokenRepo(),
	)
	return dao
}

func InitProviders() *application.ProviderGroup {
	return application.NewProviderGroup(
		providers.NewPasswordManager(10),
		providers.NewAuthTokenManager(),
	)
}

func InitServices(dao *domain.DAO, providerGroup *application.ProviderGroup) *application.ServiceGroup {
	return application.NewServiceGroup(
		services.NewAuthService(dao, providerGroup.PasswordManager(), providerGroup.AuthTokenManager()),
	)
}