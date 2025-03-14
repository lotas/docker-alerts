package notifications

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailNotifier struct {
	host        string
	port        int
	fromAddress string
	toAddresses []string
	auth        smtp.Auth
}

func NewEmailNotifier(host string, port int, fromAddress string, toAddresses []string) *EmailNotifier {
	return &EmailNotifier{
		host:        host,
		port:        port,
		fromAddress: fromAddress,
		toAddresses: toAddresses,
	}
}

func (e *EmailNotifier) SetAuth(username, password string) {
	e.auth = smtp.PlainAuth("", username, password, e.host)
}

func (e *EmailNotifier) Notify(ctx context.Context, event Event, debug bool) error {
	subject := event.Type + " " + event.Action
	body := event.Text()

	message := []string{
		"From: " + e.fromAddress,
		"To: " + strings.Join(e.toAddresses, ","),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}

	emailBody := strings.Join(message, "\r\n")

	addr := fmt.Sprintf("%s:%d", e.host, e.port)

	var err error
	if e.auth != nil {
		err = smtp.SendMail(
			addr,
			e.auth,
			e.fromAddress,
			e.toAddresses,
			[]byte(emailBody),
		)
	} else {
		// Send without authentication
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer client.Close()

		if err := client.Mail(e.fromAddress); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		for _, addr := range e.toAddresses {
			if err := client.Rcpt(addr); err != nil {
				return fmt.Errorf("failed to add recipient %s: %w", addr, err)
			}
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to start email data: %w", err)
		}

		_, err = w.Write([]byte(emailBody))
		if err != nil {
			return fmt.Errorf("failed to write email body: %w", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close email writer: %w", err)
		}

		err = client.Quit()
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (c *EmailNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	for _, n := range events {
		c.Notify(ctx, n, debug)
	}
	return nil
}
