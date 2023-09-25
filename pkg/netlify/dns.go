package netlify

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ContentType = "application/json"

type DNSZone struct {
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

type DNSRecord struct {
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

type createDNSRecordReq struct {
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}
type simpleResponse struct {
	body   []byte
	status int
}
type NetlifyClient struct {
	token      string
	httpClient *http.Client
}

func NewNetlifyClient(token string) *NetlifyClient {
	temp := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ignore expired SSL certificates
			},
		}}
	return &NetlifyClient{token: token, httpClient: temp}
}

func (n NetlifyClient) doRequest(httpMethod string, targetURL string, b []byte) (r simpleResponse, err error) {
	bearer := "Bearer " + n.token
	request, err := http.NewRequest(httpMethod, targetURL, bytes.NewBuffer(b))
	if err != nil {
		return r, err
	}
	request.Header.Add("Authorization", bearer)
	request.Header.Set("Content-Type", ContentType)
	response, err := n.httpClient.Do(request)
	if err != nil {
		return r, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return r, fmt.Errorf("failed to read body: %s", err.Error())
	}
	return simpleResponse{body: body, status: response.StatusCode}, nil
}

func (n NetlifyClient) GetDNSZones() (dnsZones []DNSZone, err error) {
	response, err := n.doRequest("GET", "https://api.netlify.com/api/v1/dns_zones", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response.body, &dnsZones)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response to DNS zones list: %s", err.Error())
	}
	return dnsZones, nil
}

func (n NetlifyClient) GetDNSRecords(zoneID string) ([]DNSRecord, error) {
	var records []DNSRecord
	url := "https://api.netlify.com/api/v1/dns_zones/" + zoneID + "/dns_records"
	response, err := n.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response.body, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response to DNS records: %s", err.Error())
	}
	return records, nil
}

func (n NetlifyClient) GetDNSRecord(zoneID string, hostname string) (DNSRecord, error) {
	records, err := n.GetDNSRecords(zoneID)
	if err != nil {
		return DNSRecord{}, err
	}
	for _, r := range records {
		if r.Hostname == hostname {
			return r, nil
		}
	}
	return DNSRecord{}, fmt.Errorf("failed to find requested record \"%s\" in DNS zone \"%s\" response to DNS records", hostname, zoneID)
}

func (n NetlifyClient) DeleteDNSRecord(zoneID string, recordID string) error {
	url := "https://api.netlify.com/api/v1/dns_zones/" + zoneID + "/dns_records/" + recordID
	response, err := n.doRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to delete requested record \"%s\": %s", recordID, err.Error())
	}
	if response.status == 204 {
		return nil
	}
	return fmt.Errorf("failed to delete requested record \"%s\": Expected HTTP 204, got %d", recordID, response.status)
}

func (n NetlifyClient) CreateADNSRecord(zoneID string, hostname string, IPaddress string) error {
	requestBody := createDNSRecordReq{
		Type:     "A",
		Hostname: hostname,
		Value:    IPaddress,
		TTL:      3600,
	}
	data, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Failed to marshal new DNS record to JSON: " + err.Error())
	}
	url := "https://api.netlify.com/api/v1/dns_zones/" + zoneID + "/dns_records"
	response, err := n.doRequest("POST", url, data)
	if err != nil {
		return err
	}
	if response.status == 201 {
		return nil
	}
	return fmt.Errorf("failed to create requested record for \"%s\": Expected HTTP 201, got %d", hostname, response.status)
}
