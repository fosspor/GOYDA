package yandexllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HTTP       *http.Client
	FolderID   string
	APIKey     string
	Completion string
}

func New(folderID, apiKey string) *Client {
	return &Client{
		HTTP:       &http.Client{Timeout: 60 * time.Second},
		FolderID:   strings.TrimSpace(folderID),
		APIKey:     strings.TrimSpace(apiKey),
		Completion: "https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
	}
}

func (c *Client) Enabled() bool {
	return c.FolderID != "" && c.APIKey != ""
}

// CompletionRaw вызывает Yandex Foundation Models completion (если заданы ключи).
// Формат тела — упрощённый; при ошибке API вернётся текст ответа для отладки.
func (c *Client) CompletionRaw(ctx context.Context, userPrompt string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("yandex llm not configured")
	}
	body := map[string]any{
		"modelUri": fmt.Sprintf("gpt://%s/yandexgpt/latest", c.FolderID),
		"completionOptions": map[string]any{
			"stream":      false,
			"temperature": 0.3,
			"maxTokens":   2000,
		},
		"messages": []map[string]any{
			{"role": "system", "text": "Ты помощник для туристических маршрутов по Краснодарскому краю. Отвечай кратко структурированным JSON с полями title, summary, stops (массив {name, reason, seasonHint})."},
			{"role": "user", "text": userPrompt},
		},
	}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Completion, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+c.APIKey)
	req.Header.Set("x-folder-id", c.FolderID)

	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("yandex api %s: %s", res.Status, string(b))
	}
	return string(b), nil
}
