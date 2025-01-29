package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Docker struct {
		Host string `mapstructure:"host" env:"DOCKER_HOST"`
		TLS  bool   `mapstructure:"tls" env:"DOCKER_TLS"`
	} `mapstructure:"docker"`

	Notifications struct {
		Telegram struct {
			Token  string `mapstructure:"token" env:"TELEGRAM_TOKEN"`
			ChatID string `mapstructure:"chat_id" env:"TELEGRAM_CHAT_ID"`
		} `mapstructure:"telegram"`

		Slack struct {
			WebhookURL string `mapstructure:"webhook_url" env:"SLACK_WEBHOOK_URL"`
			Channel    string `mapstructure:"channel" env:"SLACK_CHANNEL"`
		} `mapstructure:"slack"`

		Email struct {
			SMTPHost     string   `mapstructure:"smtp_host" env:"EMAIL_SMTP_HOST"`
			SMTPPort     int      `mapstructure:"smtp_port" env:"EMAIL_SMTP_PORT"`
			FromAddress  string   `mapstructure:"from_address" env:"EMAIL_FROM_ADDRESS"`
			ToAddresses  []string `mapstructure:"to_addresses" env:"EMAIL_TO_ADDRESSES"`
			SMTPUsername string   `mapstructure:"smtp_username" env:"EMAIL_SMTP_USERNAME"`
			SMTPPassword string   `mapstructure:"smtp_password" env:"EMAIL_SMTP_PASSWORD"`
		} `mapstructure:"email"`
	} `mapstructure:"notifications"`

	ConfigFile string `mapstructure:"-"` // Don't unmarshal this field
}

// NewConfig creates a new Config instance with default values
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig loads configuration from all sources
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	if err := bindFlags(cmd, v); err != nil {
		return nil, fmt.Errorf("error binding flags: %w", err)
	}

	if err := loadConfigFile(v); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	bindEnvVariables(v)

	config := NewConfig()
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("docker.host", "unix:///var/run/docker.sock")
	v.SetDefault("docker.tls", false)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	flags := cmd.Flags()

	// Bind each flag to its respective config key
	if err := v.BindPFlag("docker.host", flags.Lookup("docker-host")); err != nil {
		return err
	}
	if err := v.BindPFlag("docker.tls", flags.Lookup("docker-tls")); err != nil {
		return err
	}
	// Add more flag bindings as needed

	return nil
}

func loadConfigFile(v *viper.Viper) error {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.docker-alerts")
	v.AddConfigPath("/etc/docker-alerts")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
	}

	return nil
}

func bindEnvVariables(v *viper.Viper) {
	v.SetEnvPrefix("DA") // DA for Docker Alerts
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
