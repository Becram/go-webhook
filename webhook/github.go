package webhook

// parse errors
import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"go-notify/mail"

	"github.com/go-playground/webhooks/v6/github"
)

// type Update interface {
// 	prTemplate() string
// 	prTemplate() string
// }

type Release struct {
	Name    string `json:"name"`
	Body    string `json:"body"`
	Arthur  string `json:"arthur"`
	History string `json:"history"`
	Labels  string `json:"labels"`
}

// type PullRequest struct {
// 	Title   string `json:"name"`
// 	Body    string `json:"body"`
// 	Arthur  string `json:"arthur"`
// 	History string `json:"history"`
// }

var temp *template.Template

func init() {
	// template.Must takes the reponse of template.ParseFiles and does error checking
	temp = template.Must(template.ParseFiles(os.Getenv("SG_EMAIL_TMPL_FILE")))
}

func updateTemplate(rel Release) string {
	fmt.Printf("Updating template %s\n", rel.History)
	var buffer bytes.Buffer
	err := temp.Execute(&buffer, rel)
	if err != nil {
		log.Fatalln(err)
	}

	return buffer.String()

}

func getArthurVersion(str string) (string, error) {
	if !strings.HasPrefix(str, "Bump arthur") {
		return "", errors.New("Not a arthur release")
	}
	var re = regexp.MustCompile(`(?m)v[0-9]\.[0-9]\.[0-9]+[0-9]`)
	match := re.FindAllString(str, -1)
	if len(match) < 1 {
		return "", errors.New("No version match")
	}
	return match[0], nil
}

func GetReleaseData(w http.ResponseWriter, req *http.Request) {
	// PrintBody(req)

	hook, err := github.New(github.Options.Secret(os.Getenv("GH_SECRET")))
	if err != nil {
		log.Fatal("wrong secret addressed: %w", err)
	}
	payload, err := hook.Parse(req, github.ReleaseEvent, github.PullRequestEvent)
	// json.NewEncoder(w).Encode(payload)

	if err != nil {
		if err == github.ErrEventNotFound {
			log.Println(err)
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
	// fmt.Println(payload.(github.Webhook))
	switch payload.(type) {

	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		if *&release.Action == "published" {
			// Do whatever you want from here...
			fmt.Printf("Release Name:  %s\t%s\n", *release.Release.Name, *release.Release.Body)
			out := Release{Name: *release.Release.Name, Body: *release.Release.Body, Arthur: *&release.Release.Author.Login, History: *&release.Release.AssetsURL}
			mail.SendEmail(*release.Release.Name, updateTemplate(out))
		}

	case github.PullRequestPayload:
		pr := payload.(github.PullRequestPayload)
		if *&pr.Action == "opened" || *&pr.Action == "edited" || *&pr.Action == "labeled" && *&pr.PullRequest.Head.User.Login == "arthur-crm" {
			version, err := getArthurVersion(*&pr.PullRequest.Title)
			if err == nil {
				fmt.Printf("PR Title:  %s\t%s\t PR link:%s\n", *&pr.PullRequest.Title, *&pr.PullRequest.Body, *&pr.PullRequest.HTMLURL)
				out := Release{Name: "Arthur Version\t" + version, Body: *&pr.PullRequest.Body, Arthur: *&pr.PullRequest.User.Login, History: *&pr.PullRequest.HTMLURL}
				mail.SendEmail(*&pr.PullRequest.Title, updateTemplate(out))
			}
			fmt.Printf("no new version release detected for arthur\n")

		}
	}
}
