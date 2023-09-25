package main

import (
	"errors"
	"fmt"
	"io"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

const DefaultTimeout = 300

type IPSyncConfig struct {
	NetlifyToken string
	DNSTarget    string
	Timeout      int
}

func loadConfig() (*IPSyncConfig, error) {
	conf := &IPSyncConfig{
		NetlifyToken: os.Getenv("NETLIFY_TOKEN"),
		DNSTarget:    os.Getenv("DNS_TARGET"),
	}
	timeout := os.Getenv("IPSYNC_TIMEOUT")
	intTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		logrus.WithField("timeout", timeout).Warn("Invalid timeout")
		intTimeout = DefaultTimeout
	}
	conf.Timeout = intTimeout
	if conf.NetlifyToken == "" {
		logrus.WithField("netlify token", conf.NetlifyToken).Error("Invalid netlify token")
		return nil, errors.New("netlify token is empty")
	}
	if conf.DNSTarget == "" {
		logrus.WithField("DNS target", conf.NetlifyToken).Error("Invalid DNS target")
		return nil, errors.New("DNS target is empty")
	}
	return conf, nil
}

func configLogrus() {
	formatter := new(logrus.TextFormatter)
	formatter.DisableColors = true
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.InfoLevel)
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err == nil {
		logrus.AddHook(hook)
	}
}

func getTLD(url string) (string, error) {
	parts := strings.Split(url, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid url")
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
}

func getPublicIP() (string, error) {
	url := "http://api.ipify.org"
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(body), nil
}
