package notifications

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRoundTripper implements http.RoundTripper for testing HTTP requests.
type mockRoundTripper struct {
	resp       *http.Response
	err        error
	lastReq    *http.Request
	bodyReader *bytes.Buffer
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastReq = req
	if m.bodyReader != nil {
		b, _ := ioutil.ReadAll(req.Body)
		m.bodyReader.Write(b)
	}
	return m.resp, m.err
}

func newMockClient(resp *http.Response, err error, bodyReader *bytes.Buffer) *http.Client {
	return &http.Client{
		Transport: &mockRoundTripper{
			resp:       resp,
			err:        err,
			bodyReader: bodyReader,
		},
	}
}

func TestTelegramNotifier_sendMessage_Success(t *testing.T) {
	var bodyBuf bytes.Buffer
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"ok":true}`)),
	}
	client := newMockClient(mockResp, nil, &bodyBuf)
	notifier := &TelegramNotifier{
		token:  "dummy-token",
		chatID: "12345",
		client: client,
	}

	ctx := context.Background()
	err := notifier.sendMessage(ctx, "12345", "Hello, Telegram!", false)
	require.NoError(t, err)
	assert.Contains(t, bodyBuf.String(), "chat_id=12345")
	assert.Contains(t, bodyBuf.String(), "text=Hello%2C+Telegram%21")
}

func TestTelegramNotifier_sendMessage_HTTPError(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewBufferString("Internal Server Error")),
	}
	client := newMockClient(mockResp, nil, nil)
	notifier := &TelegramNotifier{
		token:  "dummy-token",
		chatID: "12345",
		client: client,
	}

	ctx := context.Background()
	err := notifier.sendMessage(ctx, "12345", "fail test", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "telegram API returned non-200 status code")
}

func TestTelegramNotifier_sendMessage_RequestError(t *testing.T) {
	client := newMockClient(nil, assert.AnError, nil)
	notifier := &TelegramNotifier{
		token:  "dummy-token",
		chatID: "12345",
		client: client,
	}

	ctx := context.Background()
	err := notifier.sendMessage(ctx, "12345", "fail test", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send telegram message")
}

func TestTelegramNotifier_Notify(t *testing.T) {
	var bodyBuf bytes.Buffer
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"ok":true}`)),
	}
	client := newMockClient(mockResp, nil, &bodyBuf)
	notifier := &TelegramNotifier{
		token:  "dummy-token",
		chatID: "12345",
		client: client,
	}

	event := Event{
		Type:   "container",
		Action: "start",
		Name:   "test-container",
		Image:  "test-image:latest",
	}
	ctx := context.Background()
	err := notifier.Notify(ctx, event, false)
	require.NoError(t, err)
	assert.Contains(t, bodyBuf.String(), "test-container")
	assert.Contains(t, bodyBuf.String(), "test-image%3Alatest")
}

func TestTelegramNotifier_NotifyMultiple(t *testing.T) {
	var bodyBuf bytes.Buffer
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"ok":true}`)),
	}
	client := newMockClient(mockResp, nil, &bodyBuf)
	notifier := &TelegramNotifier{
		token:  "dummy-token",
		chatID: "12345",
		client: client,
	}

	events := []Event{
		{Type: "container", Action: "start", Name: "c1", Image: "img1"},
		{Type: "container", Action: "stop", Name: "c2", Image: "img2"},
	}
	ctx := context.Background()
	err := notifier.NotifyMultiple(ctx, events, false)
	require.NoError(t, err)
	assert.Contains(t, bodyBuf.String(), "c1")
	assert.Contains(t, bodyBuf.String(), "c2")
}
