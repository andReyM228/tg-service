package car

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"tg_service/internal/service/car"
)

type Handler struct {
	carService car.Service
}

func NewHandler(service car.Service) Handler {
	return Handler{
		carService: service,
	}
}

func (h Handler) Get(id int64) (string, string, error) {
	car, err := h.carService.GetCar(id)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("имя: %s, модель: %s, цена: %d", car.Name, car.Model, car.Price), car.Image, nil
}

func (h Handler) Update(ctx *fiber.Ctx) error {

	return nil
}

func (h Handler) Create(ctx *fiber.Ctx) error {
	return nil
}

func (h Handler) Delete(ctx *fiber.Ctx) error {

	return nil
}
