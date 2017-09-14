/**
 * Copyright (c) 2016 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package cmd

import (
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"

	"github.com/mainflux/mainflux-core/models"
)

// CreateChannel - creates new channel and generates UUID
func CreateChannel(msg string) string {
	url := UrlHTTP + "/channels"
	sr := strings.NewReader(msg)
	resp, err := netClient.Post(url, "application/json", sr)
	s := PrettyHttpResp(resp, err)

	return s
}

// GetChannels - gets all channels
func GetChannels(limit int) string {
	url := UrlHTTP + "/channels?climit=" + strconv.Itoa(limit)
	resp, err := netClient.Get(url)
	s := PrettyHttpResp(resp, err)

	return s
}

// GetChannel - gets channel by ID
func GetChannel(id string) string {
	url := UrlHTTP + "/channels/" + id
	s := PrettyHttpResp(netClient.Get(url))

	return s
}

// UpdateChannel - publishes SenML message on the channel
func UpdateChannel(id string, msg string) string {
	url := UrlHTTP + "/channels/" + id
	sr := strings.NewReader(msg)
	req, err := http.NewRequest("PUT", url, sr)
	if err != nil {
		return err.Error()
	}

	req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(msg)))

	s := PrettyHttpResp(netClient.Do(req))

	return s
}

// DeleteChannel - removes channel
func DeleteChannel(id string) string {
	url := UrlHTTP + "/channels/" + id
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err.Error()
	}
	s := PrettyHttpResp(netClient.Do(req))

	return s
}

// DeleteAllChannels - removes all channels
func DeleteAllChannels() string {
	url := UrlHTTP + "/channels"
	resp, _ := netClient.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)

	var channels []models.Channel
	json.Unmarshal([]byte(body), &channels)
	s := ""
	for i := 0; i < len(channels); i++ {
		s = s + DeleteChannel(channels[i].ID) + "\n\n"
	}

	return s
}

// PlugChannel - plugs list of devices into the channel
func PlugChannel(id string, devices string) string {
	url := UrlHTTP + "/channels/" + id + "/plug"
	sr := strings.NewReader(devices)

	s := PrettyHttpResp(netClient.Post(url, "application/json", sr))

	return s
}

// UnplugChannel - unplugs list of devices from the channel
func UnplugChannel(id string, devices string) string {
	url := UrlHTTP + "/channels/" + id + "/unplug"
	sr := strings.NewReader(devices)

	s := PrettyHttpResp(netClient.Post(url, "application/json", sr))

	return s
}
