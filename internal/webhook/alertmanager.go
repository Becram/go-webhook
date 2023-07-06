package webhook

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getDeployment(queue string) (string, error) {

	var queueMap map[string]string
	err := viper.UnmarshalKey("worker", &queueMap)
	if err != nil {
		log.Fatalln("cannot parse the config file", err)
	}

	for k, v := range queueMap {
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
