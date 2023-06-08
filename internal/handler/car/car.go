package car

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"tg_service/internal/service/car"
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

func (h Handler) Get(id int64, token string) (string, string, error) {
	car, err := h.carService.GetCar(id, token)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("имя: %s, модель: %s, цена: %d", car.Name, car.Model, car.Price), car.Image, nil
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
