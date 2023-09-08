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

	"tg_service/internal/domain"
	"tg_service/internal/repository"

	"github.com/andReyM228/lib/log"
)

type Repository struct {
	log    log.Logger
	client *http.Client
	rabbit rabbit.Rabbit
}

func NewRepository(log log.Logger, client *http.Client, rabbit rabbit.Rabbit) Repository {
	return Repository{
		log:    log,
		client: client,
		rabbit: rabbit,
	}
}

//func (r Repository) Ge(id int64, token string) (domain.Car, error) {
//	url := fmt.Sprintf("http://localhost:3000/v1/user-service/car/%d", id)
//
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
//	}
//
//	req.Header.Add("Authorization", token)
//
//	resp, err := r.client.Do(req)
//	if err != nil {
//		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
//	}
//
//	if resp.StatusCode > 399 {
//		switch resp.StatusCode {
//		case 404:
//			return domain.Car{}, repository.NotFound{What: "car by id"}
//
//		case 401:
//			return domain.Car{}, repository.Unauthorized{Cause: ""}
//
//		default:
//			return domain.Car{}, repository.InternalServerError{Cause: ""}
//		}
//	}
//
//	data, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
//	}
//
//	var car domain.Car
//	if err := json.Unmarshal(data, &car); err != nil {
//		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
//	}
//
//	r.log.Infof("%v", car)
//
//	return car, nil
//}

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
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/cars/%s", label)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	if resp.StatusCode > 399 {
		switch resp.StatusCode {
		case 404:
			return domain.Cars{}, repository.NotFound{What: "cars"}

		case 401:
			return domain.Cars{}, repository.Unauthorized{Cause: ""}

		default:
			return domain.Cars{}, repository.InternalServerError{Cause: ""}
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	r.log.Infof("%v", cars)

	return cars, nil
}

func (r Repository) GetUserCars(token string) (domain.Cars, error) {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/user-cars")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	if resp.StatusCode > 399 {
		switch resp.StatusCode {
		case 404:
			return domain.Cars{}, repository.NotFound{What: "cars"}

		case 401:
			return domain.Cars{}, repository.Unauthorized{Cause: ""}

		default:
			return domain.Cars{}, repository.InternalServerError{Cause: ""}
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, repository.InternalServerError{Cause: err.Error()}
	}

	r.log.Infof("%v", cars)

	return cars, nil
}

func (r Repository) BuyCar(chatID, carID int64) error {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/buy-car/:chat_id/:car_id")

	url = strings.Replace(url, ":chat_id", strconv.FormatInt(chatID, 10), 1)
	url = strings.Replace(url, ":car_id", strconv.FormatInt(carID, 10), 1)

	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 500:
		return repository.InternalServerError{Cause: ""}
	}

	return nil
}

func (r Repository) SellCar(chatID, carID int64) error {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/sell-car/:chat_id/:car_id")

	url = strings.Replace(url, ":chat_id", strconv.FormatInt(chatID, 10), 1)
	url = strings.Replace(url, ":car_id", strconv.FormatInt(carID, 10), 1)

	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 500:
		return repository.InternalServerError{Cause: ""}
	}

	return nil
}

func (r Repository) Update() error {

	return nil
}

func (r Repository) Create() error {

	return nil
}

func (r Repository) Delete(id int64) error {

	return nil
}
