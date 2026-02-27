package encoder

import "encoding/base64"

func Base64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Apply(mode, input string) string {
	switch mode {
	case "base64":
		return Base64(input)
	default:
		return input
	}
}
