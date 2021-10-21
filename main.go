package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"test.car/handlers"
	"test.car/repository"
	"test.car/service"
)

func main() {

	dsn := flag.String("dsn", "user=postgres password=1234 dbname=test sslmode=disable", "datacars")

	db, err := OpenDB(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	repository := repository.Repository{
		DB:    db,
		Redis: client,
	}

	service := service.Service{}

	service.SetCarRepo(repository)
	service.SetContext(context.Background())

	callHandler := handlers.Handler{}

	callHandler.SetService(service)

	addr := flag.String("addr", ":8080", "Servers address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", callHandler.Home)
	mux.HandleFunc("/get", callHandler.GetCarByVIN)
	mux.HandleFunc("/delete", callHandler.DeleteExitingCar)
	mux.HandleFunc("/create", callHandler.CreateNewCar)
	mux.HandleFunc("/update", callHandler.UpdateExitingCar)

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Printf("Start service at %s", *addr)

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
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
