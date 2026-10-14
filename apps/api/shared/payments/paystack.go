package payments

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Provider interface {
	Initialize(context.Context, InitRequest) (InitResponse, error)
	VerifyWebhook([]byte, string) bool
}
type InitRequest struct {
	Email       string
	AmountMinor int64
	Currency    string
	Reference   string
	CallbackURL string
}
type InitResponse struct {
	Reference        string `json:"reference"`
	AuthorizationURL string `json:"authorization_url"`
}
type Paystack struct {
	SecretKey string
	BaseURL   string
	Client    *http.Client
}

func (p *Paystack) Initialize(ctx context.Context, input InitRequest) (InitResponse, error) {
	body, _ := json.Marshal(map[string]any{"email": input.Email, "amount": input.AmountMinor, "currency": input.Currency, "reference": input.Reference, "callback_url": input.CallbackURL})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(p.BaseURL, "/")+"/transaction/initialize", bytes.NewReader(body))
	if err != nil {
		return InitResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+p.SecretKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := p.Client.Do(req)
	if err != nil {
		return InitResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return InitResponse{}, errors.New("payment provider initialization failed")
	}
	var output struct {
		Status bool         `json:"status"`
		Data   InitResponse `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil || !output.Status {
		return InitResponse{}, errors.New("invalid payment provider response")
	}
	return output.Data, nil
}

func (p *Paystack) VerifyWebhook(body []byte, signature string) bool {
	mac := hmac.New(sha512.New, []byte(p.SecretKey))
	_, _ = mac.Write(body)
	return hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(signature))
}
func NewPaystack(key, base string) *Paystack {
	if base == "" {
		base = "https://api.paystack.co"
	}
	return &Paystack{SecretKey: key, BaseURL: base, Client: http.DefaultClient}
}
func PaymentReference() string { return "nox_" + uuid.NewString() }
