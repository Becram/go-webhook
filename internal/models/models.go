package models

import "html/template"

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  Content
	Template *template.Template
}

type Content struct {
	App     string
	Version string
	Title   string
	Body    string
	Arthur  string
	History string
}
