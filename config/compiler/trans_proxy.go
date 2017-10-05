package compiler

import (
	"os"
	"strings"
)

func transformProxy() Transform {
	var (
		noProxy    = getenv("no_proxy")
		httpProxy  = getenv("http_proxy")
		httpsProxy = getenv("https_proxy")
	)
	return transformEnv(map[string]string{
		"no_proxy":    noProxy,
		"NO_PROXY":    noProxy,
		"http_proxy":  httpProxy,
		"HTTP_PROXY":  httpProxy,
		"HTTPS_PROXY": httpsProxy,
		"https_proxy": httpsProxy,
	})
}

func getenv(name string) (value string) {
	name = strings.ToUpper(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	name = strings.ToLower(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	return
}
