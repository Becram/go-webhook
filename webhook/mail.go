package webhook

import (
	"bytes"
	"context"
	"fmt"
	"go-webhook/config"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/sendgrid"
)

type (
	Client struct {
		cfg      *config.Config
		notify   *notify.Notify
		sendFunc sendFunc
	}
	sendFunc func(ctx context.Context, subject, message string) error
)

func New(cfg *config.Config) *Client {
	n := notify.New()
	c := &Client{
		cfg:      cfg,
		notify:   n,
		sendFunc: n.Send,
	}
	if cfg.Notify.Email != nil {
		c.notify.UseServices(c.email())
	}

	return c
}

// type Release struct {
// 	Name    string `json:"name"`
// 	Body    string `json:"body"`
// 	Arthur  string `json:"arthur"`
// 	History string `json:"history"`
// 	Labels  string `json:"labels"`
// }

type PullRequest struct {
	Title   string `json:"name"`
	Body    string `json:"body"`
	Arthur  string `json:"arthur"`
	History string `json:"history"`
	Labels  string `json:"labels"`
}

func (c *Client) SendEmail(template string, subject string) error {
	// Create a telegram service. Ignoring error for demo simplicity.
	var emails = strings.Split(os.Getenv("SG_TO_LIST"), ",")
	fmt.Printf("Sending email to %s\n", emails[:])
	err := c.sendFunc(
		context.Background(),
		subject,
		template,
	)
	if err != nil {
		return err
	}
	return nil
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

func (c *Client) email() notify.Notifier {
	cfg := c.cfg.Notify.Email
	m := sendgrid.New(cfg.Api, getEnv("SG_FROM", "notify@arthuronline.co.uk"), getEnv("SG_FROM_NAME", "Go webhook app"))

	m.AddReceivers(cfg.Receivers...)
	return m
}
