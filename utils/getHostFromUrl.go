package utils

import (
	"strings"
)

func GetHostFromURL(url string) string {
	urlParts := strings.Split(url, "/")
	if len(urlParts) > 2 {
		host := urlParts[2]
		if strings.HasPrefix(host, "http://") {
			host = strings.TrimPrefix(host, "http://")
			host = "http://" + host
		} else if strings.HasPrefix(host, "https://") {
			host = strings.TrimPrefix(host, "https://")
			host = "https://" + host
		}
		return host
	}
	return ""
}
