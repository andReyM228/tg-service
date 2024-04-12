package repositories

import "tg_service/internal/domain"

type (
	CarRepo interface {
		Get(id int64, token string) (domain.Car, error)
		GetAll(token, label string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(chatID, carID int64, token string) error
		SellCar(chatID, carID int64, token string) error
	}

	UserRepo interface {
		Get(id int64) (domain.User, error)
		Create(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}

	RedisRepo interface {
		Create(key string, value interface{}) error
		GetString(key string) (string, error)
		GetBytes(key string) ([]byte, error)
	}
)
