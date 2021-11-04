package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/go-redis/redis"
	"test.car/definition"
)

/*
Структура Repository – ключевая для обработки запросов в БД.
Функции:
1) CreateCar – принимая в качестве аргумента объект типа Car добавляет данный по новому автомобилю в БД.
ВАЖНО: при добавлении нового автомобиля, данные записываются в Redis. Функции обновления БД и Redis выполняются
конкурентно;
2) GetCar – принимая в качестве аргумента VIN осуществляет поиск в БД Redis нужного автомобиля по его VIN;
Все автомобили из реляционной базы загружаются в Redis с помощью отдельной job, описанной в loadallcars;
3) UpdateCar - принимая в качестве аргумента объект типа Car обновляет данные по существующему в БД автомобилю.
В случае если автомобиль для обновления не найден - выходим из функции с ошибкой.
ВАЖНО: при обновлении данных по автомобилю, данные записываются в Redis. Функции обновления БД и Redis выполняются
конкурентно;
4) DeleteCar - принимая в качестве аргумента VIN удаляет данные об автомобиле из БД.
В случае если автомобиль для обновления не найден - выходим из функции с ошибкой.
ВАЖНО: при удалении данных по автомобилю, данные удаляются и в Redis. Функции обновления БД и Redis выполняются
конкурентно.

*Для работы с БД использован драйвер pgx.
*/

type Repository struct {
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Context context.Context
	WG      sync.WaitGroup
}

func (r *Repository) SetContext(Context context.Context) {
	r.Context = Context
}

func (r *Repository) CreateCar(car definition.Car) error {

	r.WG.Add(2)

	go func(definition.Car) error {

		sql := `INSERT INTO datacars (vin, brand, model, price, carstatus, odometer) VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := r.DB.Exec(r.Context, sql, car.GetVIN(), car.GetBrand(), car.GetModel(), car.GetPrice(), car.GetCarStatus(), car.GetOdometer())

		if err != nil {
			return err
		}

		r.WG.Done()

		return nil
	}(car)

	go func(definition.Car) error {

		carInJSON, err := json.Marshal(car)

		if err != nil {
			return err
		}

		r.Redis.Set(fmt.Sprint("car:", car.GetVIN()), carInJSON, 0)

		r.WG.Done()

		return nil
	}(car)

	r.WG.Wait()

	return nil
}

func (r *Repository) GetCar(VIN string) (car definition.Car, err error) {

	carInJSON, err := r.Redis.Get(fmt.Sprint("car:", VIN)).Bytes()

	if err != nil {

		log.Println("Executed 'GET CAR' query from DB")

		sql := `SELECT vin, brand, model, price, carstatus, odometer from datacars WHERE vin = $1`

		allCars, err := r.DB.Query(r.Context, sql, VIN)

		if err != nil {
			log.Println(err)
			return car, err
		}

		defer allCars.Close()

		for allCars.Next() {

			if err := allCars.Scan(&car.VIN, &car.Brand, &car.Model, &car.Price, &car.CarStatus, &car.Odometer); err != nil {
				log.Println(err)
			}
		}
		return car, err
	}

	if err = json.Unmarshal(carInJSON, &car); err != nil {
		return car, err
	}

	return car, nil
}

func (r *Repository) UpdateCar(car definition.Car) error {

	r.WG.Add(2)

	go func(definition.Car) error {

		sql := `UPDATE datacars SET brand = $1, model = $2, price = $3, carstatus = $4, odometer = $5 WHERE vin = $6`

		result, err := r.DB.Exec(r.Context, sql, car.GetBrand(), car.GetModel(), car.GetPrice(), car.GetCarStatus(), car.GetOdometer(), car.GetVIN())

		if err != nil {
			return err
		}

		if result.RowsAffected() != 1 {
			return errors.New("not finded car for update")
		}

		r.WG.Done()

		return nil
	}(car)

	go func(definition.Car) error {

		carInJSON, err := json.Marshal(car)

		if err != nil {
			return err
		}

		r.Redis.Set(fmt.Sprint("car:", car.GetVIN()), carInJSON, 0)

		r.WG.Done()

		return nil
	}(car)

	r.WG.Wait()

	return nil
}

func (r *Repository) DeleteCar(VIN string) error {

	sql := `DELETE FROM datacars WHERE vin = $1`

	result, err := r.DB.Exec(r.Context, sql, VIN)

	if err != nil {
		return err
	}

	switch result.RowsAffected() {
	case 1:
		if err := r.Redis.Del(fmt.Sprint("car:", VIN)).Err(); err != nil {
			return err
		}
	default:
		return errors.New("not finded car for delete")
	}

	return nil
}
