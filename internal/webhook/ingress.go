package webhook

import (
	"errors"
	"log"

	networkingv1 "k8s.io/api/networking/v1"
)

func GetHosts(ingress *networkingv1.Ingress) ([]string, error) {
	var hosts []string
	name := ingress.Spec.Rules
	for _, host := range name {
		log.Println(host)
		hosts = append(hosts, host.Host)
	}
	if len(hosts) == 0 {
		return nil, errors.New("no hosts found in the url")
	}

	return hosts, nil
}
