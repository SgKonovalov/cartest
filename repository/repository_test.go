package repository

import (
	"context"
	"testing"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"test.car/definition"
)

/*
Тестирование в пакете repository осуществляется на предмет проверки работоспособности
основного функционала, связанного с выполнением запросов в БД.
*/

func TestCreateCar(t *testing.T) {

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	repo := Repository{
		DB:      dbPool,
		Context: context.Background(),
	}

	testCarWTerr := definition.NewCar("TestVINWTerr", "TestBrandWTerr", "TestModelWTerr", "В пути", 700, 100)

	if err = repo.CreateCar(testCarWTerr); err != nil {
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

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	repo := Repository{
		DB:      dbPool,
		Context: context.Background(),
	}

	testCarWTerr := definition.NewCar("testvinUpdate", "TestBrandUpdate", "TestModelWupdate", "В пути", 700, 100)

	if err = repo.UpdateCar(testCarWTerr); err != nil {
		t.Errorf("Can updated car with %s vin status - %v", testCarWTerr.GetCarStatus(), err)
	}
}

func TestDeleteCar(t *testing.T) {

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	repo := Repository{
		DB:      dbPool,
		Context: context.Background(),
	}

	VIN := "testvinDelete"

	if err = repo.DeleteCar(VIN); err != nil {
		t.Errorf("Can deleted car with %s vin status - %v", VIN, err)
	}
}
