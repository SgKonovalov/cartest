package repository

import (
	"context"
	"sync"
	"testing"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"test.car/definition"
)

func BenchmarkUpdateCar(b *testing.B) {

	b.ReportAllocs()

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		b.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	repository := Repository{
		DB:      dbPool,
		Redis:   client,
		Context: context.Background(),
		WG:      sync.WaitGroup{},
	}

	b.ResetTimer()
	if err = repository.UpdateCar(definition.NewCar("TestVINWTerr", "BenchTestBrandWTerr", "BenchTestModelWTerr", "В пути", 700, 100)); err != nil {
		b.Errorf("Can't updated car reason - %v", err)
	}
}

func BenchmarkGetCar(b *testing.B) {

	b.ReportAllocs()

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		b.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	repository := Repository{
		DB:      dbPool,
		Redis:   client,
		Context: context.Background(),
		WG:      sync.WaitGroup{},
	}

	b.ResetTimer()

	car, err := repository.GetCar("TestVINWTerr")
	if err != nil {
		b.Errorf("Can't finded car with %s vin - %v", car.GetVIN(), err)
	}

	if car.GetVIN() != "TestVINWTerr" {
		b.Errorf("Can't founded car VIN %s, expected - %s", car.GetVIN(), "TestVINWTerr")
	}
}

func BenchmarkDeleteCar(b *testing.B) {

	b.ReportAllocs()

	databaseUrl := "postgres://postgres:1234@localhost:5432/test"

	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		b.Errorf("Error by connecting to DB - %v", err)
	}

	defer dbPool.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	repository := Repository{
		DB:      dbPool,
		Redis:   client,
		Context: context.Background(),
		WG:      sync.WaitGroup{},
	}

	b.ResetTimer()
	if err = repository.DeleteCar("TestVINWTerr"); err != nil {
		b.Errorf("Can't deleted car reason - %v", err)
	}
}
