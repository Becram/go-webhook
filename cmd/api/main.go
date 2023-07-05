package main

import (
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/Becram/go-webhook/internal/config"
	"github.com/Becram/go-webhook/internal/handlers"
	"github.com/Becram/go-webhook/internal/models"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	e := runCron()
	if e != nil {
		log.Fatal(err)
	}
	defer close(app.MailChan)
	listenForMail()

	log.Println("Serving application at ", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes((&app)),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

	app.InProduction = false

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	app.UseCache = false

	return nil
}

func runCron() error {
	task := func(in string) {
		log.Println("run something here")
		// webhook.Values()

	}
	s := gocron.NewScheduler(time.UTC)
	s.SingletonModeAll()
	j, err := s.Cron("*/1 * * * *").Do(task, "test")
	if err != nil {
		log.Fatalln("error scheduling job", err)
	}
	s.StartAsync()
	log.Printf("Next Run: %s", j.NextRun())

	return nil

}
