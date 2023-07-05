package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Becram/go-webhook/internal/config"
	"github.com/Becram/go-webhook/internal/models"
	"github.com/Becram/go-webhook/internal/repository"
	"github.com/Becram/go-webhook/internal/webhook"
	"github.com/go-playground/webhooks/github"
)

var Repo *Repository

// The function sets the global variable Repo to the input parameter r.
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// The function creates a new instance of a repository with the given app configuration.
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// The function sets the global variable Repo to the input parameter r.
func NewHandlers(r *Repository) {
	Repo = r
}

// The Home function renders the home page template in Go.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	log.Println("path / called")
}

func (m *Repository) GHWebhook(w http.ResponseWriter, req *http.Request) {
	// PrintBody(req)

	hook, err := github.New(github.Options.Secret(os.Getenv("GH_SECRET")))
	if err != nil {
		log.Fatal("wrong secret addressed: %w", err)
		log.Println(req.Header)

	}
	payload, err := hook.Parse(req, github.ReleaseEvent, github.PullRequestEvent)

	if err != nil {
		if err == github.ErrEventNotFound {
			log.Println(err)
		}
	}
	// fmt.Println(payload.(github.PullRequestPayload))

	switch payload.(type) {

	case github.ReleasePayload:
		fmt.Println("Release Webhook triggered")
		// go HandleReleaseEvent(payload.(github.ReleasePayload))

	case github.PullRequestPayload:

		data, err := HandlePullRequestEvent(payload.(github.PullRequestPayload))
		if err != nil {
			log.Println("Couldnt process", err)
		}
		m.App.MailChan <- data

	default:
		fmt.Println("Not a  release or pr event")

	}

}

func (m *Repository) AlertWebhook(w http.ResponseWriter, r *http.Request) {

	var msg models.HookMessage
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	e := json.Unmarshal(bodyBytes, &msg)
	if e != nil {
		log.Println(e)
		// nozzle.printError("opening config file", err.Error())
	}

	log.Println(msg.CommonLabels)

	if msg.Status == "firing" && msg.Receiver == "go-webhook" {
		log.Println("Firing Alert detected")
		for _, i := range msg.Alerts {
			log.Println(i.Labels["queue"])
			webhook.RestartPod(i.Labels["queue"])
		}
	}

}
