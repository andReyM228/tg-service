package car

import (
	"tg_service/internal/domain"
	"tg_service/internal/service/car"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	carService car.Service
	tgbot      *tgbotapi.BotAPI
}

func NewHandler(service car.Service, tgbot *tgbotapi.BotAPI) Handler {
	return Handler{
		carService: service,
		tgbot:      tgbot,
	}
}

func (h Handler) Get(id int64, token string) (domain.Car, error) {
	car, err := h.carService.GetCar(id, token)
	if err != nil {
		return domain.Car{}, err
	}

	return car, nil
}

func (h Handler) GetAll(token, label string) (domain.Cars, error) {
	cars, err := h.carService.GetCars(token, label)
	if err != nil {
		return domain.Cars{}, err
	}

	return cars, nil
}

func (h Handler) BuyCar(chatID, carID int64) error {
	err := h.carService.BuyCar(chatID, carID)
	if err != nil {
		return err
	}

	return nil
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
