package user

import (
	"errors"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"tg_service/internal/domain"
	"tg_service/internal/repository"
	"tg_service/internal/repository/user"
)

type Service struct {
	users user.Repository
	log   log.Logger
}

func NewService(users user.Repository, log log.Logger) Service {
	return Service{
		users: users,
		log:   log,
	}
}

func (s Service) GetUser(userID int64) (domain.User, error) {
	user, err := s.users.Get(userID)
	if err != nil {
		if errors.As(err, &repository.InternalServerError{}) {
			s.log.Error(err.Error())
			return domain.User{}, errs.InternalError{}
		}
		s.log.Debug(err.Error())

		return domain.User{}, errs.NotFoundError{What: err.Error()}
	}

	return user, nil
}

func (s Service) CreateUser(user domain.User) error {
	err := s.users.Create(user)
	if err != nil {
		if errors.As(err, &repository.InternalServerError{}) {
			s.log.Error(err.Error())
			return errs.InternalError{}
		}
		s.log.Debug(err.Error())

		return errs.NotFoundError{What: err.Error()}
	}

	return nil
}

func (s Service) Login(password string, chatID int64) (int64, error) {
	userID, err := s.users.Login(password, chatID)
	if err != nil {
		switch err.(type) {
		case repository.BadRequest:
			return 0, errs.BadRequestError{Cause: "wrong body"}
		case repository.Unauthorized:
			return 0, errs.UnauthorizedError{Cause: "wrong password"}
		case repository.NotFound:
			return 0, errs.NotFoundError{What: "user"}
		default:
			return 0, errs.InternalError{Cause: ""}
		}
	}
	return userID, nil
}
