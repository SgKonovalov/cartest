package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"test.car/definition"
)

/*
Структура AdderOfAllCars – центральная структура, реализующая возможность добавления данных из реляционной БД в Redis.
Функции:
1) AddAllCars – делает запрос в реляционную БД передаёт их функции CarIntoRedis;
2) CarIntoRedis – заполняет Redis получеными данными.
*/

type AdderOfAllCars struct {
	Redis redis.Client
	DB    *sql.DB
	Ctx   context.Context
}

func (aac *AdderOfAllCars) AddAllCars() error {

	sql := `SELECT vin, brand, model, price, carstatus, odometer from datacars`

	allCars, err := aac.DB.QueryContext(aac.Ctx, sql)

	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Car uploaded from DB")

	defer allCars.Close()

	for allCars.Next() {
		var car definition.Car

		if err := allCars.Scan(&car.VIN, &car.Brand, &car.Model, &car.Price, &car.CarStatus, &car.Odometer); err != nil {
			log.Println(err)
			return err
		}

		carInJSON, err := json.Marshal(car)

		if err != nil {
			log.Println(err)
			return err
		}

		if err = aac.CarIntoRedis(carInJSON, car.GetVIN()); err != nil {
			log.Println(err)
			return err
		}

	}

	log.Println("Cars added to Redis")

	return nil
}

func (aac *AdderOfAllCars) CarIntoRedis(carInJSON []byte, VIN string) error {

	if err := aac.Redis.Set(fmt.Sprint(fmt.Sprint("car:", VIN)), carInJSON, 0).Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func main() {

	log.Println("Job is started")

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

	adder := AdderOfAllCars{
		Redis: *client,
		DB:    db,
		Ctx:   context.Background(),
	}

	if err := adder.AddAllCars(); err != nil {
		log.Fatal(err)
	}

	log.Println("Job is done")

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
