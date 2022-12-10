package netlify

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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
	token string
}

func NewNetlifyClient(token string) *NetlifyClient {
	return &NetlifyClient{token: token}
}

func (n NetlifyClient) httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ignore expired SSL certificates
			},
		},
	}
}

// Returns the DNS zone ID for the given hostname
func (n NetlifyClient) GetDnsZone(hostname string) string {
	var dns_zones DnsZones
	url_parts := strings.Split(hostname, ".")
	if len(url_parts) < 2 {
		msg := "Invalid hostname: " + hostname
		logrus.Error(msg, "get_dns_zone")
		os.Exit(1)
	}
	target := url_parts[len(url_parts)-2] + "." + url_parts[len(url_parts)-1]
	url := "https://api.netlify.com/api/v1/dns_zones"
	bearer := "Bearer " + n.token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}

	req.Header.Add("Authorization", bearer)
	// client := &http.Client{Transport: transCfg}

	resp, err := n.httpClient().Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := "Failed to read body: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	if err := json.Unmarshal(body, &dns_zones); err != nil {
		msg := "Failed to unmarshal JSON: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	for _, zone := range dns_zones {
		if zone.Name == target {
			return zone.ID
		}
	}
	return ""
}

// Return the DNS record ID for the given DNS zone and hostname
func (n NetlifyClient) GetDnsRecord(zone_id string, hostname string) (string, string) {
	var dns_records DnsRecords
	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records"
	bearer := "Bearer " + n.token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}

	req.Header.Add("Authorization", bearer)

	resp, err := n.httpClient().Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := "Failed to read body: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	if err := json.Unmarshal(body, &dns_records); err != nil {
		msg := "Failed to unmarshal JSON: " + err.Error()
		logrus.Error(msg, "get_dns_zone")
	}
	for _, record := range dns_records {
		if record.Hostname == hostname {
			return record.ID, record.Value
		}
	}
	return "", ""
}

// Deletes the given DNS record
func (n NetlifyClient) DeleteDnsRecord(zone_id string, record_id string) bool {
	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records/" + record_id
	bearer := "Bearer " + n.token
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		logrus.Error(msg, "delete_dns_record")
	}

	req.Header.Add("Authorization", bearer)

	resp, err := n.httpClient().Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		logrus.Error(msg, "delete_dns_record")
	}
	return resp.StatusCode == 204
}

// Creates a new DNS record
func (n NetlifyClient) CreateDnsRecord(zone_id string, hostname string, ip string) bool {
	new_dns := CreateDnsRecord{
		Type:     "A",
		Hostname: hostname,
		Value:    ip,
		TTL:      3600,
	}
	data, err := json.Marshal(new_dns)
	if err != nil {
		msg := "Failed to marshal DNS JSON: " + err.Error()
		logrus.Error(msg, "create_dns_record")
	}
	json_data := bytes.NewBuffer(data)

	url := "https://api.netlify.com/api/v1/dns_zones/" + zone_id + "/dns_records"
	bearer := "Bearer " + n.token
	req, err := http.NewRequest("POST", url, json_data)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		logrus.Error(msg, "create_dns_record")
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", CT)

	resp, err := n.httpClient().Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		logrus.Error(msg, "delete_dns_record")
	}
	return resp.StatusCode == 201
}
