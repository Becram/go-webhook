package webhook

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/sendgrid"
)

type Release struct {
	Name    string `json:"name"`
	Body    string `json:"body"`
	Arthur  string `json:"arthur"`
	History string `json:"history"`
	Labels  string `json:"labels"`
}

type PullRequest struct {
	Title   string `json:"name"`
	Body    string `json:"body"`
	Arthur  string `json:"arthur"`
	History string `json:"history"`
	Labels  string `json:"labels"`
}

// This is a method attached to the `Release` struct. It takes a `template` string as input and sends
// an email notification using the `sendgrid` service with the `template` string as the email body. It
// returns an error if there is any issue with sending the email.
func (rel Release) Send(template string) error {
	// Create a telegram service. Ignoring error for demo simplicity.
	fmt.Printf("Sending email to %s\n", os.Getenv("SG_TO_LIST"))
	sendgridService := sendgrid.New(getEnv("SG_API_KEY", ""), getEnv("SG_FROM", "notify@arthuronline.co.uk"), getEnv("SG_FROM_NAME", "Go webhook app"))
	sendgridService.AddReceivers(getEnv("SG_TO_LIST", "bikram.dhoju@gmail.com"))
	notify.UseServices(sendgridService)
	err := notify.Send(context.Background(),
		rel.Name,
		template,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pr PullRequest) Send(template string) error {
	// Create a telegram service. Ignoring error for demo simplicity.
	emails := strings.Split(os.Getenv("SG_TO_LIST"), ",")

	fmt.Printf("Sending email to %s\n", emails[:])
	sendgridService := sendgrid.New(getEnv("SG_API_KEY", ""), getEnv("SG_FROM", "notify@arthuronline.co.uk"), getEnv("SG_FROM_NAME", "Go webhook app"))
	sendgridService.AddReceivers(emails...)
	notify.UseServices(sendgridService)
	err := notify.Send(context.Background(),
		pr.Title,
		template,
	)
	if err != nil {
		return err
	}
	return nil
}

// This is a method attached to the `Release` struct that takes a pointer to a `template.Template` as
// input and returns a string. It updates the template with the values of the `Release` struct and
// returns the resulting string.
func (rel Release) getTemplate(t *template.Template) string {
	fmt.Printf("Updating template for release %s\n", rel.History)
	var buffer bytes.Buffer
	err := t.Execute(&buffer, rel)
	if err != nil {
		log.Fatalln(err)
	}

	return buffer.String()
}

func (pr PullRequest) getTemplate(t *template.Template) string {
	fmt.Printf("Updating template for pr %s\n", pr.Title)
	var buffer bytes.Buffer
	err := t.Execute(&buffer, pr)
	if err != nil {
		log.Fatalln(err)
	}

	return buffer.String()
}
