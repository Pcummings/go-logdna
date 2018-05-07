package logdna

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

const IngestBaseURL = "https://logs.logdna.com/logs/ingest"

type LogdnaConfig struct {
	IngestionKey string
	LogLevel     []string
	Tags         []string
	AppName      string
	Environment  string
}
type Client struct {
	config  *LogdnaConfig
	apiUrl  url.URL
	payload payloadJSON
}
type LineJSON struct {
	Msg        string `json:"message"`
	File       string `json:"file"`
	Linenumber string `json:"linenumber"`
}
type logLineJSON struct {
	Line  string `json:"line"`
	App   string `json:"app"`
	Level string `json:"level"`
	Env   string `json:"env"`
}
type payloadJSON struct {
	Lines []logLineJSON `json:"lines"`
}

func NewClient(config *LogdnaConfig) (*Client, error) {
	var client Client
	client.apiUrl = makeIngestURL(config)
	if len(config.LogLevel) == 0 {
		config.LogLevel = []string{"INFO", "WARN", "ERROR", "DEBUG"}
	}
	client.config = config
	return &client, nil
}

func (c *Client) Info(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	c.configurePayload(msg, filename, strconv.Itoa(line), "INFO")
	if contains(c.config.LogLevel, "INFO") {
		err := c.do()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Client) Warn(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	c.configurePayload(msg, filename, strconv.Itoa(line), "WARN")
	if contains(c.config.LogLevel, "WARN") {
		err := c.do()
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (c *Client) Error(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	c.configurePayload(msg, filename, strconv.Itoa(line), "ERROR")
	if contains(c.config.LogLevel, "ERROR") {
		err := c.do()
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (c *Client) Debug(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	c.configurePayload(msg, filename, strconv.Itoa(line), "DEBUG")
	if contains(c.config.LogLevel, "DEBUG") {
		err := c.do()
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (c *Client) configurePayload(msg string, file string, line string, level string) {
	message := &LineJSON{
		Msg:        msg,
		File:       file,
		Linenumber: line,
	}
	e, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	logLine := logLineJSON{
		Line:  string(e),
		App:   c.config.AppName,
		Env:   c.config.Environment,
		Level: level,
	}
	c.payload.Lines = append(c.payload.Lines, logLine)
}
func (c *Client) do() error {
	jsonPayload, err := json.Marshal(c.payload)
	if err != nil {
		return err
	}

	jsonReader := bytes.NewReader(jsonPayload)

	resp, err := http.Post(c.apiUrl.String(), "application/json", jsonReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.payload = payloadJSON{}
	return err
}
func makeIngestURL(cfg *LogdnaConfig) url.URL {
	u, _ := url.Parse(IngestBaseURL)
	if cfg.IngestionKey == "" {
		cfg.IngestionKey = os.Getenv("LOGDNA_KEY")
	}
	if cfg.IngestionKey == "" {
		fmt.Println("LogDNA Ingestion Key not provided")
		os.Exit(0)
	}
	u.User = url.User(cfg.IngestionKey)
	values := url.Values{}
	values.Set("hostname", getHostName())
	values.Set("mac", getMacAddr())
	values.Set("ip", getIpAddr())
	values.Set("now", strconv.FormatInt(time.Time{}.UnixNano(), 10))
	u.RawQuery = values.Encode()

	return *u
}
