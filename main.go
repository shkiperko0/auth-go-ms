package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"eam-auth-go-ms/xarv"
)

func CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Response()
		req := c.Request()
		fmt.Println(req.Method, " ", req.URL.Path, " ", req.Host, " ", c.Request().Header.Get("Origin"))

		if c.Request().Method == "OPTIONS" {
			//origin := c.Request().Header.Get("Origin") 1231
			//res.Header().Set("Access-Control-Allow-Origin", origin)
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Header().Set("Access-Control-Allow-Methods" /*entry.Method+*/, "GET,HEAD,PUT,PATCH,POST,DELETE")
			res.Header().Set("Access-Control-Allow-Headers", "*")
			res.WriteHeader(http.StatusNoContent)
			return nil
			//} else if len(res.Header().Get("Access-Control-Allow-Origin")) < 1 {
			//origin := c.Request().Header.Get("Origin")
			//res.Header().Set("Access-Control-Allow-Origin", origin)
		}

		//c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func main() {
	// Init echo server
	e := echo.New()

	e.Use(CORS)
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

	authServerAddress := os.Getenv("AUTH_SERVER")
	xarv.NewProxyHttpHandler(db, e, authServerAddress)

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
