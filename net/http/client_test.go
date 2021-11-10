package http_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aceaura/libra/net/http"
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
	form := url.Values{}
	form.Add("platform", "pc")
	form.Add("sa", "pcindex_entry")
	resp, body, err := http.Get("https://top.baidu.com/board",
		http.ClientOption.Form(form),
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
	resp, body, err := http.NewClient().Form(url.Values{
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
	resp, body, err := http.NewClient().Protocol("https").Get(
		"top.baidu.com/board?platform=pc&sa=pcindex_entry")
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
	_, body, err := http.NewClient().Timeout(
		time.Duration(1) * time.Microsecond).Get(
		"https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), "context deadline exceeded") < 0 {
		t.Fatal("expected an error with timeout")
	}
	t.Log(string(body))
}

func TestClientContentType(t *testing.T) {
	resp, body, err := http.NewClient(
		http.ClientOption.ContentType("exception"),
		http.ClientOption.Retry(3),
	).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
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
		t.Fatalf("unexpected content type getting from client: %s",
			contentType)
	}
	t.Log(string(body))
}

func TestClientResponseBodyReader(t *testing.T) {
	var body []byte
	var err error
	respBodyFunc := func(r io.Reader) error {
		body, err = ioutil.ReadAll(r)
		return err
	}
	resp, _, err := http.NewClient(
		http.ClientOption.ResponseBodyFunc(respBodyFunc),
		http.ClientOption.Retry(3),
	).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
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
	respBodyFunc := func(r io.Reader) error {
		panic(errors.New(text))
	}
	_, _, err := http.NewClient(
		http.ClientOption.ResponseBodyFunc(respBodyFunc),
		http.ClientOption.Safety(),
	).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), text) < 0 {
		t.Fatal("expected an error with safety exception")
	}
}

func TestClientRequestBody(t *testing.T) {
	route := http.Route{"/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(body))
	}}
	text := "this is a body text"
	reqBodyFunc := func() (io.Reader, error) {
		return strings.NewReader(text), nil
	}
	http.Serve("localhost:1989",
		http.ServerOption.Background(),
		http.ServerOption.Routes(route))

	resp, body, err := http.NewClient(
		http.ClientOption.RequestBody(reqBodyFunc),
	).Get("localhost:1989")
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
	t.Log(string(body))
}

func TestClientRetry(t *testing.T) {
	const (
		clientTimout = 1
		serverSleep  = 100
		maxRetry     = 3
	)
	var count int = 1
	route := http.Route{"/", func(w http.ResponseWriter, r *http.Request) {
		if count < maxRetry {
			count++
			time.Sleep(time.Second * serverSleep)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(body))
	}}

	http.Serve("localhost:1989",
		http.ServerOption.Background(),
		http.ServerOption.Routes(route),
	)
	text := "this is a body text"
	reqBodyFunc := func() (io.Reader, error) {
		return strings.NewReader(text), nil
	}
	resp, body, err := http.NewClient(
		http.ClientOption.Timeout(time.Second*clientTimout),
		http.ClientOption.Retry(maxRetry),
		http.ClientOption.RequestBody(reqBodyFunc),
	).Get("localhost:1989")
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

func TestProxy(t *testing.T) {
	http.NewServer().Proxy().Background().Serve("localhost:1990")
	resp, body, err := http.Get("www.baidu.com", http.ClientOption.Proxy("http://localhost:1990"))
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

func TestContext(t *testing.T) {
	const (
		clientTimout = 1
		serverSleep  = 100
	)
	route := http.Route{"/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * serverSleep)
	}}
	http.Serve("localhost:1989",
		http.ServerOption.Background(),
		http.ServerOption.Routes(route),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*clientTimout)
	defer cancel()
	text := "this is a body text"
	reqBodyFunc := func() (io.Reader, error) {
		return strings.NewReader(text), nil
	}
	_, _, err := http.NewClient(
		http.ClientOption.RequestBody(reqBodyFunc),
		http.ClientOption.Context(ctx),
	).Get("localhost:1989")
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return
		}
	}
	t.Fatalf("expected an deadline exceeded error")
}
