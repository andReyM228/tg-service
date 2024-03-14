package car

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_service/internal/domain"
	"tg_service/internal/services"
)

type Handler struct {
	carService services.CarService
	tgbot      *tgbotapi.BotAPI
}

func NewHandler(carService services.CarService, tgbot *tgbotapi.BotAPI) Handler {
	return Handler{
		carService: carService,
		tgbot:      tgbot,
	}
}

func (h Handler) GetCar(id int64, token string) (domain.Car, error) {
	car, err := h.carService.GetCar(id, token)
	if err != nil {
		return domain.Car{}, err
	}

	return car, nil
}

func (h Handler) GetAllCars(token, label string) (domain.Cars, error) {
	cars, err := h.carService.GetCars(token, label)
	if err != nil {
		return domain.Cars{}, err
	}

	return cars, nil
}

func (h Handler) GetUserCars(token string) (domain.Cars, error) {
	cars, err := h.carService.GetUserCars(token)
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

func (h Handler) SellCar(chatID, carID int64) error {
	err := h.carService.SellCar(chatID, carID)
	if err != nil {
		return err
	}

	return nil
}
