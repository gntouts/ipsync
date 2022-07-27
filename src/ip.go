package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type IpClient struct{}

type IpInfo struct {
	Ip string `json:"Ip"`
	Ts int64  `json:"Ts"`
}

type IpResponse struct {
	IP string `json:"ip"`
}

func NewIpClient() *IpClient {
	return &IpClient{}
}

func (IpClient) get_ip() IpInfo {
	var result IpResponse
	var err error
	req, err := http.NewRequest("GET", "https://ip4.seeip.org/json", nil)
	if err != nil {
		msg := "Failed to create request: " + err.Error()
		log_err(msg, "get_ip")
	}

	client := &http.Client{Transport: transport_config()}
	res, err := client.Do(req)
	if err != nil {
		msg := "Failed to send request: " + err.Error()
		log_err(msg, "get_ip")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log_info("response.body: "+string(body), "get_ip")

	if err != nil {
		msg := "Failed to read body: " + err.Error()
		log_err(msg, "get_ip")
	}
	if err := json.Unmarshal(body, &result); err != nil {
		msg := "Failed to unmarshal JSON: " + err.Error()
		log_err(msg, "get_ip")
	}
	if err != nil {
		return IpInfo{}
	}
	sec := time.Now().Unix()
	return IpInfo{result.IP, sec}
}
