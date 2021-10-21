package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"test.car/definition"
)

func TestCreateCar(t *testing.T) {

	dsn := "user=postgres password=1234 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	repo := Repository{
		DB: db,
	}

	testCarWTerr := definition.NewCar("TestVINWTerr", "TestBrandWTerr", "TestModelWTerr", "В пути", 700, 100)

	if err = repo.CreateCar(context.Background(), testCarWTerr); err != nil {
		t.Errorf("Can added car with %s vin - %v", testCarWTerr.GetVIN(), err)
	}
}

func TestGetCar(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	repo := Repository{
		Redis: client,
	}

	car, err := repo.GetCar("testvin")
	if err != nil {
		t.Errorf("Can founded car - %v", err)
	}

	if car.GetVIN() != "testvin" {
		t.Errorf("Can founded car VIN %s, expected - %s", car.GetVIN(), "testvin")
	}
}

func TestUpdateCar(t *testing.T) {

	dsn := "user=postgres password=1234 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	repo := Repository{
		DB: db,
	}

	testCarWTerr := definition.NewCar("testvinUpdate", "TestBrandUpdate", "TestModelWupdate", "В пути", 700, 100)

	if err = repo.UpdateCar(context.Background(), testCarWTerr); err != nil {
		t.Errorf("Can updated car with %s vin status - %v", testCarWTerr.GetCarStatus(), err)
	}
}

func TestDeleteCar(t *testing.T) {

	dsn := "user=postgres password=1234 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	repo := Repository{
		DB: db,
	}

	VIN := "testvinDelete"

	if err = repo.DeleteCar(context.Background(), VIN); err != nil {
		t.Errorf("Can deleted car with %s vin status - %v", VIN, err)
	}
}
