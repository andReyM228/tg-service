package tg_handlers

import (
	"fmt"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"tg_service/internal/handler/car"
	"tg_service/internal/handler/user"
)

type Handler struct {
	tgbot       *tgbotapi.BotAPI
	userHandler user.Handler
	carHandler  car.Handler
	errChan     chan errs.TgError
	chatGPT     gpt3.ChatGPT
}

func NewHandler(tgbot *tgbotapi.BotAPI, userHandler user.Handler, carHandler car.Handler, errChan chan errs.TgError, chatGPT gpt3.ChatGPT) Handler {
	return Handler{
		tgbot:       tgbot,
		userHandler: userHandler,
		carHandler:  carHandler,
		errChan:     errChan,
		chatGPT:     chatGPT,
	}
}

func (h Handler) RegistrationHandler(update tgbotapi.Update) {
	err := h.userHandler.Create(update.Message.Chat.ID, update)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) AllCarsHandler(update tgbotapi.Update) {
	inlineKeyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("BMW", "all_car_data:bmw"),
			tgbotapi.NewInlineKeyboardButtonData("Audi", "all_car_data:audi"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Mercedes", "all_car_data:mercedes"),
			tgbotapi.NewInlineKeyboardButtonData("Toyota", "all_car_data:toyota"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Citroen", "all_car_data:citroen"),
			tgbotapi.NewInlineKeyboardButtonData("Shkoda", "all_car_data:shkoda"),
		},
	}

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите опцию:")
	msg.ReplyMarkup = inlineKeyboardMarkup

	if _, err := h.tgbot.Send(msg); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) GetMyCarsHandler(update tgbotapi.Update, loginUsers map[int64]string) {

	cars, err := h.carHandler.GetUserCars(loginUsers[update.Message.Chat.ID])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	for _, car := range cars.Cars {
		photo, inlineKeyboard, err := h.carHandler.PrepareCars(car, true)
		if err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.Message.Chat.ID,
			}
			continue
		}

		if _, err := h.tgbot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photo)); err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.Message.Chat.ID,
			}
			continue
		}

		message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("имя: %s, модель: %s, цена: %d", car.Name, car.Model, car.Price))

		message.ReplyMarkup = inlineKeyboard

		if _, err := h.tgbot.Send(message); err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.Message.Chat.ID,
			}
			continue
		}
	}
}

func (h Handler) GetUserHandler(update tgbotapi.Update) {
	msg := strings.Split(update.Message.Text, ":")

	id, err := strconv.Atoi(msg[1])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	userResp, err := h.userHandler.Get(int64(id))
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, userResp)); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) GetCarHandler(update tgbotapi.Update, loginUsers map[int64]string) {
	msg := strings.Split(update.Message.Text, ":")

	if len(msg) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	id, err := strconv.Atoi(msg[1])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	carResp, err := h.carHandler.Get(int64(id), loginUsers[update.Message.Chat.ID])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	photo, inlineKeyboard, err := h.carHandler.PrepareCars(carResp, false)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	if _, err := h.tgbot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photo)); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("имя: %s, модель: %s, цена: %d", carResp.Name, carResp.Model, carResp.Price))

	message.ReplyMarkup = inlineKeyboard

	if _, err := h.tgbot.Send(message); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) StartHandler(update tgbotapi.Update) {
	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "введите /registration чтобы зарегестрироваться, и /login чтобы войти")); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) LoginHandler(update tgbotapi.Update) {
	err := h.userHandler.Login(update.Message.Chat.ID, update)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.Message.Chat.ID,
		}

		return
	}
}
