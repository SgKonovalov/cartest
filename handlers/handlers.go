package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

*Для обработки HTTP запросов использован фреймворк Gin.

*/

type Handler struct {
	Service service.Service
}

func (h *Handler) SetService(Service service.Service) {
	h.Service = Service
}

func (h *Handler) Home(c *gin.Context) {}

func (h *Handler) CreateNewCar(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.Header("Allow", http.MethodPost)
		c.IndentedJSON(http.StatusNotFound, gin.H{"wrong method type": "only POST allowed"})
		return
	}

	log.Println("CREATE: create function in process")

	var Car *definition.Car

	if err := c.BindJSON(&Car); err != nil {
		return
	}

	if err := h.Service.CreateCar(*Car); err != nil {
		log.Printf("CREATE: couldn't created car - %v", err)
		c.IndentedJSON(http.StatusNotModified, gin.H{"couldn't create car": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Car by VIN %s - successfully created", Car.GetVIN()))

	log.Println("CREATE: update function executing SUCCESSFUL")
}

func (h *Handler) GetCarByVIN(c *gin.Context) {

	VIN := c.Param("vin")

	log.Printf("SEARCH: looking for car by VIN %s", VIN)

	if VIN == "" {
		log.Println("Incorrect VIN")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"incorrect data format": "wrong VIN entered"})
		return
	}

	car, err := h.Service.GetCar(VIN)

	if err != nil {
		log.Println(err)
	}

	if car.VIN == "" {
		log.Printf("Couldn't find car by VIN : %s", VIN)
		c.IndentedJSON(http.StatusNotFound, gin.H{"incorrect data format": "searched VIN not found"})
	}

	c.IndentedJSON(http.StatusOK, car)

	log.Println("SEARCH: looking for car by VIN executing SUCCESSFUL")
}

func (h *Handler) UpdateExitingCar(c *gin.Context) {

	if c.Request.Method != http.MethodPut {
		c.Header("Allow", http.MethodPut)
		c.IndentedJSON(http.StatusNotFound, gin.H{"wrong method type": "only PUT allowed"})
		return
	}

	log.Println("UPDATE: update function in process")

	var Car *definition.Car

	if err := c.BindJSON(&Car); err != nil {
		c.IndentedJSON(http.StatusNotModified, gin.H{"couldn't update car": "internal server error"})
		return
	}

	if err := h.Service.UpdateCar(*Car); err != nil {

		log.Printf("UPDATE: couldn't updated car - %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"couldn't updated car": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Car by VIN %s - successfully updated", Car.GetVIN()))

	log.Println("UPDATE: update function executing SUCCESSFUL")
}

func (h *Handler) DeleteExitingCar(c *gin.Context) {

	if c.Request.Method != http.MethodDelete {
		c.Header("Allow", http.MethodDelete)
		c.IndentedJSON(http.StatusNotFound, gin.H{"wrong method type": "only DELETE allowed"})
		return
	}

	VIN := c.Param("vin")

	if VIN == "" {
		log.Println("DELETE: Incorrect VIN")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"incorrect data format": "wrong VIN entered"})
		return
	}

	if err := h.Service.DeleteCar(VIN); err != nil {

		log.Printf("DELETE: couldn't delete car by VIN : %s", VIN)
		c.IndentedJSON(http.StatusNotModified, gin.H{"couldn't delete car": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Car by VIN %s - successfully deleted", VIN))

	log.Println("DELETE: deleting car by VIN executing SUCCESSFUL")
}
