package car

import (
	"github.com/sirupsen/logrus"
	"tg_service/internal/domain"
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
		s.log.Errorln(err)
		return domain.Car{}, err
	}

	return car, nil
}
