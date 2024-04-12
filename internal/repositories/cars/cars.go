package cars

import (
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/bus"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/rabbit"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"tg_service/internal/config"

	"tg_service/internal/domain"
	"tg_service/internal/repositories"

	"github.com/andReyM228/lib/log"
)

type Repository struct {
	log    log.Logger
	client *http.Client
	rabbit rabbit.Rabbit
	cfg    config.Config
}

func NewRepository(log log.Logger, client *http.Client, rabbit rabbit.Rabbit, cfg config.Config) Repository {
	return Repository{
		log:    log,
		client: client,
		rabbit: rabbit,
		cfg:    cfg,
	}
}

func (r Repository) Get(id int64, token string) (domain.Car, error) {
	result, err := r.rabbit.Request(bus.SubjectUserServiceGetCarByID, bus.GetCarByIDRequest{ID: id, Token: token})
	if err != nil {
		return domain.Car{}, err
	}

	if result.StatusCode != 200 {
		return domain.Car{}, errs.InternalError{Cause: fmt.Sprintf("status code is %d", result.StatusCode)}
	}

	var car domain.Car

	err = json.Unmarshal(result.Payload, &car)
	if err != nil {
		return domain.Car{}, errs.InternalError{Cause: err.Error()}
	}

	return car, nil
}

func (r Repository) GetAll(token, label string) (domain.Cars, error) {
	url := fmt.Sprintf(r.cfg.Extra.UrlGetAllCars, label)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	if resp.StatusCode > 399 {
		switch resp.StatusCode {
		case 404:
			return domain.Cars{}, errs.NotFoundError{What: "cars"}

		case 401:
			return domain.Cars{}, errs.UnauthorizedError{}

		default:
			return domain.Cars{}, errs.InternalError{}
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	return cars, nil
}

func (r Repository) GetUserCars(token string) (domain.Cars, error) {
	url := fmt.Sprintf(r.cfg.Extra.UrlGetUserCars)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	if resp.StatusCode > 399 {
		switch resp.StatusCode {
		case 404:
			return domain.Cars{}, errs.NotFoundError{What: "cars"}

		case 401:
			return domain.Cars{}, repositories.Unauthorized{}

		default:
			return domain.Cars{}, errs.InternalError{}
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	return cars, nil
}

// TODO: add header "Auth token"

func (r Repository) BuyCar(chatID, carID int64, token string) error {
	url := fmt.Sprintf(r.cfg.Extra.UrlBuyCar, chatID, carID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	url = strings.Replace(url, ":chat_id", strconv.FormatInt(chatID, 10), 1)
	url = strings.Replace(url, ":car_id", strconv.FormatInt(carID, 10), 1)

	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 500:
		return errs.InternalError{Cause: err.Error()}
	}

	return nil
}

func (r Repository) SellCar(chatID, carID int64, token string) error {
	url := fmt.Sprintf(r.cfg.Extra.UrlSellCar, chatID, carID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	url = strings.Replace(url, ":chat_id", strconv.FormatInt(chatID, 10), 1)
	url = strings.Replace(url, ":car_id", strconv.FormatInt(carID, 10), 1)

	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 500:
		return errs.InternalError{Cause: err.Error()}
	}

	return nil
}
