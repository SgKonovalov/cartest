package service

import (
	"errors"
	"log"

	"test.car/definition"
	"test.car/repository"
)

/*
Структура Service – широкая реализация методов Repository.
Так она дополняет функции вышеуказанной структуры логгированием и дополнительными проверками при необходимости.
Константы: статус автомобиля относительно продажи.
Функции:
1) CreateCar – дополнена проверкой цены на равенство 0, а так же switch-case по CarStatus.
Если статус не соответствует жёстко закреплённым константам - выходим из switch с ошибкой;
2) GetCar – дополнена логгированием;
3) UpdateCar – дополнена логгированием;
4) DeleteCar - дополнена логгированием.
*/

const (
	DELIVERY    = "В пути"
	STOCK       = "На складе"
	SOLD        = "Продан"
	OUT_OF_SALE = "Снят с продажи"
)

type Service struct {
	CarRepo repository.Repository
}

func (s *Service) SetCarRepo(CarRepo repository.Repository) {
	s.CarRepo = CarRepo
}

func (s *Service) CreateCar(car definition.Car) error {

	log.Printf("Creating car by VIN : %s started", car.VIN)

	if car.GetPrice() <= 0 {
		return errors.New("car price less or equal to 0")
	}

	switch car.GetCarStatus() {
	case DELIVERY, STOCK, SOLD, OUT_OF_SALE:
		if err := s.CarRepo.CreateCar(car); err != nil {
			log.Printf("ERROR occurated by creating car. Reason is %v", err)
			return err
		}
	default:
		return errors.New("inserted wrong car status")
	}

	log.Printf("Car by VIN : %s created", car.VIN)

	return nil
}

func (s *Service) GetCar(VIN string) (car definition.Car, err error) {

	log.Printf("Started getting car by VIN : %s started", VIN)

	car, err = s.CarRepo.GetCar(VIN)

	if err != nil {
		log.Printf("ERROR occurated by getting car. Reason is %v", err)
		return car, err
	}

	log.Printf("Car by VIN : %s finded", car.VIN)

	return car, nil
}

func (s *Service) UpdateCar(car definition.Car) error {

	log.Printf("Started updating car by VIN : %s started", car.GetVIN())

	switch car.GetCarStatus() {
	case DELIVERY, STOCK, SOLD, OUT_OF_SALE:
		if err := s.CarRepo.UpdateCar(car); err != nil {
			log.Printf("ERROR occurated by updating car. Reason is %v", err)
			return err
		}
	default:
		return errors.New("inserted wrong car status")
	}

	log.Printf("Car by VIN : %s updated", car.VIN)

	return nil
}

func (s *Service) DeleteCar(VIN string) error {

	log.Printf("Started deleting car by VIN : %s started", VIN)

	if err := s.CarRepo.DeleteCar(VIN); err != nil {
		log.Printf("ERROR occurated by deleting car. Reason is %v", err)
		return err
	}

	log.Printf("Car by VIN : %s deleted", VIN)

	return nil
}
