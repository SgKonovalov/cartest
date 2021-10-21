package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"test.car/definition"
	"test.car/service"
)

/*
Handler – структура handler’ов приложения. Включает поле Service для вызова методов этой структуры.

Функции:
Home - стандартный handler обработки стартовой страницы.
CreateNewCar - handler обработки запроса на добавление данных по новому автомобилю.
Принимает параметр – все данные типа Car из JSON;
GetCarByVIN - handler обработки запроса по получение данных по автомобилю.
Принимает параметр – VIN автомобиля из URL;
UpdateExitingCar - handler обработки запроса на обновление данных по существующему автомобилю.
Принимает параметр – все данные типа Car из JSON;
DeleteExitingCar - handler обработки запроса по удаление данных по автомобилю.
Принимает параметр – VIN автомобиля из URL.
*/

type Handler struct {
	Service service.Service
}

func (h *Handler) SetService(Service service.Service) {
	h.Service = Service
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

}

func (h *Handler) CreateNewCar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "Method not allowed", 500)
	}

	log.Println("CREATE: create function in process")

	decodeQ := json.NewDecoder(r.Body)

	var Car *definition.Car

	if err := decodeQ.Decode(&Car); err != nil {
		http.Error(w, "Wrong JSON format", 500)
	}

	if err := h.Service.CreateCar(*Car); err != nil {

		log.Printf("CREATE: couldn't created car - %v", err)
		resErr := definition.ResultError{
			Result: "Couldn't created car",
		}
		errInJSON, _ := json.Marshal(resErr)
		fmt.Fprintln(w, string(errInJSON))
		return
	}

	resOK := definition.ResultOk{
		Result: fmt.Sprintf("car by VIN : %s successfully created", Car.VIN),
	}

	okInJSON, _ := json.Marshal(resOK)

	fmt.Fprint(w, string(okInJSON))

	log.Println("CREATE: update function executing SUCCESSFUL")
}

func (h *Handler) GetCarByVIN(w http.ResponseWriter, r *http.Request) {

	VIN := r.URL.Query().Get("vin")

	log.Printf("SEARCH: looking for car by VIN %s", VIN)

	if VIN == "" {
		log.Println("Incorrect VIN")
		http.Error(w, "incorrect VIN format", 400)
		return
	}

	car, err := h.Service.GetCar(VIN)

	if err != nil {
		log.Println(err)
	}

	if car.VIN == "" {

		log.Printf("Couldn't find car by VIN : %s", VIN)
		resErr := definition.ResultError{
			Result: fmt.Sprintf("Couldn't find car by VIN : %s", VIN),
		}
		errInJSON, _ := json.Marshal(resErr)
		fmt.Fprintln(w, string(errInJSON))
		return
	}

	carInJSON, _ := json.Marshal(car)
	fmt.Fprintln(w, string(carInJSON))

	log.Println("SEARCH: looking for car by VIN executing SUCCESSFUL")
}

func (h *Handler) UpdateExitingCar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "Method not allowed", 500)
	}

	log.Println("UPDATE: update function in process")

	decodeQ := json.NewDecoder(r.Body)

	var Car *definition.Car

	if err := decodeQ.Decode(&Car); err != nil {
		http.Error(w, "Wrong JSON format", 500)
	}

	if err := h.Service.UpdateCar(*Car); err != nil {

		log.Printf("UPDATE: couldn't updated car - %v", err)
		resErr := definition.ResultError{
			Result: "Couldn't updated car",
		}
		errInJSON, _ := json.Marshal(resErr)
		fmt.Fprintln(w, string(errInJSON))
		return
	}

	resOK := definition.ResultOk{
		Result: fmt.Sprintf("car by VIN : %s successfully updated", Car.VIN),
	}

	okInJSON, _ := json.Marshal(resOK)

	fmt.Fprint(w, string(okInJSON))

	log.Println("UPDATE: update function executing SUCCESSFUL")
}

func (h *Handler) DeleteExitingCar(w http.ResponseWriter, r *http.Request) {

	VIN := r.URL.Query().Get("vin")

	if VIN == "" {
		log.Println("DELETE: Incorrect VIN")
		http.Error(w, "incorrect VIN format", 400)
		return
	}

	if err := h.Service.DeleteCar(VIN); err != nil {

		log.Printf("DELETE: couldn't delete car by VIN : %s", VIN)
		resErr := definition.ResultError{
			Result: fmt.Sprintf("Couldn't find car by VIN : %s", VIN),
		}
		errInJSON, _ := json.Marshal(resErr)
		fmt.Fprintln(w, string(errInJSON))
		return
	}

	resOK := definition.ResultOk{
		Result: fmt.Sprintf("Car by VIN : %s successfully deleted", VIN),
	}

	okInJSON, _ := json.Marshal(resOK)

	fmt.Fprint(w, string(okInJSON))

	log.Println("DELETE: deleting car by VIN executing SUCCESSFUL")
}
