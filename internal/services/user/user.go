package user

import (
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"tg_service/internal/domain"
	"tg_service/internal/repositories"
)

type Service struct {
	userRepo repositories.UserRepo
	log      log.Logger
}

func NewService(userRepo repositories.UserRepo, log log.Logger) Service {
	return Service{
		userRepo: userRepo,
		log:      log,
	}
}

func (s Service) GetUser(userID int64) (domain.User, error) {
	user, err := s.userRepo.Get(userID)
	if err != nil {
		//if errors.As(err, &repositories.InternalServerError{}) {
		//	s.log.Error(err.Error())
		//	return domain.User{}, errs.InternalError{}
		//}
		//s.log.Debug(err.Error())
		//
		//return domain.User{}, errs.NotFoundError{What: err.Error()}
		return domain.User{}, err
	}

	return user, nil
}

func (s Service) CreateUser(user domain.User) error {
	err := s.userRepo.Create(user)
	if err != nil {
		//if errors.As(err, &repositories.InternalServerError{}) {
		//	s.log.Error(err.Error())
		//	return errs.InternalError{}
		//}
		//s.log.Debug(err.Error())
		//
		//return errs.NotFoundError{What: err.Error()}
		return err
	}

	return nil
}

func (s Service) Login(password string, chatID int64) (int64, error) {
	userID, err := s.userRepo.Login(password, chatID)
	if err != nil {
		switch err.(type) {
		case repositories.BadRequest:
			return 0, errs.BadRequestError{Cause: "wrong body"}
		case repositories.Unauthorized:
			return 0, errs.UnauthorizedError{Cause: "wrong password"}
		case repositories.NotFound:
			return 0, errs.NotFoundError{What: "user"}
		default:
			return 0, errs.InternalError{Cause: ""}
		}
	}

	return userID, nil
}
