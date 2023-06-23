package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/gpt3"
	"github.com/andReyM228/lib/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	stdLog "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tg_service/internal/config"
	"tg_service/internal/domain"
	car_handler "tg_service/internal/handler/car"
	user_handler "tg_service/internal/handler/user"
	"tg_service/internal/repository/cars"
	"tg_service/internal/repository/user"
	"tg_service/internal/service/car"
	user_service "tg_service/internal/service/user"
)

const (
	characteristicsAnswerBody = `
Engine: %s,	
Drive Type: %s,	
Power: %s,	
Acceleration: %s,	
Top Speed: %s,	
Fuel Economy: %s,	
Transmission: %s,
						`
	characteristicsRequest = "опиши мне главные характеристики машины %s в виде одного json на английском, с полями: engine, power, acceleration, top_speed, fuel_economy, transmission, drive_type"
)

type App struct {
	config      config.Config
	serviceName string
	tgbot       *tgbotapi.BotAPI
	logger      log.Logger
	usersRepo   user.Repository
	carsRepo    cars.Repository
	userService user_service.Service
	carService  car.Service
	userHandler user_handler.Handler
	carHandler  car_handler.Handler
	clientHTTP  *http.Client
	errChan     chan errs.TgError
	loginUsers  map[int64]string
	chatGPT     gpt3.ChatGPT
}

func New(name string) App {
	return App{
		serviceName: name,
	}
}

func (a *App) initGPT() {
	a.chatGPT = gpt3.Init(a.config.ChatGPT.Key, a.config.ChatGPT.Model)
}

func (a *App) Run(ctx context.Context) {
	a.populateConfig()
	a.initLogger()
	a.initGPT()
	a.listenErrs(ctx)
	a.initTgBot()
	a.initHTTPClient()
	a.initRepos()
	a.initServices()
	a.initHandlers()
	a.listenTgBot()
}

func (a *App) listenErrs(ctx context.Context) {
	a.errChan = make(chan errs.TgError)

	go func() {
		for {
			select {
			case err := <-a.errChan:
				go func(err errs.TgError) {
					errs.HandleError(err.Err, a.logger, a.tgbot, err.ChatID)
				}(err)
			case <-ctx.Done():
				a.logger.Debug("ctx is done")
				os.Exit(1)

			}
		}
	}()
}

func (a *App) initTgBot() {
	var err error
	a.tgbot, err = tgbotapi.NewBotAPI(a.config.TgBot.Token)
	if err != nil {
		a.errChan <- errs.TgError{
			Err: err,
		}
		return
	}

}

