package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goexpenseslog/internal/repos"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	var Configuration *goconfig.Configuration
	var Bot *TelegramBot

	var RawMessagesRepository *repos.RawMessageRepository
	var ExpensesRepository *repos.ExpensesRepository

	log.Println("Starting bot service")

	Configuration = goconfig.NewConfiguration(
		&goconfig.Options{
			Source: "../config/config.yml",
			Format: "yaml",
		},
	)

	Configuration.Set("telegram.token", "")
	Configuration.Set("control.userId", 0)
	Configuration.Set("db.path", "../data/data.db")
	Configuration.Set("proxy.socks5", "")

	Configuration.Apply()

	var db *sql.DB
	if db, err = sql.Open("sqlite3", Configuration.Get("db.path").String()); err != nil {
		log.Fatalf("open database error: %v", err)
	}
	defer db.Close()

	RawMessagesRepository = repos.NewRawMessgeRepository(db)
	ExpensesRepository = repos.NewExpensesRepository(db)

	botService := NewBot(RawMessagesRepository, ExpensesRepository, Configuration)

	Bot = New(Configuration.Get("telegram.token").String(), Configuration)

	Bot.SetCommandHandler("start", botService.HandleStart)
	Bot.SetCommandHandler("help", botService.HandleStart)
	Bot.SetCommandHandler("list", botService.HandleUserData)

	Bot.SetCommandHandler("listall", botService.HandleList)

	Bot.SetCommandHandler("sum_this_month", botService.HandleMonthSummary)
	Bot.SetCommandHandler("sum_prev_month", botService.HandlePrevMonthSummary)

	Bot.SetDefaultHandler(botService.HandleDefault)

	log.Fatal(Bot.Listen(context.Background()))
}
