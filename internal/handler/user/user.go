package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"tg_service/internal/service/user"
)

type Handler struct {
	userService user.Service
}

func NewHandler(service user.Service) Handler {
	return Handler{
		userService: service,
	}
}

func (h Handler) Get(id int64) (string, error) {
	user, err := h.userService.GetCar(id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("имя: %s, фамилия: %s, телефон: %s, почта: %s, количество машин: %d,", user.Name, user.Surname, user.Phone, user.Email, len(user.Cars)), nil
}

func (h Handler) Update(ctx *fiber.Ctx) error {

	return nil
}

func (h Handler) Create(ctx *fiber.Ctx) error {
	return nil
}

func (h Handler) Delete(ctx *fiber.Ctx) error {

	return nil
}
