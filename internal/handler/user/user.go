package user

import (
	"fmt"
	"log"

	"tg_service/internal/domain"
	"tg_service/internal/service/user"

	"github.com/andReyM228/lib/auth"
	"github.com/andReyM228/lib/errs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService                 user.Service
	tgbot                       *tgbotapi.BotAPI
	loginMap                    map[int64]string
	processingRegistrationUsers *domain.ProcessingRegistrationUsers
}

func NewHandler(service user.Service, tgbot *tgbotapi.BotAPI, loginMap map[int64]string, processingRegistrationUsers *domain.ProcessingRegistrationUsers) Handler {
	return Handler{
		userService:                 service,
		tgbot:                       tgbot,
		loginMap:                    loginMap,
		processingRegistrationUsers: processingRegistrationUsers,
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

func (h Handler) Create(chatID int64, update tgbotapi.Update) error {
	processUser := h.processingRegistrationUsers.GetOrCreate(chatID)

	switch processUser.Step {
	case domain.RegistrationStepStart:
		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите имя")); err != nil {
			log.Fatal(err)
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepName)
	case domain.RegistrationStepName:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}
			return nil
		}
		h.processingRegistrationUsers.SetName(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите фамилию")); err != nil {
			log.Fatal(err)
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepSurname)
	case domain.RegistrationStepSurname:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		h.processingRegistrationUsers.SetSurname(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите номер телефона")); err != nil {
			log.Fatal(err)
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepPhone)

	case domain.RegistrationStepPhone:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		h.processingRegistrationUsers.SetPhone(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите електронную почту")); err != nil {
			log.Fatal(err)
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepEmail)

	case domain.RegistrationStepEmail:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		h.processingRegistrationUsers.SetEmail(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите пароль")); err != nil {
			log.Fatal(err)
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepPassword)

	case domain.RegistrationStepPassword:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				log.Fatal(err)
			}

			return nil
		}

		h.processingRegistrationUsers.SetPassword(chatID, update.Message.Text)

		if err := h.userService.CreateUser(h.processingRegistrationUsers.GetOrCreate(chatID).User); err != nil {
			return err
		}

		h.processingRegistrationUsers.Delete(chatID)
		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация успешна")); err != nil {
			log.Fatal(err)
		}

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
