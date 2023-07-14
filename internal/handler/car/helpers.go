package car

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"tg_service/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h Handler) PrepareCars(car domain.Car, myCar bool) (tgbotapi.FileBytes, tgbotapi.InlineKeyboardMarkup, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(car.Image)
	if err != nil {
		return tgbotapi.FileBytes{}, tgbotapi.InlineKeyboardMarkup{}, err
	}

	photo := tgbotapi.FileBytes{Name: "image.jpg", Bytes: imageBytes}

	var buttons []tgbotapi.InlineKeyboardButton

	buyButton := tgbotapi.NewInlineKeyboardButtonData("buy", fmt.Sprintf("buy_data:%s", strconv.Itoa(car.ID)))
	viewButton := tgbotapi.NewInlineKeyboardButtonData("view", fmt.Sprintf("view_data:%s %s", car.Name, car.Model))
	characteristicsButton := tgbotapi.NewInlineKeyboardButtonData("characteristics", fmt.Sprintf("characteristics_data:%s %s", car.Name, car.Model))
	sellButton := tgbotapi.NewInlineKeyboardButtonData("sell", fmt.Sprintf("sell_data:%s", strconv.Itoa(car.ID)))

	buttons = append(buttons, buyButton, viewButton)

	if myCar {
		buttons = append(buttons, sellButton)
	}

	buttons = append(buttons, characteristicsButton)

	row := tgbotapi.NewInlineKeyboardRow(buttons...)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	return photo, inlineKeyboard, nil
}
