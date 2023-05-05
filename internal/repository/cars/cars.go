package cars

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"tg_service/internal/domain"
	"tg_service/internal/repository"
)

type Repository struct {
	log    *logrus.Logger
	client *http.Client
}

func NewRepository(log *logrus.Logger, client *http.Client) Repository {
	return Repository{
		log:    log,
		client: client,
	}
}

func (r Repository) Get(id int64) (domain.Car, error) {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/car/%d", id)

	resp, err := r.client.Get(url)
	if err != nil {
		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
	}

	switch resp.StatusCode {
	case 404:
		return domain.Car{}, repository.NotFound{What: "car by id"}

	case 500:
		return domain.Car{}, repository.InternalServerError{Cause: ""}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
	}

	var car domain.Car
	if err := json.Unmarshal(data, &car); err != nil {
		return domain.Car{}, repository.InternalServerError{Cause: err.Error()}
	}

	r.log.Infoln(car)

	return car, nil
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
