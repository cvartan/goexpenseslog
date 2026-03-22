package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goexpenseslog/internal/repos"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	var Configuration *goconfig.Configuration

	var RawMessagesRepository *repos.RawMessageRepository
	var ExpensesRepository *repos.ExpensesRepository

	log.Println("Starting bot service")

	Configuration = goconfig.NewConfiguration(
		&goconfig.Options{
			Path:     "../config",
			Filename: "config.yml",
			Format:   "YAML",
		},
	)

	Configuration.Apply()

	var db *sql.DB
	if db, err = sql.Open("sqlite3", Configuration.Get("db.path").String()); err != nil {
		log.Fatalf("open database error: %v", err)
	}
	defer db.Close()

	RawMessagesRepository = repos.NewRawMessgeRepository(db)
	ExpensesRepository = repos.NewExpensesRepository(db)

	restService := NewRestService(RawMessagesRepository, ExpensesRepository)

	http.HandleFunc("GET /rawmessages", restService.GetAllRawMessages)
	http.HandleFunc("GET /expenses", restService.GetAllExpenses)

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(
				"%s:%d",
				Configuration.Get("api.server").String(),
				Configuration.Get("api.port").Int(),
			),
			nil,
		),
	)
}
