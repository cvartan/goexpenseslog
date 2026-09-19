package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goexpenseslog/internal/handlers/bothandlers"
	"github.com/cvartan/goexpenseslog/internal/repos"
	"github.com/cvartan/goexpenseslog/internal/repos/sqliteimpl"
	"github.com/cvartan/goexpenseslog/pkg/bot"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var appcontext context.Context = context.Background()

	var err error

	var Configuration *goconfig.Configuration = goconfig.NewConfiguration(
		&goconfig.Options{
			Source: "../config/config.yml",
			Format: "yaml",
		},
	)

	Configuration.Add("telegram.token", goconfig.String)
	Configuration.Add("db.path", goconfig.String)
	Configuration.Add("control.userid", goconfig.Int)
	Configuration.Add("proxy.socks5", goconfig.String)

	appcontext = context.WithValue(appcontext, "config", Configuration)

	var Logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	appcontext = context.WithValue(appcontext, "logger", Logger)

	var Bot *bot.TelegramBot = bot.New(Configuration.Get("telegram.token").String())

	Configuration.Bind(Bot)

	var RawMessagesRepository repos.RawMessageRepository
	var ExpensesRepository repos.ExpensesRepository

	Logger.Info("Starting bot service")
	defer Logger.Info("Bot closed")

	var db *sql.DB
	if db, err = sql.Open("sqlite3", Configuration.Get("db.path").String()); err != nil {
		Logger.Error("open database error", "error", err)
		return
	}
	defer db.Close()

	RawMessagesRepository = sqliteimpl.NewRawMessgeRepository(db)
	ExpensesRepository = sqliteimpl.NewExpensesRepository(db)

	botService := bothandlers.NewBot(RawMessagesRepository, ExpensesRepository)

	Bot.SetCommandHandler("start", botService.HandleStart)
	Bot.SetCommandHandler("help", botService.HandleStart)
	Bot.SetCommandHandler("list", botService.HandleUserData)

	Bot.SetCommandHandler("listall", botService.HandleList)

	Bot.SetCommandHandler("sum_this_month", botService.HandleMonthSummary)
	Bot.SetCommandHandler("sum_prev_month", botService.HandlePrevMonthSummary)

	Bot.SetDefaultHandler(botService.HandleDefault)

	err = Bot.Listen(appcontext)

	if err != nil {
		Logger.Error("Bot finished with error", "error", err)
		os.Exit(1)
	}
}
