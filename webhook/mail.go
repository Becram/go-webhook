package webhook

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"go-webhook/utility"

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

func (rel Release) Send(template string) error {
	// Create a telegram service. Ignoring error for demo simplicity.
	fmt.Printf("Sending email to %s\n", os.Getenv("SG_TO_LIST"))
	sendgridService := sendgrid.New(utility.GetEnv("SG_API_KEY", ""), utility.GetEnv("SG_FROM", "notify@arthuronline.co.uk"), utility.GetEnv("SG_FROM_NAME", "Go webhook app"))
	sendgridService.AddReceivers(utility.GetEnv("SG_TO_LIST", "bikram.dhoju@gmail.com"))
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
	fmt.Printf("Sending email to %s\n", os.Getenv("SG_TO_LIST"))
	sendgridService := sendgrid.New(utility.GetEnv("SG_API_KEY", ""), utility.GetEnv("SG_FROM", "notify@arthuronline.co.uk"), utility.GetEnv("SG_FROM_NAME", "Go webhook app"))
	sendgridService.AddReceivers(utility.GetEnv("SG_TO_LIST", "bikram.dhoju@gmail.com"))
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
