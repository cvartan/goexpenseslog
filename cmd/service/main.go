package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goconfig/reader/yamlreader"
	"github.com/cvartan/goexpenseslog/internal/controller/restcontroller"
	"github.com/cvartan/goexpenseslog/internal/repos"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	var configPath string = *flag.String("config", "../config/config.yml", "set path to configuration file")
	var Configuration *goconfig.ConfigurationManager

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

	restService := restcontroller.NewRestService(RawMessagesRepository, ExpensesRepository)

	http.HandleFunc("GET /rawmessages", restService.GetAllRawMessages)
	http.HandleFunc("GET /expenses", restService.GetAllExpenses)

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(
				"%s:%d",
				Configuration.Get("api.server").(string),
				Configuration.Get("api.port").(int),
			),
			nil,
		),
	)
}
