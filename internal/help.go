package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type helpCommand struct {
}

func (cmd *helpCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	text := `Это бот для ЖК Фонвизинский в Москве. На данный момент позволяет узнавать количество работающих лифтов.
Для начала добавьте хотя бы одну квартиру командой /add_flat. По номеру квартиры определяется ваш подъезд.
После этого командой /get_lift_state вы можете узнать, сколько лифтов работает в вашем подъезде. Команда /set_lift_state позволяет обновить эту информацию для соседей.
`
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	return msg, false
}

type helpCommandCreator struct {
}

func (cc *helpCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "help",
		Description: "получить краткую справку",
	}
}

func (cc *helpCommandCreator) Create() Command {
	return &helpCommand{}
}

func newHelpCommandCreator() CommandCreator {
	return &helpCommandCreator{}
}
