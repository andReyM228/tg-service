package user

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

func (r Repository) Get(id int64) (domain.User, error) {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/user/%d", id)

	resp, err := r.client.Get(url)
	if err != nil {
		return domain.User{}, err
	}

	if resp.StatusCode > 300 {
		return domain.User{}, repository.InternalServerError{}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return domain.User{}, err
	}

	r.log.Infoln(user)

	return user, nil
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
