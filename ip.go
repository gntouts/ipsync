package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type IpInfo struct {
	Ip string
	Ts int64
}

type IpResponse struct {
	IP string `json:"ip"`
}

func GetIp() (IpInfo, error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err.Error())
		}
	}()

	res, err := http.Get("https://ip4.seeip.org/json")
	if err != nil {
		return IpInfo{}, errors.New("Failed to send request: " + err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return IpInfo{}, errors.New("Failed to read body: " + err.Error())
	}
	var result IpResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return IpInfo{}, errors.New("Failed to unmarshal JSON: " + err.Error())
	}

	sec := time.Now().Unix()
	return IpInfo{result.IP, sec}, err
}
