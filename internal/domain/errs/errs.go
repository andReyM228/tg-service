package errs

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TgError struct {
	Err    error
	ChatID int64
}

type InternalError struct {
	Cause string
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Cause)
}

type NotFoundError struct {
	What string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.What)
}

type BadRequestError struct {
	Cause string
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Cause)
}

type ForbiddenError struct {
	Cause string
}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Cause)
}

type UnauthorizedError struct {
	Cause string
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Cause)
}

func HandleError(err error, log *logrus.Logger, tgbot *tgbotapi.BotAPI, chatID int64) {
	switch err.(type) {
	case *NotFoundError:
		log.Debug(err)

		_, err := tgbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		if err != nil {
			log.Error(err)
		}
	case *BadRequestError:
		log.Debug(err)

		_, err := tgbot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		if err != nil {
			log.Error(err)
		}
	case *ForbiddenError:
		log.Debug(err)

		_, err := tgbot.Send(tgbotapi.NewMessage(chatID, "вам это запрещено"))
		if err != nil {
			log.Error(err)
		}
	case *UnauthorizedError:
		log.Debug(err)

		_, err := tgbot.Send(tgbotapi.NewMessage(chatID, "неуспешная авторизация"))
		if err != nil {
			log.Error(err)
		}
	default:
		log.Error(err)

		if chatID == 0 {
			return
		}

		_, err := tgbot.Send(tgbotapi.NewMessage(chatID, "что-то пошло не так"))
		if err != nil {
			log.Error(err)
		}
	}
}
