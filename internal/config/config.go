package config

import (
	"fmt"
	"time"

	"github.com/alexflint/go-arg"
)

type Config struct {
	TelegramToken  string `arg:"--telegram-token,env:DA_TELEGRAM_TOKEN"`
	TelegramChatID string `arg:"--telegram-chat-id,env:DA_TELEGRAM_CHAT_ID"`

	SlackWebhookURL string `arg:"--slack-webhook-url,env:DA_SLACK_WEBHOOK_URL"`
	SlackChannel    string `arg:"--slack-channel,env:DA_SLACK_CHANNEL"`

	EmailSMTPHost     string   `arg:"--email-smtp-host,env:DA_EMAIL_SMTP_HOST"`
	EmailSMTPPort     int      `arg:"--email-smtp-port,env:DA_EMAIL_SMTP_PORT" default:"587"`
	EmailFrom         string   `arg:"--email-from,env:DA_EMAIL_FROM_ADDRESS"`
	EmailTo           []string `arg:"--email-to,env:DA_EMAIL_TO_ADDRESSES"`
	EmailSMTPUsername string   `arg:"--email-username,env:DA_EMAIL_SMTP_USERNAME"`
	EmailSMTPPassword string   `arg:"--email-password,env:DA_EMAIL_SMTP_PASSWORD"`

	NoDebounce      bool `arg:"--no-debounce,env:DA_NO_DEBOUNCE"`
	DebounceSeconds int  `arg:"--debounce-seconds,env:DA_DEBOUNCE_SECONDS" default:"5"`
	Debug           bool `arg:"--debug,env:DA_DEBUG"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	arg.MustParse(&cfg)
	return &cfg, nil
}

func (c *Config) DebounceDuration() time.Duration {
	if c.DebounceSeconds < 1 {
		c.DebounceSeconds = 1
	}

	return time.Duration(c.DebounceSeconds) * time.Second
}

func (c *Config) PrintValues() {
	fmt.Println("Config values")
	fmt.Println("-------------")
	fmt.Printf("TelegramToken:     %s\n", c.TelegramToken)
	fmt.Printf("TelegramChatID:    %s\n", c.TelegramChatID)
	fmt.Printf("SlackWebhookURL:   %s\n", c.SlackWebhookURL)
	fmt.Printf("SlackChannel:      %s\n", c.SlackChannel)
	fmt.Printf("EmailSMTPHost:     %s\n", c.EmailSMTPHost)
	fmt.Printf("EmailSMTPPort:     %d\n", c.EmailSMTPPort)
	fmt.Printf("EmailFrom:         %s\n", c.EmailFrom)
	fmt.Printf("EmailTo:           %v\n", c.EmailTo)
	fmt.Printf("EmailSMTPUsername: %s\n", c.EmailSMTPUsername)
	fmt.Printf("EmailSMTPPassword: %s\n", c.EmailSMTPPassword)
	fmt.Printf("NoDebounce:        %t\n", c.NoDebounce)
	fmt.Printf("DebounceSeconds:   %d\n", c.DebounceSeconds)
	fmt.Printf("Debug:             %t\n", c.Debug)
}
