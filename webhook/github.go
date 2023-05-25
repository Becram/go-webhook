package webhook

// parse errors
import (
	"encoding/json"
	"fmt"
	"go-webhook/config"
	"go-webhook/psql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/webhooks/v6/github"
)

type Mail interface {
	getTemplate(*template.Template) string
}

var temp *template.Template

func init() {
	// template.Must takes the reponse of template.ParseFiles and does error checking
	temp = template.Must(template.ParseFiles(os.Getenv("SG_EMAIL_TMPL_FILE")))
}

// The function handles incoming webhook data from GitHub and triggers different actions based on the
// type of event received.
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
		// go HandleReleaseEvent(payload.(github.ReleasePayload))

	case github.PullRequestPayload:
		go HandlePullRequestEvent(payload.(github.PullRequestPayload))

	default:
		fmt.Println("Not a  release or pr event")

	}

}

// The function handles a pull request event, validates its labels, extracts version information, sends
// an email notification, and creates a notification in a PostgreSQL database.
func HandlePullRequestEvent(payload github.PullRequestPayload) {
	labels := getLabels(payload)
	fmt.Printf("PR detected with following labels %s\n", labels)
	if payload.Action == "closed" && payload.PullRequest.Merged && validateLabels(labels) {
		// if payload.Action == "open" || payload.Action == "edited" && validateLabels(labels) {
		version, err := getArthurVersion(payload.PullRequest.Title)
		if err != nil {
			fmt.Printf("Version not found in title error: %s", err)
		}

		app := getAppName(labels)[0]
		title := payload.PullRequest.Title
		body := payload.PullRequest.Body
		arthur := payload.PullRequest.MergedBy.Login
		// arthur := "arthur"
		history := "https://github.com/arthur-crm/" + app + "/releases/tag/" + version

		if err == nil {
			var mail Mail = &PullRequest{Title: "Arthur Version " + version, Body: body, Arthur: arthur, History: history}
			data, err := json.Marshal(mail)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", data)
			emailTPL := mail.getTemplate(temp)
			sg_cfg, err := config.Load()
			if err != nil {
				log.Fatal(err)
			}
			c := New(sg_cfg)
			status := c.SendEmail(emailTPL, title)
			if status != nil {
				fmt.Printf("Email %s", status)
			}
			notification := &psql.Notification{Arthur: arthur, Title: title, Timestamp: time.Now(), Version: version, Body: body}
			psql.CreateNotification(*notification)
		}
	}
}
