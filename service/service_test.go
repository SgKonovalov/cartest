package service

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"test.car/definition"
	"test.car/repository"
)

func TestCreateCar(t *testing.T) {

	dsn := "user=postgres password=1234 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	repo := repository.Repository{
		DB: db,
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

	dsn := "user=postgres password=1234 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Errorf("Error by connecting to DB - %v", err)
	}

	repo := repository.Repository{
		DB: db,
	}

	service := Service{
		CarRepo: repo,
	}

	testCarWerr := definition.NewCar("testvinUpdate", "Wrong", "Wrong", "Wrong", 700, 100)

	if err = service.UpdateCar(testCarWerr); err == nil {
		t.Errorf("Can updated car with %s status - %v", testCarWerr.GetCarStatus(), err)
	}
}
