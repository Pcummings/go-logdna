package logdna

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	protocol = "https://"
	domain   = "logs.logdna.com"
	endPoint = "logs/ingest"
)

type LogdnaConfig struct {
	IngestionKey string
	LogLevel     []string
	Tags         []string
	AppName      string
	Environment  string
}
type Client struct {
	Config *LogdnaConfig
}
type Payload struct {
	Line  string
	App   string
	Level string
	Env   string
}

// NewClient returns a client for Logdna
func NewClient(config *LogdnaConfig) (*Client, error) {
	client := new(Client)
	if config.IngestionKey == "" {
		config.IngestionKey = os.Getenv("LOGDNA_KEY")
	}
	if config.IngestionKey == "" {
		fmt.Println("LogDNA Ingestion Key not provided")
		os.Exit(0)
	}
	client.Config = config
	return client, nil
}
func (c *Client) Info(msg string) {
	payload := c.configurePayload()
	payload.Level = "INFO"
	payload.Line = msg
	if contains(c.Config.Tags, "INFO") {
		c.do(payload)
	}
}
func (c *Client) Notice(msg string) {
	payload := c.configurePayload()
	payload.Level = "NOTICE"
	payload.Line = msg
	if contains(c.Config.Tags, "NOTICE") {
		c.do(payload)
	}
}
func (c *Client) Warn(msg string) {
	payload := c.configurePayload()
	payload.Level = "WARN"
	payload.Line = msg
	if contains(c.Config.Tags, "WARN") {
		c.do(payload)
	}
}
func (c *Client) Error(msg string) {
	payload := c.configurePayload()
	payload.Level = "ERROR"
	payload.Line = msg
	if contains(c.Config.Tags, "ERROR") {
		c.do(payload)
	}
}
func (c *Client) Fatal(msg string) {
	payload := c.configurePayload()
	payload.Level = "FATAL"
	payload.Line = msg
	if contains(c.Config.Tags, "FATAL") {
		c.do(payload)
	}
}
func (c *Client) configurePayload() *Payload {
	payload := new(Payload)
	payload.App = c.Config.AppName
	payload.Env = c.Config.Environment
	return payload
}
func (c *Client) do(payload *Payload) {
	u := &url.URL{
		Scheme: protocol,
		Host:   domain,
		Path:   fmt.Sprintf("/%s", endPoint),
	}
	q := u.Query()
	q.Set("hostname", getHostName())
	q.Set("mac", getMacAddr())
	q.Set("ip", getIpAddr())
	q.Set("now", strconv.Itoa(int(time.Now().Unix())))
	if len(c.Config.Tags) > 0 {
		q.Set("tags", strings.Join(c.Config.Tags, ","))
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		//handle error
	}
	body := bytes.NewReader(payloadBytes)
	r, err := http.NewRequest("POST", u.String(), body)
	r.SetBasicAuth(c.Config.IngestionKey, "")
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		//handle error
	}
	defer resp.Body.Close()
}
