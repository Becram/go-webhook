package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Becram/go-webhook/internal/models"

	// mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/sendgrid"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			err := Notify(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func getTemplate(t *template.Template, m models.TemplateContent) (string, error) {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, m)
	if err != nil {
		return " ", err
	}
	return buffer.String(), nil
}

func Notify(m models.MsgData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("notify: %v", e)
		}
	}()
	var msgToSend string
	var subject string

	api := envOrDefault("SG_API_KEY", " ")
	slackToken := envOrDefault("SLACK_TOKEN", "")
	to := strings.Split(viper.GetString("notification.email.to"), ",")
	from := viper.GetString("notification.email.from")
	sender := viper.GetString("notification.email.sender")
	subject = viper.GetString("notification.email.subject")
	slackChannel := viper.GetString("notification.slack.channel")

	notifier := notify.New()
	service := m.Type
	switch service {
	case "telegram":
		tg, err := telegram.New(env("TELEGRAM_TOKEN"))
		if err != nil {
			return fmt.Errorf("telegram: %w", err)
		}
		tg.SetParseMode("MarkdownV2")
		chatIDStr := env("TELEGRAM_CHAT_ID")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("telegram parse chat_id: %w", err)
		}
		tg.AddReceivers(chatID)
		notifier.UseServices(tg)
	case "email":
		msgToSend, err = getTemplate(m.Template, m.TemplateContent)
		if err != nil {
			log.Println("cannot build message")
		}
		log.Println("Email being sent to:", to)
		mail := sendgrid.New(api, from, sender)
		mail.AddReceivers(to...)
		notifier.UseServices(mail)
	case "slack":
		msgToSend = m.Body
		log.Println("Slack msg being sent to:", slackChannel)
		subject = ":fire_extinguisher:\n"
		slack := slack.New(slackToken)
		slack.AddReceivers(slackChannel)
		notifier.UseServices(slack)
	}
	err = notifier.Send(context.Background(), subject, msgToSend)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func envOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func env(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("environment variable %s is not set", key))
}
