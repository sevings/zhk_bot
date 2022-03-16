package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

var liftStates []string
var liftStateMarkup tgbotapi.ReplyKeyboardMarkup

func init() {
	liftStates = []string{"Неизвестно", "Работает", "Авария", "На обслуживании", "Пожар"}

	liftStateMarkup = tgbotapi.NewReplyKeyboard()
	liftStateMarkup.OneTimeKeyboard = true
	for _, state := range liftStates {
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(state))
		liftStateMarkup.Keyboard = append(liftStateMarkup.Keyboard, row)
	}
}

type setLeftStateCommand struct {
	largeFreightElevatorState   int
	smallFreightElevatorState   int
	smallPassengerElevatorState int
}

func (cmd *setLeftStateCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	building := getBuilding(683)

	if upd.Message.IsCommand() {
		return cmd.setStateMessage(upd, building, "большого грузового")
	}

	stateStr := upd.Message.Text

	state := -1
	for i, s := range liftStates {
		if s == stateStr {
			state = i
			break
		}
	}

	if state < 0 {
		text := fmt.Sprintf("Неверный статус.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	if cmd.largeFreightElevatorState < 0 {
		cmd.largeFreightElevatorState = state
		return cmd.setStateMessage(upd, building, "малого грузового")
	}

	if cmd.smallFreightElevatorState < 0 {
		cmd.smallFreightElevatorState = state
		return cmd.setStateMessage(upd, building, "малого пассажирского")
	}

	cmd.smallPassengerElevatorState = state

	text := fmt.Sprintf(`Установлены статусы лифтов.
Большой грузовой лифт: %s.
Малый грузовой лифт: %s.
Малый пассажирский лифт: %s.
`,
		strings.ToLower(liftStates[cmd.largeFreightElevatorState]),
		strings.ToLower(liftStates[cmd.smallFreightElevatorState]),
		strings.ToLower(liftStates[cmd.smallPassengerElevatorState]),
	)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

func (cmd *setLeftStateCommand) setStateMessage(upd *tgbotapi.Update, building int, lift string) (tgbotapi.MessageConfig, bool) {
	text := fmt.Sprintf("Выберите статус <b>%s</b> лифта в %d корпусе.", lift, building)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = liftStateMarkup
	return msg, true
}

type setLeftStateCommandCreator struct {
}

func (cc *setLeftStateCommandCreator) Text() string {
	return "set_lift_state"
}

func (cc *setLeftStateCommandCreator) Create() Command {
	return &setLeftStateCommand{-1, -1, -1}
}

func newSetLiftStateCommandCreator() CommandCreator {
	return &setLeftStateCommandCreator{}
}
