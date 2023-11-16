package menus

import (
	botgolang "github.com/mail-ru-im/bot-golang"
)

var ACTION_GET_STAND = "/getStand"
var ACTION_CHECK_AVAILABILITY = "/checkAvailability"
var ACTION_CHECK_MY_STANDS = "/checkMyStands"
var ACTION_RELEASE = "/release"
var ACTION_ADD_STAND = "/addStand"
var ACTION_TO_QUEUE = "/toQueue"
var ACTION_ROOT_MENU = "/rootMenu"

type Button struct {
	name   string
	action string
}

var BUTTONS = map[string]Button{
	ACTION_GET_STAND:          {name: "Занять песок", action: ACTION_GET_STAND},
	ACTION_CHECK_AVAILABILITY: {name: "Узнать, какие есть пески", action: ACTION_CHECK_AVAILABILITY},
	ACTION_CHECK_MY_STANDS:    {name: "Статус моих бронирований", action: ACTION_CHECK_MY_STANDS},
	ACTION_RELEASE:            {name: "Освободить", action: ACTION_RELEASE},
	ACTION_ADD_STAND:          {name: "Добавить песок", action: ACTION_ADD_STAND},
	ACTION_TO_QUEUE:           {name: "Встать в очередь", action: ACTION_TO_QUEUE},
	ACTION_ROOT_MENU:          {name: "Меню", action: ACTION_ROOT_MENU},
}

func CreateBaseMenu(bot *botgolang.Bot) botgolang.Keyboard {
	MENU := CreateCustomMenu(bot, []Button{BUTTONS[ACTION_GET_STAND], BUTTONS[ACTION_CHECK_AVAILABILITY], BUTTONS[ACTION_CHECK_MY_STANDS], BUTTONS[ACTION_RELEASE]})
	MENU.DeleteRow(4)
	return MENU
}

func CreateCustomMenu(bot *botgolang.Bot, buttons []Button) botgolang.Keyboard {
	MENU := botgolang.NewKeyboard()
	for index := range buttons {
		MENU.AddRow()
		MENU.AddButton(index, botgolang.NewCallbackButton(buttons[index].name, buttons[index].action))
	}

	MENU.AddRow()
	MENU.AddButton(len(buttons), botgolang.NewCallbackButton(BUTTONS[ACTION_ROOT_MENU].name, BUTTONS[ACTION_ROOT_MENU].action))

	return MENU
}
