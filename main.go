package main

import (
	"os"
	"strconv"
	"time"

	"log/syslog"

	"github.com/sirupsen/logrus"

	netlify "github.com/gntouts/ipsync/pkg/netlify"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

const TIMEOUT = 20 // seconds

type Config struct {
	Token   string
	Dns     string
	Timeout int
}

func getConfig() Config {
	token := os.Getenv("NETLIFY_TOKEN")
	if token == "" {
		logrus.Error("NETLIFY_TOKEN not set")
		os.Exit(1)
	}
	dns_target := os.Getenv("DNS_TARGET")
	if dns_target == "" {
		logrus.Error("DNS_TARGET not set", "get_dns_target")
		os.Exit(1)
	}
	timeout := os.Getenv("IPSYNC_TIMEOUT")
	intVar, err := strconv.Atoi(timeout)
	if err != nil {
		intVar = 20
	}
	return Config{
		Token:   token,
		Dns:     dns_target,
		Timeout: intVar,
	}
}

func main() {
	// configure logurs
	formatter := new(logrus.TextFormatter)
	formatter.DisableColors = true

	logrus.SetFormatter(formatter)
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err == nil {
		logrus.AddHook(hook)
	}
	logrus.Info("Started monitoring IP address", "main")

	// retrieve config from ENV
	config := getConfig()
	target := config.Dns

	netlify := netlify.NewNetlifyClient(config.Token)

	zone := netlify.GetDnsZone(target)
	record_id, record_ip := netlify.GetDnsRecord(zone, target)
	logrus.Info("Netlify IP is set "+record_ip, "main")

	for {
		current, err := GetIp()
		if err != nil {
			logrus.Error(err.Error())
		}

		if current.Ip != record_ip {
			netlify.DeleteDnsRecord(zone, record_id)
			changed := netlify.CreateDnsRecord(zone, target, current.Ip)
			logrus.Info("Local IP changed to "+current.Ip, "main")
			var msg string
			if changed {
				msg = "Updated IP address to " + current.Ip
			} else {
				msg = "Failed to update IP address to " + current.Ip
			}
			logrus.Info(msg, "main")
			record_id, record_ip = netlify.GetDnsRecord(zone, target)
		}
		time.Sleep(time.Duration(config.Timeout) * time.Second)
	}
}
