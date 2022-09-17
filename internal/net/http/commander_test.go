package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/cloudlibraries/libra/internal/core/device"
	"github.com/cloudlibraries/libra/internal/net/http"
)

func TestComannder(t *testing.T) {
	const (
		url      = "https://www.baidu.com"
		interval = 10 * time.Millisecond
	)
	c := http.NewCommander(
		http.WithCommanderBackground(),
		http.WithCommanderSafety(),
		http.WithCommanderContext(context.Background()),
		http.WithCommanderName("HttpCollector"),
		http.WithCommanderRequestBacklog(1000),
		http.WithCommanderResponseBacklog(1000),
		http.WithCommanderReportBacklog(1),
		http.WithCommanderTPSLimit(20),
		http.WithCommanderParallel(10),
		http.WithCommanderParallelTick(100*time.Millisecond),
		http.WithCommanderParallelIncrease(1),
	)
	device.Bus().Integrate(c)
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

func TestInvoke(t *testing.T) {
	const (
		url      = "https://www.baidu.com"
		interval = 10 * time.Millisecond
	)
	c := http.NewCommander(
		http.WithCommanderBackground(),
		http.WithCommanderSafety(),
		http.WithCommanderContext(context.Background()),
		http.WithCommanderName("HttpCollector"),
		http.WithCommanderRequestBacklog(1000),
		http.WithCommanderResponseBacklog(1000),
		http.WithCommanderReportBacklog(1),
		http.WithCommanderTPSLimit(20),
		http.WithCommanderParallel(10),
		http.WithCommanderParallelTick(100*time.Millisecond),
		http.WithCommanderParallelIncrease(1),
	)
	device.Bus().Integrate(c)
	if err := c.Serve(); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}

	req := &http.ServiceRequest{
		URL:     url,
		Timeout: 10 * time.Second,
		Retry:   3,
	}

	if resp := c.Invoke(req); resp.Err != nil {
		t.Fatal(resp.Err)
	} else {
		t.Log(resp.Body)
	}
}
