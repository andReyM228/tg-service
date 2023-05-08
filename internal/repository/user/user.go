package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
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
		return domain.User{}, repository.InternalServerError{Cause: err.Error()}
	}

	switch resp.StatusCode {
	case 404:
		return domain.User{}, repository.NotFound{What: "user by id"}

	case 500:
		return domain.User{}, repository.InternalServerError{Cause: ""}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, repository.InternalServerError{Cause: err.Error()}
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return domain.User{}, repository.InternalServerError{Cause: err.Error()}
	}

	r.log.Infoln(user)

	return user, nil
}

func (r Repository) Update() error {

	return nil
}

// сделать норм репоситорские ошибки
func (r Repository) Create(user domain.User) error {
	url := fmt.Sprintf("http://localhost:3000/v1/user-service/user")

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	reader := io.Reader(buf)

	resp, err := r.client.Post(url, "application/json", reader)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 500:
		return repository.InternalServerError{Cause: ""}
	}

	return nil
}

func (r Repository) Delete(id int64) error {

	return nil
}
