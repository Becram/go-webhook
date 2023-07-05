package webhook

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

func getDeployment(queue string) (string, error) {
	queueDeployments := map[string]string{
		"remittance-queue": "remittance-worker",
		"batch":            "batch-queue",
		"elastic_search":   "es-queue",
		"webhook":          "webhook-worker",
		"webhook1":         "webhook-worker1",
	}

	for k, v := range queueDeployments {
		if k == queue {
			log.Println("matched", k, v)
			return v, nil
		}
	}
	return "", errors.New("queue doesn't have any deployment map")

}

func RestartPod(queue string) {
	deployment, err := getDeployment(queue)
	if err != nil {
		log.Println(err)
	}
	log.Println("deployment found:", deployment)
	update := UpdateDeployment(deployment)
	log.Println("restarting deployment", update)
}
