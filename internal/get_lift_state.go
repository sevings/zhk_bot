package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

type getLeftStateCommand struct {
}

func (cmd *getLeftStateCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	building := getBuilding(683)
	n := 2

	text := fmt.Sprintf("В %d корпусе ", building)
	switch n {
	case 0:
		text += "не работает <b>ни один</b> лифта."
	case 1:
		text += "работает <b>один</b> лифт."
	case 2:
		text += "работают <b>два</b> лифта."
	case 3:
		text += "работают <b>все</b> лифты."
	}

	at := time.Now().Add(-24 * time.Hour)
	text += fmt.Sprintf(" Обновлено %s.", formatDate(at))

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	return msg, false
}

type getLeftStateCommandCreator struct {
}

func (cc *getLeftStateCommandCreator) Text() string {
	return "get_lift_state"
}

func (cc *getLeftStateCommandCreator) Create() Command {
	return &getLeftStateCommand{}
}

func newGetLiftStateCommandCreator() CommandCreator {
	return &getLeftStateCommandCreator{}
}

func formatDate(date time.Time) string {
	today := time.Now()

	if today.YearDay() == date.YearDay() && today.Year() == date.Year() {
		return "сегодня в " + date.Format("15:04")
	}

	yesterday := today.Add(-24 * time.Hour)
	if yesterday.YearDay() == date.YearDay() && yesterday.Year() == date.Year() {
		return "вчера в " + date.Format("15:04")
	}

	var str = date.Format("_2")

	switch date.Month() {
	case 0:
		str += " января"
	case 1:
		str += " февраля"
	case 2:
		str += " марта"
	case 3:
		str += " апреля"
	case 4:
		str += " мая"
	case 5:
		str += " июня"
	case 6:
		str += " июля"
	case 7:
		str += " августа"
	case 8:
		str += " сентября"
	case 9:
		str += " октября"
	case 10:
		str += " ноября"
	case 11:
		str += " декабря"
	}

	if today.Year() != date.Year() {
		str += " " + date.Format("2006")
	}

	return str
}
