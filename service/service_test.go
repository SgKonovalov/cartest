package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"test.car/definition"
	"test.car/repository"
)

func TestCreateCar(t *testing.T) {

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	repo := repository.Repository{
		DB: dbPool,
	}

	service := Service{
		CarRepo: repo,
	}

	testCarWerr := definition.NewCar("TestVINWerr", "Wrong", "Wrong", "Wrong", 700, 100)

	if err = service.CreateCar(testCarWerr); err == nil {
		t.Errorf("Can added car with %s status - %v", testCarWerr.GetCarStatus(), err)
	}
}

func TestUpdateCar(t *testing.T) {

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	repo := repository.Repository{
		DB: dbPool,
	}

	service := Service{
		CarRepo: repo,
	}

	testCarWerr := definition.NewCar("testvinUpdate", "Wrong", "Wrong", "Wrong", 700, 100)

	if err = service.UpdateCar(testCarWerr); err == nil {
		t.Errorf("Can updated car with %s status - %v", testCarWerr.GetCarStatus(), err)
	}
}
