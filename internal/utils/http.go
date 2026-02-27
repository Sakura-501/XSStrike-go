package utils

import (
	"encoding/json"
	"net/url"
	"strings"
)

// ExtractHeaders parses a raw header string (line separated) to a key-value map.
func ExtractHeaders(raw string) map[string]string {
	headers := map[string]string{}
	if raw == "" {
		return headers
	}

	normalized := strings.ReplaceAll(raw, "\\n", "\n")
	for _, line := range strings.Split(normalized, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.TrimSuffix(value, ",")
		if key != "" {
			headers[key] = value
		}
	}
	return headers
}

// GetURL returns URL without query string when isGET is true.
func GetURL(rawURL string, isGET bool) string {
	if !isGET {
		return rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	u.RawQuery = ""
	return u.String()
}

// ParseParams extracts params from URL query, form body, or JSON body.
func ParseParams(rawURL, data string, jsonData bool) map[string]string {
	params := map[string]string{}

	u, err := url.Parse(rawURL)
	if err == nil && u.RawQuery != "" {
		for key, values := range u.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			} else {
				params[key] = ""
			}
		}
		return params
	}

	if strings.TrimSpace(data) == "" {
		return nil
	}

	if jsonData {
		decoded := map[string]interface{}{}
		if err := json.Unmarshal([]byte(data), &decoded); err == nil {
			for key, value := range decoded {
				switch v := value.(type) {
				case string:
					params[key] = v
				default:
					encoded, _ := json.Marshal(v)
					params[key] = string(encoded)
				}
			}
			return params
		}
	}

	for _, part := range strings.Split(data, "&") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		key := kv[0]
		value := ""
		if len(kv) == 2 {
			value = kv[1]
		}
		params[key] = value
	}

	if len(params) == 0 {
		return nil
	}
	return params
}
