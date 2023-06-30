package app

import (
	"context"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/gpt3"
	"github.com/andReyM228/lib/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	stdLog "log"
	"net/http"
	"os"
	"strings"
	"tg_service/internal/config"
	car_handler "tg_service/internal/handler/car"
	user_handler "tg_service/internal/handler/user"
	"tg_service/internal/repository/cars"
	"tg_service/internal/repository/user"
	"tg_service/internal/service/car"
	user_service "tg_service/internal/service/user"
	"tg_service/internal/tg_handlers"
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
	tgHandler   tg_handlers.Handler
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
					go a.tgHandler.BuyDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "view_data"):
					go a.tgHandler.ViewDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "characteristics_data"):
					go a.tgHandler.CharacteristicsDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "all_car_data"):
					go a.tgHandler.AllCarDataButton(update, a.loginUsers)
					continue

				}
			}

			continue
		}

		switch {
		case strings.Contains(update.Message.Text, "/start"):
			go a.tgHandler.StartHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/get-car"):
			go a.tgHandler.GetCarHandler(update, a.loginUsers)
			continue

		case strings.Contains(update.Message.Text, "/all-cars"):
			go a.tgHandler.AllCarsHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/get-user"):
			go a.tgHandler.GetUserHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/registration"):
			go a.tgHandler.RegistrationHandler(update, updates)
			continue

		case strings.Contains(update.Message.Text, "/login"):
			go a.tgHandler.LoginHandler(update, updates)
			continue

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
	a.tgHandler = tg_handlers.NewHandler(a.tgbot, a.userHandler, a.carHandler, a.errChan, a.chatGPT)

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
