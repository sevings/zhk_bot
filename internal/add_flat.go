package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type addFlatCommand struct {
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

	floor := getFloor(flat)
	building := getBuilding(flat)
	text := fmt.Sprintf("Корпус %d, этаж %d, квартира %d. Сохранено.", building, floor, flat)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type addFlatCommandCreator struct {
}

func (cc *addFlatCommandCreator) Text() string {
	return "add_flat"
}

func (cc *addFlatCommandCreator) Create() Command {
	return &addFlatCommand{}
}

func newAddFlatCommandCreator() CommandCreator {
	return &addFlatCommandCreator{}
}
