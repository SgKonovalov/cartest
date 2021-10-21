package definition

/*
Центральная структура Car (модель данных).
Содержит все ключевые поля, а также именование полей структуры в формате JSON
+ конструкторы и геттеры/сеттеры для всех полей структуры.
*/

type Car struct {
	VIN       string `json:"VIN"`
	Brand     string `json:"Brand"`
	Model     string `json:"Model"`
	Price     uint   `json:"Price"`
	CarStatus string `json:"CarStatus"`
	Odometer  uint   `json:"Odometer"`
}

func NewCar(VIN, Brand, Model, CarStatus string, Price, Odometer uint) (car Car) {

	car.VIN = VIN
	car.Brand = Brand
	car.Model = Model
	car.Price = Price
	car.CarStatus = CarStatus
	car.Odometer = Odometer

	return car
}

func (c *Car) SetVIN(VIN string) {
	c.VIN = VIN
}

func (c *Car) GetVIN() string {
	return c.VIN
}

func (c *Car) SetBrand(Brand string) {
	c.Brand = Brand
}

func (c *Car) GetBrand() string {
	return c.Brand
}

func (c *Car) SetModel(Model string) {
	c.Model = Model
}

func (c *Car) GetModel() string {
	return c.Model
}

func (c *Car) SetPrice(Price uint) {
	c.Price = Price
}

func (c *Car) GetPrice() uint {
	return c.Price
}

func (c *Car) SetCarStatus(CarStatus string) {
	c.CarStatus = CarStatus
}

func (c *Car) GetCarStatus() string {
	return c.CarStatus
}

func (c *Car) SetOdometer(Odometer uint) {
	c.Odometer = Odometer
}

func (c *Car) GetOdometer() uint {
	return c.Odometer
}
