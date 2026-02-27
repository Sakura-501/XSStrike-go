package utils

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

// URLPathToMap converts URL path segments to a map where key/value are the same segment value.
func URLPathToMap(rawURL string) map[string]string {
	result := map[string]string{}
	u, err := url.Parse(rawURL)
	if err != nil {
		return result
	}
	parts := strings.Split(strings.TrimPrefix(u.EscapedPath(), "/"), "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		decoded, decodeErr := url.PathUnescape(part)
		if decodeErr != nil {
			decoded = part
		}
		result[decoded] = decoded
	}
	return result
}

// MapToURLPath builds a URL by appending map values as path segments to scheme://host from rawURL.
func MapToURLPath(rawURL string, data map[string]string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	base := u.Scheme + "://" + u.Host

	if len(data) == 0 {
		return base
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	segments := make([]string, 0, len(keys))
	for _, k := range keys {
		segments = append(segments, url.PathEscape(data[k]))
	}
	return base + "/" + strings.Join(segments, "/")
}

// JSONToMap decodes JSON object data to map[string]string.
func JSONToMap(data string) (map[string]string, error) {
	decoded := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &decoded); err != nil {
		return nil, err
	}
	result := map[string]string{}
	for key, value := range decoded {
		switch v := value.(type) {
		case string:
			result[key] = v
		default:
			raw, _ := json.Marshal(v)
			result[key] = string(raw)
		}
	}
	return result, nil
}

// MapToJSON encodes map data to JSON string.
func MapToJSON(data map[string]string) (string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// FlattenParams replaces currentParam value with payload and returns query string.
func FlattenParams(currentParam string, params map[string]string, payload string) string {
	if len(params) == 0 {
		return ""
	}
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		value := params[key]
		if key == currentParam {
			value = payload
		}
		pairs = append(pairs, key+"="+value)
	}
	return "?" + strings.Join(pairs, "&")
}
