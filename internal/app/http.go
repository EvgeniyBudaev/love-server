package app

import (
	profileRepo "github.com/EvgeniyBudaev/love-server/internal/adapter/psqlRepo/profile"
	identityEntity "github.com/EvgeniyBudaev/love-server/internal/entity/identity"
	profileHandler "github.com/EvgeniyBudaev/love-server/internal/handler/profile"
	userHandler "github.com/EvgeniyBudaev/love-server/internal/handler/user"
	"github.com/EvgeniyBudaev/love-server/internal/middlewares"
	profileUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/profile"
	userUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

var prefix = "/api/v1"

const (
	EMOJI_COIN       = "\U0001FA99"
	EMOJI_SMILE      = "\U0001F642"
	EMOJI_SUNGLASSES = "\U0001F60E"
)

var bot *tgbotapi.BotAPI
var err error

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func delay(seconds uint8) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func printSystemMessageWithDelay(chatId int64, delayInSec uint8, message string) {
	bot.Send(tgbotapi.NewMessage(chatId, message))
	delay(delayInSec)
}

func printIntro(chatId int64) {
	printSystemMessageWithDelay(chatId, 1, "Привет! "+EMOJI_SUNGLASSES)
	printSystemMessageWithDelay(chatId, 5, "Нажми на кнопку App,"+
		" чтобы перейти на главную страницу приложения")
}

func (app *App) StartHTTPServer() error {
	// Telegram Bot
	//if bot, err = tgbotapi.NewBotAPI(app.config.TelegramBotToken); err != nil {
	//	log.Panic(err)
	//}
	//bot.Debug = true
	//log.Printf("Authorized on account %s", bot.Self.UserName)
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//updates := bot.GetUpdatesChan(u)
	//for update := range updates {
	//	chatId := update.Message.Chat.ID
	//	if isStartMessage(&update) {
	//		log.Printf("Начало общения: [%s] %s", update.Message.From.UserName, update.Message.Text)
	//		printIntro(chatId)
	//	}
	//	//if update.Message != nil {
	//	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//	//	msg := tgbotapi.NewMessage(chatId, update.Message.Text)
	//	//	msg.ReplyToMessageID = update.Message.MessageID
	//	//	bot.Send(msg)
	//	//}
	//}

	app.fiber.Static("/static", "./static")
	im := identityEntity.NewIdentity(app.config, app.Logger)
	pr := profileRepo.NewRepositoryProfile(app.Logger, app.db.psql)
	imc := userUseCase.NewUseCaseUser(app.Logger, im)
	puc := profileUseCase.NewUseCaseProfile(app.Logger, pr)
	imh := userHandler.NewHandlerUser(app.Logger, imc)
	ph := profileHandler.NewHandlerProfile(app.Logger, puc)
	grp := app.fiber.Group(prefix)
	middlewares.InitFiberMiddlewares(
		app.fiber, app.config, app.Logger, grp, imh, ph, InitPublicRoutes, InitProtectedRoutes)
	if err := app.fiber.Listen(app.config.Port); err != nil {
		app.Logger.Fatal("error func StartHTTPServer, method Listen by path internal/app/http.go", zap.Error(err))
	}
	return nil
}
