package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/aceaura/libra/core/net/http"
)

func TestCollector(t *testing.T) {
	const (
		url      = "https://top.baidu.com/board?platform=pc&sa=pcindex_entry"
		interval = 10 * time.Millisecond
	)
	c := http.NewCollector(
		http.CollectorOption.Background(),
		http.CollectorOption.Safety(),
		http.CollectorOption.Context(context.Background()),
		http.CollectorOption.Name("test_collector"),
		http.CollectorOption.RequestBacklog(1000),
		http.CollectorOption.ResponseBacklog(1000),
		http.CollectorOption.TPSLimit(1000),
		http.CollectorOption.ParallelIncrease(1),
		http.CollectorOption.ParallelTick(100*time.Millisecond),
		http.CollectorOption.ParallelIncrease(1),
	)
	if err := c.Serve(); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}

	timer := time.NewTicker(interval).C

	for {
		select {
		case <-timer:
			req := &http.ServiceRequest{
				URL:     url,
				Timeout: 10 * time.Second,
				Retry:   3,
			}
			c.RequestChan() <- req
		case resp := <-c.ResponseChan():
			if len(resp.Body) == 0 {
				t.Fatal("expected a body with content")
			}
		}
	}
}
