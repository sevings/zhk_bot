package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	goconf "github.com/zpatrick/go-config"
	"log"
	"strconv"
	"strings"
)

const unrecognisedText = "Неизвестная команда. Попробуйте /help."

type zhkBot struct {
	cfg    *goconf.Config
	api    *tgbotapi.BotAPI
	admins []int64
	stop   chan interface{}
}

func loadConfig() *goconf.Config {
	toml := goconf.NewTOMLFile("configs/zhk_bot.toml")
	loader := goconf.NewOnceLoader(toml)
	config := goconf.NewConfig([]goconf.Provider{loader})
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	return config
}

func NewBot() *zhkBot {
	bot := &zhkBot{
		cfg:  loadConfig(),
		stop: make(chan interface{}),
	}

	bot.admins = bot.configInt64s("telegram.admins")

	return bot
}

func (bot *zhkBot) configString(field string) string {
	value, err := bot.cfg.String(field)
	if err != nil {
		log.Println(err)
	}

	return value
}

func (bot *zhkBot) configStrings(field string) []string {
	value, err := bot.cfg.String(field)
	if err != nil {
		log.Println(err)
	}

	return strings.Split(value, ";")
}

func (bot *zhkBot) configInt64s(field string) []int64 {
	var values []int64

	for _, str := range bot.configStrings(field) {
		value, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			values = append(values, value)
		} else {
			log.Println(err)
		}
	}

	return values
}

func (bot *zhkBot) Run() {
	token := bot.configString("telegram.token")
	if len(token) == 0 {
		return
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}

	bot.api = api

	log.Printf("Running Telegram bot %s\n", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case <-bot.stop:
			return
		case upd := <-updates:
			if upd.Message == nil || !upd.Message.IsCommand() {
				bot.sendMessageNow(upd.Message.Chat.ID, unrecognisedText)
				continue
			}

			bot.command(upd)
		}
	}
}

func (bot *zhkBot) Stop() {
	if bot.api == nil {
		return
	}

	bot.api.StopReceivingUpdates()
	close(bot.stop)
}

func (bot *zhkBot) command(upd tgbotapi.Update) {
	cmd := upd.Message.Command()
	log.Printf("%s sent a command %s", upd.Message.From.UserName, cmd)

	var reply string
	switch cmd {
	case "help":
		reply = bot.help(&upd)
	default:
		reply = unrecognisedText
	}

	bot.sendMessageNow(upd.Message.Chat.ID, reply)
}

func (bot *zhkBot) sendMessageNow(chat int64, text string) {
	msg := tgbotapi.NewMessage(chat, text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	_, err := bot.api.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (bot *zhkBot) help(_ *tgbotapi.Update) string {
	text := `(помощь по командам бота)`

	return text
}
