package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	botgolang "github.com/mail-ru-im/bot-golang"
)

type userModel struct {
	sn         string
	first_name string
	last_name  string
}

var SQL_SELECT_USER = "SELECT sn, first_name, last_name FROM users WHERE sn=?"
var SQL_CREATE_USER = "INSERT INTO users (first_name, last_name, sn) VALUES(?, ?, ?)"

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println("You have to create .env file")
		return
	}

	db, dbErr := sql.Open("sqlite3", "./standsbot.sqlite")
	if dbErr != nil {
		log.Panicln("Error in DB", dbErr)
		return
	}
	db.SetMaxOpenConns(1)
	userSelectStmt, _ := db.Prepare(SQL_SELECT_USER)
	defer userSelectStmt.Close()
	userInsertStmt, _ := db.Prepare(SQL_CREATE_USER)
	defer userInsertStmt.Close()
	defer db.Close()

	var TOKEN = os.Getenv("TOKEN")
	var API_URL = os.Getenv("API_URL")

	bot, err := botgolang.NewBot(TOKEN, botgolang.BotApiURL(API_URL), botgolang.BotDebug(true))
	if err != nil {
		log.Println("wrong token", err)
		return
	}

	ctx, finish := context.WithCancel(context.Background())
	defer finish()
	updates := bot.GetUpdatesChannel(ctx)
	for update := range updates {
		if update.Type == "newMessage" {
			sn := update.Payload.BaseEventPayload.From.ID
			users, selectErr := userSelectStmt.Exec(sn)
			if selectErr != nil {
				log.Println(selectErr)
			}
			log.Println(users)
		}
	}

}
