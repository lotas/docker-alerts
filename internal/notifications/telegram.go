package notifications

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (t *TelegramNotifier) sendMessage(ctx context.Context, chatId string, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	params := url.Values{}
	params.Add("chat_id", chatId)
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
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API returned non-200 status code: %d\n%v\n", resp.StatusCode, string(body))
	}

	return nil
}

func (t *TelegramNotifier) Notify(ctx context.Context, notification Notification, debug bool) error {
	message := notification.Markdown()

	return t.sendMessage(ctx, t.chatID, message)
}

func (t *TelegramNotifier) NotifyMultiple(ctx context.Context, notifications []Notification, debug bool) error {
	// TODO: group by chatId, notifications might override some settings

	messages := []string{}

	if len(notifications) > 1 {
		messages = append(messages, "Multiple events:")
	}

	for _, n := range notifications {
		messages = append(messages, n.Markdown())
	}

	return t.sendMessage(ctx, t.chatID, strings.Join(messages, "\n---\n"))
}
