package config

import (
	"errors"
	"os"
	"strings"
)

const (
	EnvKMSKeyID = "CTR_KMS_KEY_ID"
	EnvSecret   = "CTR_KMS_SECRET"
)

// Common holds configuration required by all CLI binaries.
type Common struct {
	KeyID  string
	Secret string
}

// LoadCommon loads the shared configuration, supplying sensible defaults if unset.
func LoadCommon() (Common, error) {
	keyID := strings.TrimSpace(os.Getenv(EnvKMSKeyID))
	if keyID == "" {
		keyID = "mock-key"
	}

	secret := os.Getenv(EnvSecret)
	if secret == "" {
		secret = "local-dev-secret"
	}

	if len(secret) < 8 {
		return Common{}, errors.New("config: CTR_KMS_SECRET must be at least 8 characters")
	}

	return Common{
		KeyID:  keyID,
		Secret: secret,
	}, nil
}

// ParseList converts a comma-separated string into a slice, trimming whitespace and skipping blanks.
func ParseList(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

// ParseAssignments converts comma-separated key=value pairs into a map.
func ParseAssignments(raw string) (map[string]string, error) {
	if raw == "" {
		return nil, nil
	}
	assignments := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, errors.New("config: constraints must use key=value format")
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key == "" {
			return nil, errors.New("config: constraint key cannot be empty")
		}
		assignments[key] = val
	}
	return assignments, nil
}
