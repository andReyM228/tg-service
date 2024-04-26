package services

import "tg_service/internal/domain"

type (
	CarService interface {
		GetCar(carID int64, token string) (domain.Car, error)
		GetCars(token, label string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(token, txHash string, carID int64) error
		SellCar(chatID, carID int64, token string) error
	}

	UserService interface {
		GetUser(userID int64) (domain.User, error)
		CreateUser(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}

	CacheService interface {
		AddToken(chatID int64, token string) error
		GetToken(chatID int64) (string, error)
	}
)
