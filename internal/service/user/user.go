package user

import (
	"errors"
	"github.com/sirupsen/logrus"
	"tg_service/internal/domain"
	"tg_service/internal/domain/errs"
	"tg_service/internal/repository"
	"tg_service/internal/repository/user"
)

type Service struct {
	users user.Repository
	log   *logrus.Logger
}

func NewService(users user.Repository, log *logrus.Logger) Service {
	return Service{
		users: users,
		log:   log,
	}
}

func (s Service) GetUser(userID int64) (domain.User, error) {
	user, err := s.users.Get(userID)
	if err != nil {
		if errors.As(err, &repository.InternalServerError{}) {
			s.log.Errorln(err)
			return domain.User{}, errs.InternalError{}
		}
		s.log.Debug(err)

		return domain.User{}, errs.NotFoundError{What: err.Error()}
	}

	return user, nil
}

func (s Service) CreateUser(user domain.User) error {
	err := s.users.Create(user)
	if err != nil {
		if errors.As(err, &repository.InternalServerError{}) {
			s.log.Errorln(err)
			return errs.InternalError{}
		}
		s.log.Debug(err)

		return errs.NotFoundError{What: err.Error()}
	}

	return nil
}
