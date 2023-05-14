package mail

import (
	"context"
	"fmt"
	"os"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/sendgrid"
)

var envs = []string{"SG_API_KEY", "SG_FROM", "SG_FROM_NAME", "SG_TO_LIST"}

func appInit() {
	for _, v := range envs {

		if !isEnvExist(v) {
			fmt.Printf("Environment variable %s Doesn't exit", v)
			os.Exit(0)
		}

	}

}
func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func SendEmail(subject string, body string) {
	// Create a telegram service. Ignoring error for demo simplicity.
	fmt.Printf("Sending email to %s\n", os.Getenv("SG_TO_LIST"))
	sendgridService := sendgrid.New(getenv("SG_API_KEY", ""), getenv("SG_FROM", "notify@arthuronline.co.uk"), getenv("SG_FROM_NAME", "Go webhook app"))

	sendgridService.AddReceivers(getenv("SG_TO_LIST", "bikram.dhoju@gmail.com"))

	notify.UseServices(sendgridService)
	_ = notify.Send(context.Background(),
		subject,
		body,
	)

}
