package notifications

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

func (t *TelegramNotifier) Notify(ctx context.Context, notification Notification) error {
	message := fmt.Sprintf("*%s*\n%s", notification.Title, notification.Message)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	params := url.Values{}
	params.Add("chat_id", t.chatID)
	params.Add("text", message)
	params.Add("parse_mode", "Markdown")

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
