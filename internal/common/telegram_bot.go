package common

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SendTelegramMessage sends a message to a Telegram chat via a bot.
func SendTelegramMessage(botToken, chatID, message, proxyURL string) error {
	if botToken == "" || chatID == "" {
		return nil // Silent skip if not configured
	}

	// Use GET instead of POST to bypass firewalls that block POST requests to Telegram
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	params := url.Values{}
	params.Add("chat_id", chatID)
	params.Add("text", message)
	params.Add("parse_mode", "HTML")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	transport := &http.Transport{}
	if proxyURL != "" {
		pURL, err := url.Parse(proxyURL)
		if err == nil {
			transport.Proxy = http.ProxyURL(pURL)
		}
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		return fmt.Errorf("failed to send telegram request via GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram notification failed with status: %d", resp.StatusCode)
	}

	return nil
}
