package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const brevoSendEmailPath = "/v3/smtp/email"

type BrevoConfig struct {
	APIKey      string
	BaseURL     string
	SenderEmail string
	SenderName  string
}

type BrevoProvider struct {
	apiKey      string
	baseURL     string
	senderEmail string
	senderName  string
	client      *http.Client
}

func NewBrevoProvider(cfg BrevoConfig) *BrevoProvider {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.brevo.com"
	}

	return &BrevoProvider{
		apiKey:      cfg.APIKey,
		baseURL:     baseURL,
		senderEmail: cfg.SenderEmail,
		senderName:  cfg.SenderName,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *BrevoProvider) Send(ctx context.Context, message Message) error {
	if p.apiKey == "" {
		return errors.New("brevo api key is empty")
	}
	if p.senderEmail == "" {
		return errors.New("mail sender email is empty")
	}
	if message.ToEmail == "" {
		return errors.New("mail recipient email is empty")
	}

	body := brevoSendEmailRequest{
		Sender: brevoEmailAddress{
			Name:  p.senderName,
			Email: p.senderEmail,
		},
		To: []brevoEmailAddress{
			{
				Name:  message.ToName,
				Email: message.ToEmail,
			},
		},
		Subject:     message.Subject,
		HTMLContent: message.HTMLContent,
		TextContent: message.TextContent,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+brevoSendEmailPath, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", p.apiKey)
	req.Header.Set("content-type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		resBody, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return fmt.Errorf("brevo send email failed: status=%d body=%s", res.StatusCode, string(resBody))
	}

	return nil
}

type brevoSendEmailRequest struct {
	Sender      brevoEmailAddress   `json:"sender"`
	To          []brevoEmailAddress `json:"to"`
	Subject     string              `json:"subject"`
	HTMLContent string              `json:"htmlContent,omitempty"`
	TextContent string              `json:"textContent,omitempty"`
}

type brevoEmailAddress struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}
