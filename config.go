package main

import "os"

func readConfigEnv(config map[string]string) map[string]string {
	config_vars := make([]string, 0, len(config))
	for k := range config {
		config_vars = append(config_vars, k)
	}

	for _, v := range config_vars {
		if e := os.Getenv(v); e != "" {
			config[v] = e
		}
	}

	return config
}
