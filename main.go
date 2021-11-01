package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"test.car/handlers"
	"test.car/repository"
	"test.car/service"
)

func main() {

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	//db, err := OpenDB(*dsn)

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	repository := repository.Repository{
		DB:      dbPool,
		Redis:   client,
		Context: context.Background(),
	}

	service := service.Service{}

	service.SetCarRepo(repository)

	callHandler := handlers.Handler{}

	callHandler.SetService(service)

	addr := *flag.String("addr", ":8080", "Servers address")
	flag.Parse()

	router := gin.Default()
	router.GET("/", callHandler.Home)
	router.GET("/get/:vin", callHandler.GetCarByVIN)
	router.DELETE("/delete/:vin", callHandler.DeleteExitingCar)
	router.POST("/create", callHandler.CreateNewCar)
	router.PUT("/update", callHandler.UpdateExitingCar)
	router.Run(addr)

	log.Printf("Start service at %s", addr)

}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
