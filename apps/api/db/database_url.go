package db

import (
	"net"
	"net/url"
	"strings"
)

func normalizeDatabaseURL(raw string) string {
	if raw == "" {
		return raw
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	if hasSSLMode(parsed) || !shouldDefaultSSLDisable(parsed) {
		return raw
	}

	query := parsed.Query()
	query.Set("sslmode", "disable")
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func hasSSLMode(parsed *url.URL) bool {
	return parsed.Query().Get("sslmode") != ""
}

func shouldDefaultSSLDisable(parsed *url.URL) bool {
	host := parsed.Hostname()
	if host == "" {
		return false
	}

	if host == "localhost" {
		return true
	}

	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}

	return strings.EqualFold(host, "127.0.0.1") || strings.EqualFold(host, "::1")
}
