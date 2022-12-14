package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/shkiperko0/auth-go-ms/delivery/http"
	"github.com/shkiperko0/auth-go-ms/handlers"
	"github.com/shkiperko0/auth-go-ms/interactor"
	"github.com/shkiperko0/auth-go-ms/repositories"
	"github.com/shkiperko0/auth-go-ms/usecases"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	e.Use(http.CORS_Middleware)
	e.HideBanner = true

	// Init redis
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr:     redisAddress,
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})

	grpcServerAddress := os.Getenv("GRPC_SERVER")
	lis, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// Init DB
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_DATABASE")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), nil)

	userRepo := repositories.NewUserRepository(db)
	userIter := interactor.NewUserInteractor(userRepo)
	jwtIter := interactor.NewJwtInteractor(userRepo)
	pwdIter := interactor.NewPasswordInteractor()

	authProvider := handlers.NewAuthProvider(userIter, pwdIter)
	authUC := usecases.NewAuthUseCase(jwtIter, userIter, authProvider)
	userUC := usecases.NewUserUseCase(jwtIter, userIter)

	http.NewCommonHTTPHandler(e)
	http.NewAuthHTTPHandler(e, authUC, userUC)

	if err != nil {
		log.Fatal(err)
		return
	}

	serverAddress := os.Getenv("HTTP_SERVER")
	go func() {
		err = e.Start(serverAddress)

		if err != nil {
			log.Fatal(err)
		}
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
