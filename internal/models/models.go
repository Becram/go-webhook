package models

import (
	"html/template"
	"time"
)

// MsgData holds an email message
type MsgData struct {
	Body            string
	TemplateContent TemplateContent
	Template        *template.Template
	Type            string
}

type TemplateContent struct {
	App     string
	Version string
	Title   string
	Body    string
	Arthur  string
	History string
}

type HookMessage struct {
	Alerts []struct {
		Annotations  map[string]string `json:"annotations"`
		EndsAt       time.Time         `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Labels       map[string]string `json:"labels"`
		StartsAt     time.Time         `json:"startsAt"`
		Status       string            `json:"status"`
	} `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`

	// Timestamp records when the alert notification was received
	Timestamp string `json:"@timestamp"`
}
