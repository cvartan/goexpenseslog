package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goconfig/reader/yamlreader"
	"github.com/cvartan/goexpenseslog/internal/bot"
	"github.com/cvartan/goexpenseslog/internal/bothandlers"
	"github.com/cvartan/goexpenseslog/internal/repos"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	var configPath string = *flag.String("config", "../config/config.yml", "set path to configuration file")
	var Configuration *goconfig.ConfigurationManager
	var Bot *bot.TelegramBot

	var RawMessagesRepository *repos.RawMessageRepository
	var ExpensesRepository *repos.ExpensesRepository

	log.Println("Starting bot service")

	Configuration = goconfig.New().
		SetSource(configPath).
		SetReader(&yamlreader.YamlConfigurationReader{})

	if err := Configuration.Read(); err != nil {
		log.Fatalf("read configuration error: %v", err)
	}

	var db *sql.DB
	if db, err = sql.Open("sqlite3", Configuration.Get("db.path").(string)); err != nil {
		log.Fatalf("open database error: %v", err)
	}
	defer db.Close()

	RawMessagesRepository = repos.NewRawMessgeRepository(db)
	ExpensesRepository = repos.NewExpensesRepository(db)

	botService := bothandlers.NewBot(RawMessagesRepository, ExpensesRepository, Configuration)

	Bot = bot.New(Configuration.Get("telegram.token").(string))

	Bot.SetCommandHandler("start", botService.HandleStart)
	Bot.SetCommandHandler("help", botService.HandleStart)
	Bot.SetCommandHandler("list", botService.HandleUserData)

	Bot.SetCommandHandler("listall", botService.HandleList)

	Bot.SetCommandHandler("sum_this_month", botService.HandleMonthSummary)
	Bot.SetCommandHandler("sum_prev_month", botService.HandlePrevMonthSummary)

	Bot.SetDefaultHandler(botService.HandleDefault)

	log.Fatal(Bot.Listen(context.Background()))
}
