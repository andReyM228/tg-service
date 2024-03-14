package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_service/internal/domain"
)

type (
	UserHandler interface {
		GetUser(id int64) (string, error)
		CreateUser(chatID int64, update tgbotapi.Update) error
		Login(chatID int64, update tgbotapi.Update) error
	}

	CarHandler interface {
		GetCar(id int64, token string) (domain.Car, error)
		GetAllCars(token, label string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(chatID, carID int64) error
		SellCar(chatID, carID int64) error
		PrepareCars(car domain.Car, myCar bool) (tgbotapi.FileBytes, tgbotapi.InlineKeyboardMarkup, error)
	}
)
