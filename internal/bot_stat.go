package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type botStatCommand struct {
	db *botDB
}

func (cmd *botStatCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	bs := cmd.db.getBotStat()

	text := "Статистика бота\n"
	text += "\n<b>Пользователей</b>: " + strconv.FormatInt(bs.users, 10)
	text += "\n<b>Квартир</b>: " + strconv.FormatInt(bs.flats, 10)
	text += "\n<b>Обновлений состояния лифтов</b>: " + strconv.FormatInt(bs.liftStates, 10)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	return msg, false
}

type botStatCommandCreator struct {
	db *botDB
}

func (cc *botStatCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "bot_stat",
		Description: "посмотреть статистику бота",
	}
}

func (cc *botStatCommandCreator) Create() Command {
	return &botStatCommand{db: cc.db}
}

func newBotStatCommandCreator(db *botDB) CommandCreator {
	return &botStatCommandCreator{db: db}
}
