package car

import (
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"tg_service/internal/domain"
	"tg_service/internal/repository"
	"tg_service/internal/repository/cars"
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
