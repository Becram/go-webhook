package webhook

import (
	"context"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/amazonsns"
)

// The function sends an SMS notification using Amazon SNS service.
func SendSMS() error {
	snsService, _ := amazonsns.New(getEnv("AWS_ACCESS_KEY_ID", ""), getEnv("AWS_SECRET_ACCESS_KEY", ""), getEnv("AWS_REGION", ""))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	snsService.AddReceivers(getEnv("TOPIC_ARN", ""))
	notify.UseServices(snsService)
	status := notify.Send(context.Background(),
		"Test",
		"Bikram",
	)
	if status != nil {
		return status
	}
	return nil

}
