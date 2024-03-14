package repositories

import "tg_service/internal/domain"

type (
	CarRepo interface {
		Get(id int64, token string) (domain.Car, error)
		GetAll(token, label string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(chatID, carID int64) error
		SellCar(chatID, carID int64) error
	}

	UserRepo interface {
		Get(id int64) (domain.User, error)
		Create(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}
)
