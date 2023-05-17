package user

import (
	"github.com/andReyM228/lib/auth"
	"github.com/andReyM228/lib/errs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"

	"fmt"
	"log"

	"tg_service/internal/domain"
	"tg_service/internal/service/user"
)

type Handler struct {
	userService user.Service
	tgbot       *tgbotapi.BotAPI
	loginMap    map[int64]string
}

func NewHandler(service user.Service, tgbot *tgbotapi.BotAPI, loginMap map[int64]string) Handler {
	return Handler{
		userService: service,
		tgbot:       tgbot,
		loginMap:    loginMap,
	}
}

func (h Handler) Get(id int64) (string, error) {
	user, err := h.userService.GetUser(id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("имя: %s, фамилия: %s, телефон: %s, почта: %s, количество машин: %d,", user.Name, user.Surname, user.Phone, user.Email, len(user.Cars)), nil
}

func (h Handler) Update(ctx *fiber.Ctx) error {

	return nil
}

func (h Handler) Create(updates tgbotapi.UpdatesChannel, chatID int64) error {
	var user domain.User

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите имя")); err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		user.Name = update.Message.Text
		break
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите фамилию")); err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		user.Surname = update.Message.Text
		break
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите номер телефона")); err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		user.Phone = update.Message.Text
		break
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите електронную почту")); err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		user.Email = update.Message.Text
		break
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите пароль")); err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		user.Password = update.Message.Text
		user.ChatID = update.Message.Chat.ID
		break
	}

	if err := h.userService.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (h Handler) Login(updates tgbotapi.UpdatesChannel, chatID int64) error {
	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите пароль")); err != nil {
		return errs.InternalError{}
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "процес логина прерван")); err != nil {
				return errs.InternalError{}
			}

			return nil
		}

		userID, err := h.userService.Login(update.Message.Text, chatID)
		if err != nil {
			return err
		}

		token, err := auth.CreateToken(chatID, userID)
		if err != nil {
			return errs.InternalError{}
		}

		h.loginMap[chatID] = token

		break
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "логин успешный!")); err != nil {
		return errs.InternalError{}
	}

	return nil
}

func (h Handler) Delete(ctx *fiber.Ctx) error {

	return nil
}

/*
1. при логине генерировать токен
2. токен должен содержать: уникальная инфа про юзера и время истечения
3. написать функцию которая будет докодировать токен
4. добавить авторизацию во все остальные сервисы
*/
