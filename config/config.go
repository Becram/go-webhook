package config

import "os"

type (
	// Config represents the configuration for the application.
	Config struct {
		Notify Notify `yaml:"notify"`
	}
	// Page represents a singular Pages to monitor.
	// Notify represents the notification settings for when
	// an element has changed in the DOM.
	Notify struct {
		Email *Email `yaml:"email"`
		Slack *Slack `yaml:"slack"`
	}
	// Email represents SMTP email credentials.
	Email struct {
		Api       string   `yaml:"address"`
		Receivers []string `yaml:"receivers"`
	}
	// Slack represents the Slack credentials configuration.
	Slack struct {
		Token     string `yaml:"token"`
		ChannelID string `yaml:"channel_id"`
	}
)

// The function loads configuration data and returns a Config struct with Notify information.
func Load() (*Config, error) {
	n := &Notify{Email: &Email{Api: os.Getenv("SG_API_KEY"), Receivers: []string{os.Getenv("SG_TO_LIST")}}}
	config := &Config{
		Notify: *n,
	}

	return config, nil
}
