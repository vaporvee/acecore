package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func simpleGetFromAPI(key string, url string) interface{} {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request:", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error making request:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error("Error decoding JSON:", err)
	}
	return result[key]
}
