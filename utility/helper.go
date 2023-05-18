package utility

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

func DoesEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetArthurVersion(str string) (string, error) {
	if !strings.HasPrefix(str, os.Getenv("PR_PREFIX")) {
		return "", errors.New("not a arthur release")
	}
	var re = regexp.MustCompile(`(?m)v[0-9]\.[0-9]\.[0-9]+[0-9]`)
	match := re.FindAllString(str, -1)
	if len(match) < 1 {
		return "", errors.New("no version match")
	}
	return match[0], nil
}
