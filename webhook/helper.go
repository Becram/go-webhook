package webhook

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/webhooks/v6/github"
)

// The function checks if a given environment variable exists.
func DoesEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

// The function `getEnv` returns the value of an environment variable or a fallback value if the
// variable is not set.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback

	}
	return value
}

// The function extracts the version number from a string that starts with a specific prefix and
// returns an error if the string is not a valid Arthur release or if no version number is found.
func getArthurVersion(str string) (string, error) {
	// if !strings.HasPrefix(str, os.Getenv("PR_PREFIX")) {
	// 	return "", errors.New("not a arthur release")
	// }
	var re = regexp.MustCompile(`(?m)v[0-9]\.[0-9]\.[0-9]+[0-9]`)
	match := re.FindAllString(str, -1)
	if len(match) < 1 {
		return "", errors.New("no version match")
	}
	return match[0], nil
}

// This function retrieves the labels of a pull request from a GitHub payload and returns them as a
// string slice.
func getLabels(payload github.PullRequestPayload) []string {
	labels := payload.PullRequest.Labels
	var out []string
	for _, value := range labels {
		out = append(out, value.Name)
	}
	return out
}

// This function validates if any string in a given array matches with a list of strings obtained from
// an environment variable.
func validateLabels(check []string) bool {
	svc := strings.Split(os.Getenv("ALERT_SERVICE_LIST"), ",")
	for _, i := range check {
		for _, j := range svc {
			if j == i {
				fmt.Printf("Service %s Found in alert list\n", j)
				return true
			}
		}
	}
	return false
}
