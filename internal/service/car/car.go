package car

import (
	"errors"
	"github.com/sirupsen/logrus"
	"tg_service/internal/domain"
	"tg_service/internal/domain/errs"
	"tg_service/internal/repository"
	"tg_service/internal/repository/cars"
)

type Service struct {
	cars cars.Repository
	log  *logrus.Logger
}

func NewService(cars cars.Repository, log *logrus.Logger) Service {
	return Service{
		cars: cars,
		log:  log,
	}
}

func (s Service) GetCar(carID int64) (domain.Car, error) {
	car, err := s.cars.Get(carID)
	if err != nil {
		if errors.As(err, &repository.InternalServerError{}) {
			s.log.Errorln(err)
			return domain.Car{}, errs.InternalError{}
		}
		s.log.Debug(err)
		return domain.Car{}, errs.NotFoundError{What: err.Error()}
	}

	return car, nil
}
