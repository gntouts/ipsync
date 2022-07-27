package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const CT = "application/json"

type DnsZone struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name"`
	Errors               []interface{} `json:"errors"`
	SupportedRecordTypes []string      `json:"supported_record_types"`
	UserID               string        `json:"user_id"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	Records              []interface{} `json:"records"`
	DNSServers           []string      `json:"dns_servers"`
	AccountID            string        `json:"account_id"`
	SiteID               interface{}   `json:"site_id"`
	AccountSlug          string        `json:"account_slug"`
	AccountName          string        `json:"account_name"`
	Domain               interface{}   `json:"domain"`
	Ipv6Enabled          bool          `json:"ipv6_enabled"`
	Dedicated            interface{}   `json:"dedicated"`
}

type DnsZones []struct {
	DnsZone
}

type DnsRecord struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	TTL       int    `json:"ttl"`
	Priority  int    `json:"priority"`
	DNSZoneID string `json:"dns_zone_id"`
	SiteID    string `json:"site_id"`
	Flag      int    `json:"flag"`
	Tag       string `json:"tag"`
	Managed   bool   `json:"managed"`
}

type DnsRecords []struct {
	DnsRecord
}

type CreateDnsRecord struct {
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}

type NetlifyClient struct {
}

func transport_config() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
}

func NewNetlifyClient() *NetlifyClient {
	return &NetlifyClient{}
}

// Returns the DNS zone ID for the given hostname
func (NetlifyClient) get_dns_zone(hostname string) string {
	var dns_zones DnsZones
	url_parts := strings.Split(hostname, ".")
	if len(url_parts) < 2 {
		msg := "Invalid hostname: " + hostname
		log_err(msg, "get_dns_zone")
		os.Exit(1)
	}
	target := url_parts[len(url_parts)-2] + "." + url_parts[len(url_parts)-1]
	url := "https://api.netlify.com/api/v1/dns_zones"
	bearer := "Bearer " + config().get_access_token()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		log_err(msg, "get_dns_zone")
	}

	req.Header.Add("Authorization", bearer)
	client := &http.Client{Transport: transport_config()}
	// client := &http.Client{Transport: transCfg}

	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := "Failed to read body: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	if err := json.Unmarshal(body, &dns_zones); err != nil {
		msg := "Failed to unmarshal JSON: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	for _, zone := range dns_zones {
		if zone.Name == target {
			return zone.ID
		}
	}
	return ""
}

// Return the DNS record ID for the given DNS zone and hostname
func (NetlifyClient) get_dns_record(zone_id string, hostname string) (string, string) {
	var dns_records DnsRecords
	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records"
	bearer := "Bearer " + config().get_access_token()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		log_err(msg, "get_dns_zone")
	}

	req.Header.Add("Authorization", bearer)
	client := &http.Client{Transport: transport_config()}

	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := "Failed to read body: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	if err := json.Unmarshal(body, &dns_records); err != nil {
		msg := "Failed to unmarshal JSON: " + err.Error()
		log_err(msg, "get_dns_zone")
	}
	for _, record := range dns_records {
		if record.Hostname == hostname {
			return record.ID, record.Value
		}
	}
	return "", ""
}

// Deletes the given DNS record
func (NetlifyClient) delete_dns_record(zone_id string, record_id string) bool {
	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records/" + record_id
	bearer := "Bearer " + config().get_access_token()
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		log_err(msg, "delete_dns_record")
	}

	req.Header.Add("Authorization", bearer)
	client := &http.Client{Transport: transport_config()}

	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		log_err(msg, "delete_dns_record")
	}
	return resp.StatusCode == 204
}

// Creates a new DNS record
func (NetlifyClient) create_dns_record(zone_id string, hostname string, ip string) bool {
	new_dns := CreateDnsRecord{
		Type:     "A",
		Hostname: hostname,
		Value:    ip,
		TTL:      3600,
	}
	data, err := json.Marshal(new_dns)
	if err != nil {
		msg := "Failed to marshal DNS JSON: " + err.Error()
		log_err(msg, "create_dns_record")
	}
	json_data := bytes.NewBuffer(data)

	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records"
	bearer := "Bearer " + config().get_access_token()
	req, err := http.NewRequest("POST", url, json_data)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		log_err(msg, "create_dns_record")
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", CT)
	client := &http.Client{Transport: transport_config()}
	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		log_err(msg, "delete_dns_record")
	}
	return resp.StatusCode == 201
}
