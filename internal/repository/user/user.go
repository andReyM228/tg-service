package user

import (
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/bus"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/rabbit"
	"github.com/go-playground/validator/v10"

	"net/http"

	"tg_service/internal/domain"
	"tg_service/internal/repository"

	"github.com/andReyM228/lib/log"
)

type Repository struct {
	log       log.Logger
	client    *http.Client
	rabbit    rabbit.Rabbit
	validator *validator.Validate
}

func NewRepository(log log.Logger, client *http.Client, rabbit rabbit.Rabbit, validator *validator.Validate) Repository {
	return Repository{
		log:       log,
		client:    client,
		rabbit:    rabbit,
		validator: validator,
	}
}

func (r Repository) Get(id int64) (domain.User, error) {
	result, err := r.rabbit.Request(bus.SubjectUserServiceGetUserByID, bus.GetUserByIDRequest{ID: id})
	if err != nil {
		return domain.User{}, err
	}

	if result.StatusCode != 200 {
		return domain.User{}, errs.InternalError{Cause: fmt.Sprintf("status code is %d", result.StatusCode)}
	}

	var user domain.User

	err = json.Unmarshal(result.Payload, &user)
	if err != nil {
		return domain.User{}, errs.InternalError{Cause: err.Error()}
	}

	return user, nil
}

func (r Repository) Update() error {

	return nil
}

// сделать норм репоситорские ошибки
func (r Repository) Create(user domain.User) error {
	result, err := r.rabbit.Request(bus.SubjectUserServiceCreateUser, user)
	if err != nil {
		return err
	}

	if result.StatusCode != 200 {
		return errs.InternalError{}
	}

	return nil
}

func (r Repository) Login(password string, chatID int64) (int64, error) {
	request := bus.LoginRequest{
		ChatID:   chatID,
		Password: password,
	}

	result, err := r.rabbit.Request(bus.SubjectUserServiceLoginUser, request)
	if err != nil {
		return 0, err
	}

	if result.StatusCode != 200 {
		return 0, errs.InternalError{}
	}

	var loginResp loginResponse

	if err := json.Unmarshal(result.Payload, &loginResp); err != nil {
		return 0, repository.InternalServerError{Cause: err.Error()}
	}

	err = r.validator.Struct(loginResp)
	if err != nil {
		return 0, repository.BadRequest{Cause: err.Error()}
	}

	return loginResp.UserID, nil
}

func (r Repository) Delete(id int64) error {

	return nil
}
