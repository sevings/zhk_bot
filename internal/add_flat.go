package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type addFlatCommand struct {
	db *botDB
}

func (cmd *addFlatCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	if upd.Message.IsCommand() {
		text := fmt.Sprintf("Введите номер вашей квартиры (%d—%d).", minFlat, maxFlat)
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, true
	}

	n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
	flat := int(n)
	if err != nil || flat < minFlat || flat > maxFlat {
		text := fmt.Sprintf("Неверный номер.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	cmd.db.addUserFlat(upd.Message.Chat.ID, upd.Message.Chat.UserName, flat)

	floor := getFloor(flat)
	building := getBuilding(flat)
	text := fmt.Sprintf("Корпус %d, этаж %d, квартира %d. Сохранено.", building, floor, flat)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type addFlatCommandCreator struct {
	db *botDB
}

func (cc *addFlatCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "add_flat",
		Description: "добавить квартиру проживания",
	}
}

func (cc *addFlatCommandCreator) Create() Command {
	return &addFlatCommand{db: cc.db}
}

func newAddFlatCommandCreator(db *botDB) CommandCreator {
	return &addFlatCommandCreator{db: db}
}
