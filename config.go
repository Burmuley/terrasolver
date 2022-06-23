package main

import (
	"os"
)

func readConfigEnv(config map[string]string) map[string]string {
	configVars := make([]string, 0, len(config))
	for k := range config {
		configVars = append(configVars, k)
	}

	for _, v := range configVars {
		if e := os.Getenv(v); e != "" {
			config[v] = e
		}
	}

	return config
}

func strIsInSlice(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}

	return false
}
