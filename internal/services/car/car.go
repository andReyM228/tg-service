package car

import (
	"tg_service/internal/domain"
	"tg_service/internal/repositories"

	"github.com/andReyM228/lib/log"
)

type Service struct {
	carRepo repositories.CarRepo
	log     log.Logger
}

func NewService(carRepo repositories.CarRepo, log log.Logger) Service {
	return Service{
		carRepo: carRepo,
		log:     log,
	}
}

func (s Service) GetCar(carID int64, token string) (domain.Car, error) {
	car, err := s.carRepo.Get(carID, token)
	if err != nil {
		//switch err.(type) {
		//case repositories.BadRequest:
		//	return domain.Car{}, errs.BadRequestError{Cause: err.Error()}
		//case repositories.Unauthorized:
		//	return domain.Car{}, errs.UnauthorizedError{Cause: err.Error()}
		//case repositories.NotFound:
		//	return domain.Car{}, errs.NotFoundError{What: err.Error()}
		//default:
		//	s.log.Error(err.Error())
		//	return domain.Car{}, errs.InternalError{Cause: ""}
		//}
		return domain.Car{}, err
	}

	return car, nil
}

func (s Service) GetCars(token, label string) (domain.Cars, error) {
	cars, err := s.carRepo.GetAll(token, label)
	if err != nil {
		//switch err.(type) {
		//case repositories.BadRequest:
		//	return domain.Cars{}, errs.BadRequestError{Cause: err.Error()}
		//case repositories.Unauthorized:
		//	return domain.Cars{}, errs.UnauthorizedError{Cause: err.Error()}
		//case repositories.NotFound:
		//	return domain.Cars{}, errs.NotFoundError{What: err.Error()}
		//default:
		//	s.log.Error(err.Error())
		//	return domain.Cars{}, errs.InternalError{Cause: ""}
		//}
		return domain.Cars{}, err
	}

	return cars, nil
}

func (s Service) GetUserCars(token string) (domain.Cars, error) {
	cars, err := s.carRepo.GetUserCars(token)
	if err != nil {
		//switch err.(type) {
		//case repositories.BadRequest:
		//	return domain.Cars{}, errs.BadRequestError{Cause: err.Error()}
		//case repositories.Unauthorized:
		//	return domain.Cars{}, errs.UnauthorizedError{Cause: err.Error()}
		//case repositories.NotFound:
		//	return domain.Cars{}, errs.NotFoundError{What: err.Error()}
		//default:
		//	s.log.Error(err.Error())
		//	return domain.Cars{}, errs.InternalError{Cause: ""}
		//}
		return domain.Cars{}, err
	}

	return cars, nil
}

func (s Service) BuyCar(chatID, carID int64) error {
	err := s.carRepo.BuyCar(chatID, carID)
	if err != nil {
		//switch err.(type) {
		//case repositories.BadRequest:
		//	return errs.BadRequestError{Cause: err.Error()}
		//case repositories.Unauthorized:
		//	return errs.UnauthorizedError{Cause: err.Error()}
		//case repositories.NotFound:
		//	return errs.NotFoundError{What: err.Error()}
		//default:
		//	s.log.Error(err.Error())
		//	return errs.InternalError{Cause: ""}
		//}
		return err
	}

	return nil
}

func (s Service) SellCar(chatID, carID int64) error {
	err := s.carRepo.SellCar(chatID, carID)
	if err != nil {
		//switch err.(type) {
		//case repositories.BadRequest:
		//	return errs.BadRequestError{Cause: err.Error()}
		//case repositories.Unauthorized:
		//	return errs.UnauthorizedError{Cause: err.Error()}
		//case repositories.NotFound:
		//	return errs.NotFoundError{What: err.Error()}
		//default:
		//	s.log.Error(err.Error())
		//	return errs.InternalError{Cause: ""}
		//}
		return err
	}

	return nil
}
