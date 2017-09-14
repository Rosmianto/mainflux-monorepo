/**
 * Copyright (c) 2016 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

// CreateUser - create user
func CreateUser(user string, pwd string) string {
	var err error

	msg := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pwd)
	sr := strings.NewReader(msg)
	url := UrlAuth + "/users"
	resp, err := netClient.Post(url, "application/json", sr)
	s := PrettyHttpResp(resp, err)

	return s
}

// LogginInUser - loggin in user
func LogginInUser(user string, pwd string) string {
	var err error

	msg := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pwd)
	sr := strings.NewReader(msg)
	url := UrlAuth + "/sessions"
	resp, err := netClient.Post(url, "application/json", sr)
	s := PrettyHttpResp(resp, err)

	return s
}

// GetApiKeys - Retrieved a list of created keys
func GetApiKeys(auth string) string {
	url := UrlAuth + "/api-keys"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	s := PrettyHttpResp(resp, err)

	return s
}

// CreateApiKeys - An API key can be given to the user, device or channel
func CreateApiKeys(auth string, owner string) string {
	url := UrlAuth + "/api-keys"

	req, err := http.NewRequest("POST", url, strings.NewReader(owner))
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	s := PrettyHttpResp(resp, err)

	return s
}

// DeleteApiKeys - Completely removes the key from the API key list
func DeleteApiKeys(auth string, key string) string {
	url := UrlAuth + "/api-keys/" + key

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	s := PrettyHttpResp(resp, err)

	return s
}

// GetOwnerApiKeys - Get key owner
func GetOwnerApiKeys(auth string, key string) string {
	url := UrlAuth + "/api-keys/" + key
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	s := PrettyHttpResp(resp, err)

	return s
}

// UpdateOwnerApiKeys - Updates the key scope
func UpdateOwnerApiKeys(auth string, key string, owner string) string {
	url := UrlAuth + "/api-keys/" + key
	req, err := http.NewRequest("PUT", url, strings.NewReader(owner))
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
	req.Header.Add("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	s := PrettyHttpResp(resp, err)

	return s
}
