package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/net/http"
)

func TestCollector(t *testing.T) {
	const (
		url      = "https://www.baidu.com"
		interval = 10 * time.Millisecond
	)
	c := http.NewCommander(
		http.CommanderOption.Background(),
		http.CommanderOption.Safety(),
		http.CommanderOption.Context(context.Background()),
		http.CommanderOption.Name("HttpCollector"),
		http.CommanderOption.RequestBacklog(1000),
		http.CommanderOption.ResponseBacklog(1000),
		http.CommanderOption.ReportBacklog(1),
		http.CommanderOption.TPSLimit(20),
		http.CommanderOption.ParallelInit(10),
		http.CommanderOption.ParallelTick(100*time.Millisecond),
		http.CommanderOption.ParallelIncrease(1),
	)
	device.Bus().WithDevice(c)
	if err := c.Serve(); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}

	timer := time.NewTicker(interval).C
	req := &http.ServiceRequest{
		URL:     url,
		Timeout: 10 * time.Second,
		Retry:   3,
	}

	var count = 0
	for {
		select {
		case <-timer:
			c.RequestChan() <- req
		case resp := <-c.ResponseChan():
			if len(resp.Body) == 0 {
				t.Fatal("expected a body with content")
			}
			count++
			if count%100 == 0 {
				t.Logf("finish count:%d", count)
			}
			if count == 500 {
				return
			}
		}
	}
}
