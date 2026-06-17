package client

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const defaultUploadFolder = "nox/posts"

type Config struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadFolder string
}

type Client struct {
	config Config
}

type UploadSignature struct {
	CloudName    string            `json:"cloud_name"`
	APIKey       string            `json:"api_key"`
	ResourceType string            `json:"resource_type"`
	UploadURL    string            `json:"upload_url"`
	PublicID     string            `json:"public_id"`
	Folder       string            `json:"folder"`
	Timestamp    int64             `json:"timestamp"`
	Signature    string            `json:"signature"`
	Params       map[string]string `json:"params"`
}

func New(config Config) *Client {
	config.CloudName = strings.TrimSpace(config.CloudName)
	config.APIKey = strings.TrimSpace(config.APIKey)
	config.APISecret = strings.TrimSpace(config.APISecret)
	config.UploadFolder = normalizeFolder(config.UploadFolder)
	return &Client{config: config}
}

func (c *Client) Configured() bool {
	return c != nil &&
		c.config.CloudName != "" &&
		c.config.APIKey != "" &&
		c.config.APISecret != ""
}

func (c *Client) SignUpload(resourceType string, publicID string) UploadSignature {
	timestamp := time.Now().Unix()
	params := map[string]string{
		"folder":    c.config.UploadFolder,
		"public_id": publicID,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}
	signature := sign(params, c.config.APISecret)
	params["signature"] = signature
	params["api_key"] = c.config.APIKey

	return UploadSignature{
		CloudName:    c.config.CloudName,
		APIKey:       c.config.APIKey,
		ResourceType: resourceType,
		UploadURL:    uploadURL(c.config.CloudName, resourceType),
		PublicID:     publicID,
		Folder:       c.config.UploadFolder,
		Timestamp:    timestamp,
		Signature:    signature,
		Params:       params,
	}
}

func PostPublicID(ownerPersonaID uuid.UUID) string {
	return fmt.Sprintf("posts/%s/%s", ownerPersonaID.String(), uuid.NewString())
}

func normalizeFolder(folder string) string {
	folder = strings.Trim(strings.TrimSpace(folder), "/")
	if folder == "" {
		return defaultUploadFolder
	}
	return folder
}

func sign(params map[string]string, apiSecret string) string {
	keys := make([]string, 0, len(params))
	for key, value := range params {
		if value == "" || key == "file" || key == "api_key" || key == "resource_type" || key == "signature" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+params[key])
	}

	sum := sha1.Sum([]byte(strings.Join(parts, "&") + apiSecret))
	return hex.EncodeToString(sum[:])
}

func uploadURL(cloudName string, resourceType string) string {
	return fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", url.PathEscape(cloudName), resourceType)
}
