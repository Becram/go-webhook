package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Becram/go-webhook/internal/models"

	// mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/sendgrid"
	"github.com/nikoksr/notify/service/telegram"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			Notify(msg)
		}
	}()
}

func getTemplate(t *template.Template, m models.Content) (string, error) {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, m)
	if err != nil {
		return " ", err
	}

	return buffer.String(), nil
}

func Notify(m models.MailData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("notify: %v", e)
		}
	}()
	var msgToSend string

	api := envOrDefault("SG_API_KEY", " ")
	to := strings.Split(envOrDefault("EMAIL_TO", "bikram.dhoju@gmail.com"), ",")
	from := envOrDefault("EMAIL_FROM", "arthur")
	sender := envOrDefault("EMAIL_SENDER", "Go Webhook")
	subject := envOrDefault("EMAIL_SUBJECT", "Notification")

	notifier := notify.New()
	service := envOrDefault("NOTIFY_SERVICE", "email")
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
		msgToSend, err = getTemplate(m.Template, m.Content)
		if err != nil {
			log.Println("cannot build message")
		}
		mail := sendgrid.New(api, from, sender)
		mail.AddReceivers(to...)
		notifier.UseServices(mail)
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