func (a *App) listenTgBot() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 1
	updates := a.tgbot.GetUpdatesChan(updateConfig)

	a.logger.Debug("tg_bot api started")

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				switch {
				case strings.Contains(update.CallbackQuery.Data, "buy_data"):
					data := strings.Split(update.CallbackQuery.Data, ":")
					if len(data) < 2 {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

					carID, err := strconv.Atoi(data[1])
					if err != nil {
						a.errChan <- errs.TgError{
							Err:    err,
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

					err = a.carHandler.BuyCar(update.CallbackQuery.Message.Chat.ID, int64(carID))
					if err != nil {
						a.errChan <- errs.TgError{
							Err:    err,
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

					if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "congratulations!, you bought a car")); err != nil {
						a.errChan <- errs.TgError{
							Err:    err,
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

				case strings.Contains(update.CallbackQuery.Data, "view_data"):
					data := strings.Split(update.CallbackQuery.Data, ":")
					if len(data) < 2 {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

					answ, err := a.chatGPT.GetCompletion(fmt.Sprintf("расскажи мне об этой машине: %s", data[1]))
					if err != nil {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
					}

					if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answ)); err != nil {
						a.errChan <- errs.TgError{
							Err:    err,
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

				case strings.Contains(update.CallbackQuery.Data, "characteristics_data"):
					data := strings.Split(update.CallbackQuery.Data, ":")
					if len(data) < 2 {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}

					answ, err := a.chatGPT.GetCompletion(fmt.Sprintf(characteristicsRequest, data[1]))
					if err != nil {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
					}

					var result domain.CarCharacteristics

					err = json.Unmarshal([]byte(answ), &result)
					if err != nil {
						a.errChan <- errs.TgError{
							Err:    errs.BadRequestError{},
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
					}

					answer := fmt.Sprintf(characteristicsAnswerBody, result.Engine, result.DriveType, result.Power, result.Acceleration, result.TopSpeed, result.FuelEconomy, result.Transmission)

					if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)); err != nil {
						a.errChan <- errs.TgError{
							Err:    err,
							ChatID: update.CallbackQuery.Message.Chat.ID,
						}
						continue
					}
				}
			}

			continue
		}

		switch {
		case strings.Contains(update.Message.Text, "/start"):
			if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "введите /registration чтобы зарегестрироваться, и /login чтобы войти")); err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

		case strings.Contains(update.Message.Text, "/get-car"):
			msg := strings.Split(update.Message.Text, ":")

			id, err := strconv.Atoi(msg[1])
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			carResp, carImage, err := a.carHandler.Get(int64(id), a.loginUsers[update.Message.Chat.ID])
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			imageBytes, err := base64.StdEncoding.DecodeString(carImage)
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			photo := tgbotapi.FileBytes{Name: "image.jpg", Bytes: imageBytes}

			if _, err := a.tgbot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photo)); err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			buyButton := tgbotapi.NewInlineKeyboardButtonData("buy", fmt.Sprintf("buy_data:%s", strconv.Itoa(id)))
			viewButton := tgbotapi.NewInlineKeyboardButtonData("view", fmt.Sprintf("view_data:%s %s", carResp.Name, carResp.Model))
			characteristicsButton := tgbotapi.NewInlineKeyboardButtonData("characteristics", fmt.Sprintf("characteristics_data:%s %s", carResp.Name, carResp.Model))

			row := tgbotapi.NewInlineKeyboardRow(buyButton, viewButton, characteristicsButton)

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

			message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("имя: %s, модель: %s, цена: %d", carResp.Name, carResp.Model, carResp.Price))

			message.ReplyMarkup = inlineKeyboard

			if _, err := a.tgbot.Send(message); err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

		case strings.Contains(update.Message.Text, "/get-cars"):
			msg := strings.Split(update.Message.Text, ":")

			carsResp, err := a.carHandler.GetAll(a.loginUsers[update.Message.Chat.ID])
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			var buttons []tgbotapi.InlineKeyboardButton
			for _, c := range carsResp.Cars {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("buy", fmt.Sprintf("get_full_info_data:%d", c.ID)))
			}

			row := tgbotapi.NewInlineKeyboardRow(buttons...)

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

			message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("имя: %s, модель: %s, цена: %d", carResp.Name, carResp.Model, carResp.Price))

			message.ReplyMarkup = inlineKeyboard

			if _, err := a.tgbot.Send(message); err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

		case strings.Contains(update.Message.Text, "/get-user"):
			msg := strings.Split(update.Message.Text, ":")

			id, err := strconv.Atoi(msg[1])
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			userResp, err := a.userHandler.Get(int64(id))
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

			if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, userResp)); err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}

		case strings.Contains(update.Message.Text, "/registration"):
			err := a.userHandler.Create(updates, update.Message.Chat.ID)
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}
				continue
			}
		case strings.Contains(update.Message.Text, "/login"):
			err := a.userHandler.Login(updates, update.Message.Chat.ID)
			if err != nil {
				a.errChan <- errs.TgError{
					Err:    err,
					ChatID: update.Message.Chat.ID,
				}

				continue
			}

			a.logger.Debugf("map: %v", a.loginUsers)

		}
	}
}

func (a *App) initLogger() {
	a.logger = log.Init()
}

func (a *App) initRepos() {
	a.carsRepo = cars.NewRepository(a.logger, a.clientHTTP)
	a.usersRepo = user.NewRepository(a.logger, a.clientHTTP)

	a.logger.Debug("repos created")
}

func (a *App) initServices() {
	a.carService = car.NewService(a.carsRepo, a.logger)
	a.userService = user_service.NewService(a.usersRepo, a.logger)

	a.logger.Debug("services created")
}

func (a *App) initHandlers() {
	a.loginUsers = map[int64]string{}
	a.carHandler = car_handler.NewHandler(a.carService, a.tgbot)
	a.userHandler = user_handler.NewHandler(a.userService, a.tgbot, a.loginUsers)
	a.logger.Debug("handlers created")
}

func (a *App) populateConfig() {
	cfg, err := config.ParseConfig()
	if err != nil {
		stdLog.Fatal()
	}

	a.config = cfg
}

func (a *App) initHTTPClient() {
	a.clientHTTP = http.DefaultClient
}
