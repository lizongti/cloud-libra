package http_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aceaura/libra/http"
)

func TestClientGet(t *testing.T) {
	resp, body, err := http.Get("www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientDoGet(t *testing.T) {
	resp, body, err := http.Do(http.GET, "www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientPost(t *testing.T) {
	resp, body, err := http.Post("www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientDoPost(t *testing.T) {
	resp, body, err := http.Do(http.POST, "www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientHead(t *testing.T) {
	resp, body, err := http.Head("www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) > 0 {
		t.Fatal("expected a body without content")
	}
	if len(resp.Header) == 0 {
		t.Fatal("expected a header with content")
	}
	t.Log(resp.Header)
}

func TestClientDoHead(t *testing.T) {
	resp, body, err := http.Do(http.HEAD, "www.baidu.com")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) > 0 {
		t.Fatal("expected a body without content")
	}
	if len(resp.Header) == 0 {
		t.Fatal("expected a header with content")
	}
	t.Log(resp.Header)
}

func TestClientParam(t *testing.T) {
	resp, body, err := http.Get("https://top.baidu.com/board",
		http.WithParam("platform", "pc"),
		http.WithParam("sa", "pcindex_entry"),
	)
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientForm(t *testing.T) {
	resp, body, err := http.NewClient().WithForm(url.Values{
		"platform": []string{"pc"},
		"sa":       []string{"pcindex_entry"},
	}).Get("https://top.baidu.com/board")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientProtocol(t *testing.T) {
	resp, body, err := http.NewClient().WithProtocol("https").Get("top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientTimeout(t *testing.T) {
	_, body, err := http.NewClient().WithTimeout(time.Duration(1) * time.Microsecond).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), "context deadline exceeded") < 0 {
		t.Fatal("expected an error with timeout")
	}
	t.Log(string(body))
}

func TestClientContentType(t *testing.T) {
	resp, body, err := http.NewClient().WithContentType("exception").WithRetry(3).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	contentType := resp.Request.Header["Content-Type"][0]
	if contentType != "exception" {
		t.Fatalf("unexpected content type getting from client: %s", contentType)
	}
	t.Log(string(body))
}

func TestClientResponseBodyReader(t *testing.T) {
	var body []byte
	var err error
	resp, _, err := http.NewClient().WithResponseBodyReader(func(r io.Reader) error {
		body, err = ioutil.ReadAll(r)
		return err
	}).WithRetry(3).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	t.Log(string(body))
}

func TestClientSafety(t *testing.T) {
	text := "safety exception"
	_, _, err := http.NewClient().WithResponseBodyReader(func(r io.Reader) error {
		panic(errors.New(text))
	}).WithClientSafety(true).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), text) < 0 {
		t.Fatal("expected an error with safety exception")
	}
}

func TestClientBody(t *testing.T) {
	http.Serve("localhost:1989",
		http.WithBackground(true),
		http.WithRoute("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprint(w, string(body))
		}))
	text := "this is a body text"
	resp, body, err := http.NewClient().WithBody(text).Get("localhost:1989")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
	if string(body) != text {
		t.Fatalf("unexpected body type getting from server: %s", string(body))
	}
}

func TestClientRetry(t *testing.T) {
	var count int = 1
	http.Serve("localhost:1989",
		http.WithBackground(true),
		http.WithRoute("/", func(w http.ResponseWriter, r *http.Request) {
			if count < 3 {
				count++
				time.Sleep(time.Second * 2)
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprint(w, string(body))
		}))
	text := "this is a body text"
	resp, body, err := http.NewClient().WithTimeout(time.Second * 1).WithRetry(3).WithBody(text).Get("localhost:1989")
	if err != nil {
		t.Fatalf("unexpected error getting from client: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected a status code of 200, got %v", resp.StatusCode)
	}
	if len(body) == 0 {
		t.Fatal("expected a body with content")
	}
}

func TestProxy(t *testing.T) {
	// TODO: after finish proxy package
}
