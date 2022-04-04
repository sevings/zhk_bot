package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	goconf "github.com/zpatrick/go-config"
	"log"
	"strconv"
	"strings"
)

const unrecognisedText = "Неизвестная команда. Попробуйте /help."

type Command interface {
	Exec(update *tgbotapi.Update) (tgbotapi.MessageConfig, bool)
}

type CommandCreator interface {
	Text() string
	Create() Command
}

type zhkBot struct {
	cfg      *goconf.Config
	api      *tgbotapi.BotAPI
	db       *botDB
	admins   []int64
	creators map[string]CommandCreator
	cmds     map[int64]Command
	stop     chan interface{}
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
		cfg:      loadConfig(),
		stop:     make(chan interface{}),
		creators: make(map[string]CommandCreator),
		cmds:     make(map[int64]Command),
	}

	db, err := openBotDB(bot.configString("database.source"))
	if err != nil {
		log.Println(err)
		return bot
	}
	bot.db = db

	bot.admins = bot.configInt64s("telegram.admins")

	bot.addCreator(newHelpCommandCreator())
	bot.addCreator(newAddFlatCommandCreator(db))
	bot.addCreator(newRmFlatCommandCreator(db))
	bot.addCreator(newSetLiftStateCommandCreator(db))
	bot.addCreator(newGetLiftStateCommandCreator(db))

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

func (bot *zhkBot) addCreator(cc CommandCreator) {
	bot.creators[cc.Text()] = cc
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
			if upd.Message == nil {
				continue
			}

			bot.handleMessage(upd)
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

func (bot *zhkBot) handleMessage(upd tgbotapi.Update) {
	chatID := upd.Message.Chat.ID
	var command Command

	if upd.Message.IsCommand() {
		cmd := upd.Message.Command()
		log.Printf("%s sent a command %s", upd.Message.From.UserName, cmd)

		cc := bot.creators[cmd]
		if cc == nil {
			bot.replyText(upd, unrecognisedText)
			return
		}

		_, hasOld := bot.cmds[chatID]
		if hasOld {
			msg := tgbotapi.NewMessage(chatID, "Отмена.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			bot.sendMessage(msg)
		}

		command = cc.Create()
	} else {
		command = bot.cmds[chatID]
		if command == nil {
			bot.replyText(upd, unrecognisedText)
			return
		}
	}

	msg, hasNext := command.Exec(&upd)
	if hasNext {
		bot.cmds[chatID] = command
	} else {
		delete(bot.cmds, chatID)
	}

	bot.sendMessage(msg)
}

func (bot *zhkBot) replyText(upd tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	bot.sendMessage(msg)
}

func (bot *zhkBot) sendMessage(msg tgbotapi.MessageConfig) {
	_, err := bot.api.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
