package webhook

// parse errors
import (
	"encoding/json"
	"fmt"
	"go-webhook/utility"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/webhooks/v6/github"
)

type Mail interface {
	Send(string) error
	getTemplate(*template.Template) string
}

var temp *template.Template

func init() {
	// template.Must takes the reponse of template.ParseFiles and does error checking
	temp = template.Must(template.ParseFiles(os.Getenv("SG_EMAIL_TMPL_FILE")))
}

func GetWebhookData(w http.ResponseWriter, req *http.Request) {
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
		HandleReleaseEvent(payload.(github.ReleasePayload))

	case github.PullRequestPayload:
		fmt.Println("PR Webhook triggered")
		HandlePullRequestEvent(payload.(github.PullRequestPayload))

	default:
		fmt.Println("Not a  release or pr event")

	}

}

func HandleReleaseEvent(payload github.ReleasePayload) {
	if payload.Action == "published" {
		var mail Mail = &Release{Name: *payload.Release.Name, Body: *payload.Release.Body, Arthur: *&payload.Release.Author.Login, History: *&payload.Release.AssetsURL}
		// fmt.Printf("Release %v\n", mail)
		data, err := json.Marshal(mail)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
		emailTPL := mail.getTemplate(temp)
		mail.Send(emailTPL)

	}
}

func HandlePullRequestEvent(payload github.PullRequestPayload) {
	if payload.Action == "closed" && payload.PullRequest.Merged && payload.PullRequest.Head.User.Login == "arthur-crm" {
		version, err := utility.GetArthurVersion(payload.PullRequest.Title)
		if err == nil {
			fmt.Printf("Version: %s\n", version)
			var mail Mail = &PullRequest{Title: "Arthur Version\t" + version, Body: payload.PullRequest.Body, Arthur: payload.PullRequest.User.Login, History: payload.PullRequest.HTMLURL}
			data, err := json.Marshal(mail)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", data)
			emailTPL := mail.getTemplate(temp)
			mail.Send(emailTPL)
		}

	}

}
