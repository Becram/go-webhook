package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Becram/go-webhook/internal/config"
	"github.com/Becram/go-webhook/internal/handlers"
	"github.com/Becram/go-webhook/internal/models"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var (
	app     config.AppConfig
	session *scs.SessionManager
	cfgFile string
)

func initConfig() {
	log.Info("initializing viper config")
	flag.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")
	flag.Parse()
	if cfgFile != "" {
		// Use config file from the flag.
		log.Printf("using configfile from %s", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("cannot find homedir", err)
		}
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		log.Printf("using %s/config.yaml as configure file ", home)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	initConfig()
	var show map[string]string

	e := viper.UnmarshalKey("worker", &show)
	if e != nil {
		log.Println("cannot read config", e)
	}
	for k, v := range show {
		log.Printf("Queue: %s, Deployment: %s", k, v)
	}

	err := run()
	if err != nil {
		log.Fatal(err)
	}
	// e := runCron()
	// if e != nil {
	// 	log.Fatal(err)
	// }
	// e := runCron()
	// if e != nil {
	// 	log.Fatal(err)
	// }
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

// func runCron() error {
// 	task := func(in string) {
// 		log.Println("run something here")
// 		// webhook.Values()
// func runCron() error {
// 	task := func(in string) {
// 		log.Println("run something here")
// 		// webhook.Values()

// 	}
// 	s := gocron.NewScheduler(time.UTC)
// 	s.SingletonModeAll()
// 	j, err := s.Cron("*/1 * * * *").Do(task, "test")
// 	if err != nil {
// 		log.Fatalln("error scheduling job", err)
// 	}
// 	s.StartAsync()
// 	log.Printf("Next Run: %s", j.NextRun())
// 	}
// 	s := gocron.NewScheduler(time.UTC)
// 	s.SingletonModeAll()
// 	j, err := s.Cron("*/1 * * * *").Do(task, "test")
// 	if err != nil {
// 		log.Fatalln("error scheduling job", err)
// 	}
// 	s.StartAsync()
// 	log.Printf("Next Run: %s", j.NextRun())

// 	return nil
// 	return nil

// }
// }
