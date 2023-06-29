package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Becram/go-webhook/internal/config"
	"github.com/Becram/go-webhook/internal/repository"
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
	log.Println("Nothing to do here")
}

func (m *Repository) GHWebhook(w http.ResponseWriter, req *http.Request) {
	// PrintBody(req)

	hook, err := github.New(github.Options.Secret(os.Getenv("GH_SECRET")))
	if err != nil {
		log.Fatal("wrong secret addressed: %w", err)
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
