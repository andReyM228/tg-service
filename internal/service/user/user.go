package user

import (
	"github.com/sirupsen/logrus"
	"tg_service/internal/domain"
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

func (s Service) GetCar(userID int64) (domain.User, error) {
	user, err := s.users.Get(userID)
	if err != nil {
		s.log.Errorln(err)
		return domain.User{}, err
	}

	return user, nil
}
