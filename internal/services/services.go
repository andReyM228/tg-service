package services

import "tg_service/internal/domain"

type (
	CarService interface {
		GetCar(carID int64, token string) (domain.Car, error)
		GetCars(token, label string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(chatID, carID int64) error
		SellCar(chatID, carID int64) error
	}

	UserService interface {
		GetUser(userID int64) (domain.User, error)
		CreateUser(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}
)
