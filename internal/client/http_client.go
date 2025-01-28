package client

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amavrin/go-musthave-devops/internal/metrics"
)

const (
	defaultTimeout  = 2 * time.Second
	reportInterval  = 10 * time.Second
	numberOfMetrics = 32
)

type Client struct {
	URL    string
	Client *http.Client
}

func NewClient(URL string) *Client {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &Client{
		URL:    URL,
		Client: client,
	}
}

func (c *Client) SendLoop(db *metrics.DB) {
	for {
		time.Sleep(reportInterval)
		m := db.GetMetrics()
		err := c.sendMetrics(m)
		if err != nil {
			log.Printf("Error sending metrics: %v", err)
		}
	}
}

func (c *Client) sendMetrics(m metrics.Metrics) error {
	countErr := 0
	urls := c.updateURLs(m)
	for _, url := range urls {
		err := c.send(url)
		if err != nil {
			log.Printf("Error sending metrics: %v", err)
			countErr++
			if countErr > 10 {
				return fmt.Errorf("too many consecutive errors sending metrics: %v", err)
			}
			continue
		}
		countErr = 0
	}
	return nil
}

func (c *Client) gaugeURL(name string, value metrics.Gauge) string {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", c.URL, name, value)
	return url
}

func (c *Client) counterURL(name string, value metrics.Counter) string {
	url := fmt.Sprintf("%s/update/counter/%s/%d", c.URL, name, value)
	return url
}

func (c *Client) updateURLs(m metrics.Metrics) []string {
	urls := make([]string, 0, numberOfMetrics)
	urls = append(urls, c.gaugeURL("Alloc", m.Alloc))
	urls = append(urls, c.gaugeURL("BuckHashSys", m.BuckHashSys))
	urls = append(urls, c.gaugeURL("Frees", m.Frees))
	urls = append(urls, c.gaugeURL("GCCPUFraction", m.GCCPUFraction))
	urls = append(urls, c.gaugeURL("GCSys", m.GCSys))
	urls = append(urls, c.gaugeURL("HeapAlloc", m.HeapAlloc))
	urls = append(urls, c.gaugeURL("HeapIdle", m.HeapIdle))
	urls = append(urls, c.gaugeURL("HeapInuse", m.HeapInuse))
	urls = append(urls, c.gaugeURL("HeapObjects", m.HeapObjects))
	urls = append(urls, c.gaugeURL("HeapReleased", m.HeapReleased))
	urls = append(urls, c.gaugeURL("HeapSys", m.HeapSys))
	urls = append(urls, c.gaugeURL("LastGC", m.LastGC))
	urls = append(urls, c.gaugeURL("Lookups", m.Lookups))
	urls = append(urls, c.gaugeURL("MCacheInuse", m.MCacheInuse))
	urls = append(urls, c.gaugeURL("MCacheSys", m.MCacheSys))
	urls = append(urls, c.gaugeURL("MSpanInuse", m.MSpanInuse))
	urls = append(urls, c.gaugeURL("MSpanSys", m.MSpanSys))
	urls = append(urls, c.gaugeURL("Mallocs", m.Mallocs))
	urls = append(urls, c.gaugeURL("NextGC", m.NextGC))
	urls = append(urls, c.gaugeURL("NumForcedGC", m.NumForcedGC))
	urls = append(urls, c.gaugeURL("NumGC", m.NumGC))
	urls = append(urls, c.gaugeURL("OtherSys", m.OtherSys))
	urls = append(urls, c.gaugeURL("PauseTotalNs", m.PauseTotalNs))
	urls = append(urls, c.gaugeURL("StackInuse", m.StackInuse))
	urls = append(urls, c.gaugeURL("StackSys", m.StackSys))
	urls = append(urls, c.gaugeURL("Sys", m.Sys))
	urls = append(urls, c.gaugeURL("TotalAlloc", m.TotalAlloc))
	urls = append(urls, c.counterURL("PollCounter", m.PollCounter))
	urls = append(urls, c.gaugeURL("RandomValue", m.RandomValue))
	return urls
}

func (c *Client) send(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST %s: unexpected status code: %d", url, resp.StatusCode)
	}
	return nil
}
