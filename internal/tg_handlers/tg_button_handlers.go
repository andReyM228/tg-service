package tg_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/errs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"tg_service/internal/domain"
)

func (h Handler) AllCarDataButton(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")
	if len(data) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	token, err := h.cache.GetToken(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	cars, err := h.carHandler.GetAllCars(token, data[1])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	for _, car := range cars.Cars {
		photo, inlineKeyboard, err := h.carHandler.PrepareCars(car, false)
		if err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.CallbackQuery.Message.Chat.ID,
			}
			continue
		}

		if _, err := h.tgbot.Send(tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, photo)); err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.CallbackQuery.Message.Chat.ID,
			}
			continue
		}

		message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("имя: %s, модель: %s, цена: %d", car.Name, car.Model, car.Price))

		message.ReplyMarkup = inlineKeyboard

		if _, err := h.tgbot.Send(message); err != nil {
			h.errChan <- errs.TgError{
				Err:    err,
				ChatID: update.CallbackQuery.Message.Chat.ID,
			}
			continue
		}
	}
}

func (h Handler) CharacteristicsDataButton(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")
	if len(data) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	answ, err := h.chatGPT.GetCompletion(fmt.Sprintf(characteristicsRequest, data[1]))
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
	}

	var result domain.CarCharacteristics

	err = json.Unmarshal([]byte(answ), &result)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
	}

	answer := fmt.Sprintf(characteristicsAnswerBody, result.Engine, result.DriveType, result.Power, result.Acceleration, result.TopSpeed, result.FuelEconomy, result.Transmission)

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) ViewDataButton(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")
	if len(data) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	answ, err := h.chatGPT.GetCompletion(fmt.Sprintf("расскажи мне об этой машине: %s", data[1]))
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answ)); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}
}

func (h Handler) BuyDataButton(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")
	if len(data) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	token, err := h.cache.GetToken(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	carID, err := strconv.Atoi(data[1])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	car, err := h.carHandler.GetCar(int64(carID), token)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	msg := fmt.Sprintf("you have to pay %d to address %s \n and you should send tx-hash", car.Price, h.config.CarPaymentAddress)
	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg)); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	h.processingBuyUsers.Create(update.CallbackQuery.Message.Chat.ID, int64(carID))
}

func (h Handler) SellDataButton(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")
	if len(data) < 2 {
		h.errChan <- errs.TgError{
			Err:    errs.BadRequestError{},
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	token, err := h.cache.GetToken(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	carID, err := strconv.Atoi(data[1])
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	err = h.carHandler.SellCar(update.CallbackQuery.Message.Chat.ID, int64(carID), token)
	if err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}

	if _, err := h.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "congratulations!, you sell car")); err != nil {
		h.errChan <- errs.TgError{
			Err:    err,
			ChatID: update.CallbackQuery.Message.Chat.ID,
		}
		return
	}
}
