package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	car_handler "tg_service/internal/handler/car"
	user_handler "tg_service/internal/handler/user"
	"tg_service/internal/repository/cars"
	"tg_service/internal/repository/user"
	"tg_service/internal/service/car"
	user_service "tg_service/internal/service/user"

	"tg_service/internal/config"

	"github.com/sirupsen/logrus"
)

type App struct {
	config      config.Config
	serviceName string
	tgbot       *tgbotapi.BotAPI
	logger      *logrus.Logger
	usersRepo   user.Repository
	carsRepo    cars.Repository
	userService user_service.Service
	carService  car.Service
	userHandler user_handler.Handler
	carHandler  car_handler.Handler
	clientHTTP  *http.Client
}

func New(name string) App {
	return App{
		serviceName: name,
	}
}

func (a *App) Run() {
	a.populateConfig()
	a.initLogger()
	a.initHTTPClient()
	a.initRepos()
	a.initServices()
	a.initHandlers()
	a.initTgBot()
}

func (a *App) initTgBot() {
	var err error
	a.tgbot, err = tgbotapi.NewBotAPI(a.config.TgBot.Token)
	if err != nil {
		log.Fatal(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 1
	updates := a.tgbot.GetUpdatesChan(updateConfig)

	a.logger.Debug("tg_bot api started")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch {
		case strings.Contains(update.Message.Text, "/start"):
			if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "tg-bot started")); err != nil {
				log.Fatal(err)
			}

		case strings.Contains(update.Message.Text, "/get-car"):
			msg := strings.Split(update.Message.Text, ":")

			id, err := strconv.Atoi(msg[1])
			if err != nil {
				log.Fatal(err)
			}

			carResp, err := a.carHandler.Get(int64(id))
			if err != nil {
				log.Fatal(err)
			}

			if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, carResp)); err != nil {
				log.Fatal(err)
			}

		case strings.Contains(update.Message.Text, "/get-user"):
			msg := strings.Split(update.Message.Text, ":")

			id, err := strconv.Atoi(msg[1])
			if err != nil {
				log.Fatal(err)
			}

			userResp, err := a.userHandler.Get(int64(id))
			if err != nil {
				log.Fatal(err)
			}

			if _, err := a.tgbot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, userResp)); err != nil {
				log.Fatal(err)
			}
		}

	}

}

func (a *App) initLogger() {
	a.logger = logrus.New()
	a.logger.SetLevel(logrus.DebugLevel)
}

func (a *App) initRepos() {
	a.carsRepo = cars.NewRepository(a.logger, a.clientHTTP)
	a.usersRepo = user.NewRepository(a.logger, a.clientHTTP)

	a.logger.Debug("repos created")
}

func (a *App) initServices() {
	a.carService = car.NewService(a.carsRepo, a.logger)
	a.userService = user_service.NewService(a.usersRepo, a.logger)

	a.logger.Debug("repos created")
}

func (a *App) initHandlers() {
	a.carHandler = car_handler.NewHandler(a.carService)
	a.userHandler = user_handler.NewHandler(a.userService)
	a.logger.Debug("handlers created")
}

func (a *App) populateConfig() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal()
	}

	a.config = cfg
}

func (a *App) initHTTPClient() {
	a.clientHTTP = http.DefaultClient
}
