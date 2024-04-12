package user

import (
	"fmt"
	"github.com/andReyM228/lib/auth"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/one/chain_client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_service/internal/domain"
	"tg_service/internal/services"
)

type Handler struct {
	userService                 services.UserService
	tgbot                       *tgbotapi.BotAPI
	cache                       services.CacheService
	processingRegistrationUsers *domain.ProcessingRegistrationUsers
	processingLoginUsers        *domain.ProcessingLoginUsers
	chain                       chain_client.Client
}

func NewHandler(userService services.UserService, tgbot *tgbotapi.BotAPI, cache services.CacheService, processingRegistrationUsers *domain.ProcessingRegistrationUsers, processingLoginUsers *domain.ProcessingLoginUsers, chain chain_client.Client) Handler {
	return Handler{
		userService:                 userService,
		tgbot:                       tgbot,
		cache:                       cache,
		processingRegistrationUsers: processingRegistrationUsers,
		processingLoginUsers:        processingLoginUsers,
		chain:                       chain,
	}
}

func (h Handler) GetUser(id int64) (string, error) {
	user, err := h.userService.GetUser(id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("имя: %s, фамилия: %s, телефон: %s, почта: %s, количество машин: %d,", user.Name, user.Surname, user.Phone, user.Email, len(user.Cars)), nil
}

func (h Handler) CreateUser(chatID int64, update tgbotapi.Update) error {
	processUser := h.processingRegistrationUsers.GetOrCreate(chatID)

	switch processUser.Step {
	case domain.RegistrationStepStart:
		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите имя")); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepName)
	case domain.RegistrationStepName:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				h.processingRegistrationUsers.Delete(chatID)

				return err
			}

			h.processingRegistrationUsers.Delete(chatID)

			return nil
		}

		h.processingRegistrationUsers.SetName(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите фамилию")); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepSurname)

	case domain.RegistrationStepSurname:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				h.processingRegistrationUsers.Delete(chatID)

				return err
			}

			h.processingRegistrationUsers.Delete(chatID)

			return nil
		}

		h.processingRegistrationUsers.SetSurname(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите номер телефона")); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepPhone)

	case domain.RegistrationStepPhone:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				h.processingRegistrationUsers.Delete(chatID)

				return err
			}

			h.processingRegistrationUsers.Delete(chatID)

			return nil
		}

		h.processingRegistrationUsers.SetPhone(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите електронную почту")); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepEmail)

	case domain.RegistrationStepEmail:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				h.processingRegistrationUsers.Delete(chatID)

				return err
			}

			h.processingRegistrationUsers.Delete(chatID)

			return nil
		}

		h.processingRegistrationUsers.SetEmail(chatID, update.Message.Text)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите пароль")); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		h.processingRegistrationUsers.UpdateRegistrationStep(chatID, domain.RegistrationStepPassword)

	case domain.RegistrationStepPassword:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация прервана")); err != nil {
				h.processingRegistrationUsers.Delete(chatID)

				return err
			}

			h.processingRegistrationUsers.Delete(chatID)

			return nil
		}

		h.processingRegistrationUsers.SetPassword(chatID, update.Message.Text)

		record, mnemonic, err := h.chain.GenerateAccount(fmt.Sprintf("%d", chatID))
		if err != nil {
			return err
		}

		address, err := record.GetAddress()
		if err != nil {
			return err
		}

		h.processingRegistrationUsers.SetAddress(chatID, address.String())

		if err := h.userService.CreateUser(h.processingRegistrationUsers.GetOrCreate(chatID).User); err != nil {
			h.processingRegistrationUsers.Delete(chatID)

			return err
		}

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "регистрация успешна, ваша мнемоника:")); err != nil {
			return err
		}

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, mnemonic)); err != nil {
			return err
		}

		h.processingRegistrationUsers.Delete(chatID)
	}

	return nil
}

// TODO: переместить создание токена в user-services

func (h Handler) Login(chatID int64, update tgbotapi.Update) error {
	processUser := h.processingLoginUsers.GetOrCreate(chatID)

	switch processUser.Step {
	case domain.LoginStepStart:
		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "введите пароль")); err != nil {
			h.processingLoginUsers.Delete(chatID)

			return errs.InternalError{}
		}

		h.processingLoginUsers.UpdateLoginStep(chatID, domain.LoginStepPassword)

	case domain.LoginStepPassword:
		if update.Message.Text == "/exit" {
			if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "процес логина прерван")); err != nil {
				h.processingLoginUsers.Delete(chatID)

				return errs.InternalError{}
			}

			h.processingLoginUsers.Delete(chatID)

			return nil
		}

		userID, err := h.userService.Login(update.Message.Text, chatID)
		if err != nil {
			h.processingLoginUsers.Delete(chatID)
			return err
		}

		token, err := auth.CreateToken(chatID, userID)
		if err != nil {
			h.processingLoginUsers.Delete(chatID)

			return errs.InternalError{}
		}

		err = h.cache.AddToken(chatID, token)
		if err != nil {
			return err
		}

		h.processingLoginUsers.Delete(chatID)

		if _, err := h.tgbot.Send(tgbotapi.NewMessage(chatID, "логин успешный!")); err != nil {
			return errs.InternalError{}
		}

	}

	return nil
}
