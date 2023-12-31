package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	botgolang "github.com/mail-ru-im/bot-golang"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vk.com/standsbot/menus"
)

type User struct {
	gorm.Model
	ID           string `gorm:"primarykey"`
	First_name   string
	Last_name    string
	Last_command string
}

type Stand struct {
	gorm.Model
	Name        string
	Description string
}

type Booking struct {
	gorm.Model
	UserID  string
	StandID uint
	User    User
	Stand   Stand
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println("You have to create .env file")
		return
	}

	db, dbErr := gorm.Open(sqlite.Open("standsbot.sqlite"), &gorm.Config{})
	if dbErr != nil {
		log.Panicln("Error in DB", dbErr)
		return
	}
	db.AutoMigrate(&Stand{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Booking{})

	var TOKEN = os.Getenv("TOKEN")
	var API_URL = os.Getenv("API_URL")

	bot, err := botgolang.NewBot(TOKEN, botgolang.BotApiURL(API_URL), botgolang.BotDebug(true))
	if err != nil {
		log.Println("wrong token", err)
		return
	}
	client := *http.DefaultClient
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logrus.DebugLevel)
	buttonResponseClient := botgolang.NewCustomClient(&client, API_URL, TOKEN, logger)

	ctx, finish := context.WithCancel(context.Background())
	defer finish()

	updates := bot.GetUpdatesChannel(ctx)
	for update := range updates {
		sn := update.Payload.BaseEventPayload.From.ID

		fmt.Println(update.Type)
		switch update.Type {
		case "newMessage":
			firstName := update.Payload.BaseEventPayload.From.FirstName
			lastName := update.Payload.BaseEventPayload.From.LastName
			var user User
			db.Where(User{ID: sn}).Attrs(User{ID: sn, First_name: firstName, Last_name: lastName}).FirstOrCreate(&user)

			switch update.Payload.Text {
			default:
				buttonMenu := bot.NewInlineKeyboardMessage(sn, "Привет, "+firstName+", чего изволите?", menus.CreateBaseMenu(bot))
				buttonMenu.Send()
			}
		case "callbackQuery":
			switch update.Payload.CallbackData {
			case menus.ACTION_CHECK_MY_STANDS:
				var stands []Stand
				text := ""
				db.Where(Booking{User: User{ID: sn}}).Find(&stands)
				for _, stand := range stands {
					text += "\r\n" + stand.Name
				}
				if len(stands) == 0 {
					// text = "У вас нет забронированных песков" + MENU
				}
				message := bot.NewTextMessage(sn, text)
				message.Send()

			case menus.ACTION_GET_STAND:
				var stands []Stand
				var bookings []Booking
				text := "Сейчас свободны следующие:"
				kb := botgolang.NewKeyboard()
				db.Find(&stands)
				db.Find(&bookings)
				bookingsMap := make(map[uint]bool)
				for _, book := range bookings {
					bookingsMap[book.StandID] = true
				}
				for _, stand := range stands {
					if bookingsMap[stand.ID] {
						text += "\r\n" + stand.Name + " - " + stand.Description
						kb.AddRow()
						kb.AddButton(kb.RowsCount()-1, botgolang.NewCallbackButton("Занять "+stand.Name, "GET_"+strconv.FormatUint(uint64(stand.ID), 10)))
					}
				}
				if len(stands) == 0 {
					text = "Сейчас нет свободных песков :("
					kb = menus.CreateCustomMenu(bot, []menus.Button{menus.BUTTONS[menus.ACTION_TO_QUEUE]})
				}
				responseError := buttonResponseClient.SendAnswerCallbackQuery(&botgolang.ButtonResponse{
					QueryID:   update.Payload.QueryID,
					Text:      "",
					URL:       "",
					ShowAlert: false,
				})
				message := bot.NewInlineKeyboardMessage(sn, text, kb)
				message.Send()
				if responseError != nil {
					fmt.Println(responseError)
				}
			}
		}
	}
}
