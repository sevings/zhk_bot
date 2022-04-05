package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type helpCommand struct {
}

func (cmd *helpCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	text := `(помощь по командам бота)`
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	return msg, false
}

type helpCommandCreator struct {
}

func (cc *helpCommandCreator) Text() string {
	return "help"
}

func (cc *helpCommandCreator) Create() Command {
	return &helpCommand{}
}

func newHelpCommandCreator() CommandCreator {
	return &helpCommandCreator{}
}
