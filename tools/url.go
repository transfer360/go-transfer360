package tools

import (
	"fmt"
	"strings"
)

func AddHttps(url string) string {

	if len(url) == 0 {
		return ""
	}

	if !strings.Contains(url, "https://") && !strings.Contains(url, "http://") {
		return fmt.Sprintf("https://%s", url)
	}

	if strings.Contains(url, "http://") {
		return strings.ReplaceAll(url, "http://", "https://")
	}

	return url
}
