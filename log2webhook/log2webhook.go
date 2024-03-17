package log2webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Log struct {
	File     string `json:"file"`
	Function string `json:"func"`
	Level    string `json:"level"`
	Message  string `json:"msg"`
	Time     string `json:"time"`
}

type Message struct {
	Embeds []Embed `json:"embeds"`
}
type Embed struct {
	Author      Author `json:"author"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Footer      Footer `json:"footer"`
	Timestamp   string `json:"timestamp"`
}

type Author struct {
	Name string `json:"name"`
}
type Footer struct {
	Text string `json:"text"`
}

type WebhookWriter struct{}

func (cw *WebhookWriter) Write(p []byte) (n int, err error) {
	webhook(p)
	return len(p), nil
}

func webhook(p []byte) {
	webhookURL := os.Getenv("LOG_WEBHOOK")
	if webhookURL == "" || !strings.HasPrefix(webhookURL, "http://") && !strings.HasPrefix(webhookURL, "https://") {
		return
	}
	var logJson Log
	json.Unmarshal(p, &logJson)
	var color string = "36314"
	if logJson.Level == "Error" {
		color = "16739179"
	}
	m := Message{
		Embeds: []Embed{
			{
				Author: Author{
					Name: logJson.File,
				},
				Title:       logJson.Function,
				Color:       color,
				Description: logJson.Message,
				Footer: Footer{
					Text: logJson.Level,
				},
				Timestamp: logJson.Time,
			},
		},
	}
	messageJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(messageJson))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}
