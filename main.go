package main

import (
	"time"

	"github.com/gntouts/ipsync/pkg/netlify"
	"github.com/sirupsen/logrus"
)

func main() {
	configLogrus()
	config, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration from ENV variables")
	}
	logrus.WithFields(logrus.Fields{
		"NETLIFY_TOKEN":  "REDACTED",
		"DNS_TARGET":     config.DNSTarget,
		"IPSYNC_TIMEOUT": config.Timeout,
	}).Info("Loaded config from ENV variables")
	ntlfClient := netlify.NewNetlifyClient(config.NetlifyToken)
	zones, err := ntlfClient.GetDNSZones()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get Netlify DNS Zones")
	}
	var targetZone netlify.DNSZone
	found := false
	targetTLD, err := getTLD(config.DNSTarget)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get TLD from provided DNS target")
	}
	for _, zone := range zones {
		if zone.Name == targetTLD {
			targetZone = zone
			found = true
			break
		}
	}
	if !found {
		logrus.Fatal("DNS target not found in available DNS zones")
	}
	logrus.WithField("DNS zone", targetZone.ID).Info("Found DNS zone")
	for {
		dnsRecord, err := ntlfClient.GetDNSRecord(targetZone.ID, config.DNSTarget)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get DNS record for provided DNS target")
		}
		logrus.WithField("DNS record hostname", dnsRecord.Hostname).WithField("DNS record value", dnsRecord.Value).Info("Found DNS record")

		myIP, err := getPublicIP()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get public IP")
		}
		logrus.WithField("current IP", myIP).Info("Found current IP")
		if myIP != dnsRecord.Value {
			logrus.WithField("current IP", myIP).WithField("DNS IP", dnsRecord.Value).Info("DNS record mismatch")
			err := ntlfClient.DeleteDNSRecord(targetZone.ID, dnsRecord.ID)
			if err != nil {
				logrus.WithError(err).Fatal("Failed to delete DNS record")
			}
			logrus.Info("Deleted old DNS record")
			err = ntlfClient.CreateADNSRecord(targetZone.ID, config.DNSTarget, myIP)
			if err != nil {
				logrus.WithError(err).Fatal("Failed to create updated DNS record")
			}
			logrus.Info("Created updated DNS record")

		}
		logrus.Infof("Check complete, sleeping for %d seconds", config.Timeout)
		time.Sleep(time.Second * time.Duration(config.Timeout))
	}
}
