package service

import (
	"context"
	"errors"
	"log"

	"test.car/definition"
	"test.car/repository"
)

/*
Структура Service – широкая реализация методов Repository.
Так она дополняет функции вышеуказанной структуры логгированием и дополнительными проверками при необходимости.
Функции:
1) CreateCar – дополнена проверкой цены на равенство 0, а так же switch-case по CarStatus;
2) GetCar – дополнена логгированием;
3) UpdateCar – дополнена логгированием;
4) DeleteCar - дополнена логгированием.
*/

type Service struct {
	CarRepo repository.Repository
	Ctx     context.Context
}

func (s *Service) SetCarRepo(CarRepo repository.Repository) {
	s.CarRepo = CarRepo
}

func (s *Service) SetContext(Ctx context.Context) {
	s.Ctx = Ctx
}

func (s *Service) CreateCar(car definition.Car) error {

	log.Printf("Creating car by VIN : %s started", car.VIN)

	if car.GetPrice() <= 0 {
		return errors.New("car price less or equal to 0")
	}

	switch car.GetCarStatus() {
	case "В пути", "На складе", "Продан", "Снят с продажи":
		if err := s.CarRepo.CreateCar(s.Ctx, car); err != nil {
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
	case "В пути", "На складе", "Продан", "Снят с продажи":
		if err := s.CarRepo.UpdateCar(s.Ctx, car); err != nil {
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

	if err := s.CarRepo.DeleteCar(s.Ctx, VIN); err != nil {
		log.Printf("ERROR occurated by deleting car. Reason is %v", err)
		return err
	}

	log.Printf("Car by VIN : %s deleted", VIN)

	return nil
}
