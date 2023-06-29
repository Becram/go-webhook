package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Becram/go-webhook/internal/models"
	"github.com/go-playground/webhooks/github"
)

var msg models.MailData
var temp *template.Template

func HandlePullRequestEvent(payload github.PullRequestPayload) (models.MailData, error) {
	labels := getLabels(payload)
	fmt.Printf("PR detected with following labels %s\n", labels)
	if payload.Action == "closed" && payload.PullRequest.Merged && validateLabels(labels) {
		// if payload.Action == "open" || payload.Action == "edited" && validateLabels(labels) {
		prData, err := GetData(payload)
		if err != nil {
			return models.MailData{}, errors.New("error parsing the webhook data, may be some values are missing")
		}
		temp = template.Must(template.ParseFiles(os.Getenv("TEMPLATE_EMAIL_PATH")))
		msg.To = os.Getenv("EMAIL_TO")
		msg.From = os.Getenv("EMAIL_FROM")
		msg.Subject = os.Getenv("EMAIL_SUBJECT")
		msg.Template = temp
		msg.Content = prData

		return msg, nil
	}
	return models.MailData{}, errors.New("pr request didnt match, pr type or labels")
}

func GetData(payload github.PullRequestPayload) (models.Content, error) {
	var data models.Content
	labels := getLabels(payload)
	version, ok := getArthurVersion(payload.PullRequest.Title)
	if !ok {
		log.Println("couldnt get the version in the pr")
		return models.Content{}, errors.New("version not found")
	}
	appName := getAppName(labels)[0]

	data.Version = version
	data.App = appName
	data.Title = payload.PullRequest.Title
	data.Body = payload.PullRequest.Body
	data.Arthur = payload.PullRequest.MergedBy.Login
	data.History = "https://github.com/arthur-crm/" + appName + "/releases/tag/" + version

	return data, nil

}

func getArthurVersion(str string) (string, bool) {
	// if !strings.HasPrefix(str, os.Getenv("PR_PREFIX")) {
	// 	return "", errors.New("not a arthur release")
	// }
	var re = regexp.MustCompile(`(?m)v([0-9]+)\.([0-9]+)\.([0-9]+)`)
	match := re.FindAllString(str, -1)
	if len(match) < 1 {
		return "", false
	}
	log.Println("Release version detected", match[0])
	return match[0], true
}

// This function retrieves the labels of a pull request from a GitHub payload and returns them as a
// string slice.
func getLabels(payload github.PullRequestPayload) []string {
	labels := payload.PullRequest.Labels
	var out []string
	for _, value := range labels {
		out = append(out, value.Name)
	}
	return out
}

// The function removes the string "prod" from a given slice of strings.
func getAppName(s []string) []string {
	remove := "prod"
	for i, v := range s {
		if v == remove {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// This function validates if any string in a given array matches with a list of strings obtained from
// an environment variable.
func validateLabels(check []string) bool {
	if !stringInSlice("prod", check) {
		fmt.Println("prod label not in the pr labels...skipping alert")
		return false
	}
	svc := strings.Split(os.Getenv("ALERT_SERVICE_LIST"), ",")
	for _, i := range check {
		for _, j := range svc {
			if j == i {
				fmt.Printf("Service %s found in alert list...Alerting\n", j)
				return true
			}
		}
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
