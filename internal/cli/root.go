package cli

import (
	"github.com/lotas/docker-alerts/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cfg := config.NewConfig()

	cmd := &cobra.Command{
		Use:   "docker-alerts",
		Short: "Docker event monitoring and notification system",
		Long: `Docker event monitoring and notification system.
Configuration can be provided via command line flags, environment variables, or config file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, cfg)
		},
	}

	flags := cmd.Flags()
	cmd.Flags().SortFlags = false

	flags.String("docker-host", "",
		"Docker host URL [env: DA_DOCKER_HOST]")
	flags.Bool("docker-tls", false,
		"Use TLS for Docker connection [env: DA_DOCKER_TLS]")

	// Notification flags - Telegram
	flags.String("telegram-token", "",
		"Telegram bot token [env: DA_TELEGRAM_TOKEN]")
	flags.String("telegram-chat-id", "",
		"Telegram chat ID [env: DA_TELEGRAM_CHAT_ID]")

	// Notification flags - Slack
	flags.String("slack-webhook-url", "",
		"Slack webhook URL [env: DA_SLACK_WEBHOOK_URL]")
	flags.String("slack-channel", "",
		"Slack channel [env: DA_SLACK_CHANNEL]")

	// Notification flags - Email
	flags.String("email-smtp-host", "",
		"SMTP server host [env: DA_EMAIL_SMTP_HOST]")
	flags.Int("email-smtp-port", 587,
		"SMTP server port [env: DA_EMAIL_SMTP_PORT]")
	flags.String("email-from", "",
		"Email sender address [env: DA_EMAIL_FROM_ADDRESS]")
	flags.StringSlice("email-to", []string{},
		"Email recipient addresses (comma-separated) [env: DA_EMAIL_TO_ADDRESSES]")
	flags.String("email-username", "",
		"SMTP authentication username [env: DA_EMAIL_SMTP_USERNAME]")
	flags.String("email-password", "",
		"SMTP authentication password [env: DA_EMAIL_SMTP_PASSWORD]")

	// Config file flag
	flags.StringVar(&cfg.ConfigFile, "config", "",
		"Config file path [env: DA_CONFIG]")

	// Add example usage to help text
	cmd.Example = `  # Run with config file
  docker-alerts --config=/path/to/config.yaml

  # Run with environment variables
  export DA_TELEGRAM_TOKEN=your-token
  export DA_TELEGRAM_CHAT_ID=your-chat-id
  docker-alerts

  # Run with command line flags
  docker-alerts --telegram-token=your-token --telegram-chat-id=your-chat-id

  # Send notifications via email
  docker-alerts --email-smtp-host=smtp.gmail.com --email-smtp-port=587 \
    --email-from=sender@example.com --email-to=recipient@example.com \
    --email-username=user --email-password=pass`

	return cmd
}

func run(cmd *cobra.Command, cfg *config.Config) error {
	// Load configuration from all sources
	finalConfig, err := config.LoadConfig(cmd)
	if err != nil {
		return err
	}

	// Start your application with the configuration
	return startApp(finalConfig)
}
