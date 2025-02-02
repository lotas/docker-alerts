package config

import (
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

	Debug bool `arg:"--debug,env:DA_DEBUG"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	arg.MustParse(&cfg)
	return &cfg, nil
}
