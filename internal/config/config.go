package config

import (
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/Becram/go-webhook/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/nikoksr/notify"
)

type (
	AppConfig struct {
		UseCache      bool
		InProduction  bool
		TemplateCache map[string]*template.Template
		ErrorLog      *log.Logger
		InfoLog       *log.Logger
		Session       *scs.SessionManager
		MailChan      chan models.MsgData
		Sendgrid      *notify.Notifier
	}
)
