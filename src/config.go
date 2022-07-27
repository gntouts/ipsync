package main

import (
	"os"
	"strconv"
)

type Config struct {
}

func config() Config {
	return Config{}
}

func (Config) get_access_token() string {
	token := os.Getenv("NETLIFY_TOKEN")
	if token == "" {
		log_err("NETLIFY_TOKEN not set", "get_access_token")
		os.Exit(1)
	}
	return token
}

func (Config) get_dns_target() string {
	dns_target := os.Getenv("DNS_TARGET")
	if dns_target == "" {
		log_err("DNS_TARGET not set", "get_dns_target")
		os.Exit(1)
	}
	return dns_target
}

func (Config) get_timeout() int {
	timeout := os.Getenv("IPSYNC_TIMEOUT")
	intVar, err := strconv.Atoi(timeout)
	if err != nil {
		return 20
	}
	return intVar
}
