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

func (t *TelegramNotifier) sendMessage(ctx context.Context, chatId string, message string, debug bool) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	params := url.Values{}
	params.Add("chat_id", chatId)
	params.Add("text", message)
	params.Add("parse_mode", "HTML")

	if debug {
		fmt.Printf("Sending TG message %v\n", params)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		if debug {
			fmt.Printf("Failed API call - code: %d\n%v\n", resp.StatusCode, string(body))
		}
		return fmt.Errorf("telegram API returned non-200 status code: %d\n%v\n", resp.StatusCode, string(body))
	}

	return nil
}

func (t *TelegramNotifier) Notify(ctx context.Context, event Event, debug bool) error {
	message := event.HTML()

	return t.sendMessage(ctx, t.chatID, message, debug)
}

func (t *TelegramNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	// TODO: group by chatId, allow docker lables to override chat id

	messages := []string{}

	for _, n := range events {
		messages = append(messages, n.HTML())
	}

	return t.sendMessage(ctx, t.chatID, strings.Join(messages, "\n"), debug)
}
