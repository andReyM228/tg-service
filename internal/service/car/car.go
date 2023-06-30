package car

import (
	"tg_service/internal/domain"
	"tg_service/internal/repository"
	"tg_service/internal/repository/cars"

	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
)

type Service struct {
	cars cars.Repository
	log  log.Logger
}

func NewService(cars cars.Repository, log log.Logger) Service {
	return Service{
		cars: cars,
		log:  log,
	}
}

func (s Service) GetCar(carID int64, token string) (domain.Car, error) {
	car, err := s.cars.Get(carID, token)
	if err != nil {
		switch err.(type) {
		case repository.BadRequest:
			return domain.Car{}, errs.BadRequestError{Cause: err.Error()}
		case repository.Unauthorized:
			return domain.Car{}, errs.UnauthorizedError{Cause: err.Error()}
		case repository.NotFound:
			return domain.Car{}, errs.NotFoundError{What: err.Error()}
		default:
			s.log.Error(err.Error())
			return domain.Car{}, errs.InternalError{Cause: ""}
		}
	}

	return car, nil
}

func (s Service) GetCars(token, label string) (domain.Cars, error) {
	cars, err := s.cars.GetAll(token, label)
	if err != nil {
		switch err.(type) {
		case repository.BadRequest:
			return domain.Cars{}, errs.BadRequestError{Cause: err.Error()}
		case repository.Unauthorized:
			return domain.Cars{}, errs.UnauthorizedError{Cause: err.Error()}
		case repository.NotFound:
			return domain.Cars{}, errs.NotFoundError{What: err.Error()}
		default:
			s.log.Error(err.Error())
			return domain.Cars{}, errs.InternalError{Cause: ""}
		}
	}

	return cars, nil
}

func (s Service) BuyCar(chatID, carID int64) error {
	err := s.cars.BuyCar(chatID, carID)
	if err != nil {
		switch err.(type) {
		case repository.BadRequest:
			return errs.BadRequestError{Cause: err.Error()}
		case repository.Unauthorized:
			return errs.UnauthorizedError{Cause: err.Error()}
		case repository.NotFound:
			return errs.NotFoundError{What: err.Error()}
		default:
			s.log.Error(err.Error())
			return errs.InternalError{Cause: ""}
		}
	}

	return nil
}
