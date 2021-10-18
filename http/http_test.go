package http_test

import (
	"fmt"
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
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http get received empty body")
	}
	t.Log(string(body))
}

func TestClientDoGet(t *testing.T) {
	resp, body, err := http.Do(http.GET, "www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http get received empty body")
	}
	t.Log(string(body))
}

func TestClientPost(t *testing.T) {
	resp, body, err := http.Post("www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http post received empty body")
	}
	t.Log(string(body))
}

func TestClientDoPost(t *testing.T) {
	resp, body, err := http.Do(http.POST, "www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http post received empty body")
	}
	t.Log(string(body))
}

func TestClientHead(t *testing.T) {
	resp, body, err := http.Head("www.baidu.com")
	if err != nil {
		t.Error(resp)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) > 0 {
		t.Error("http head received body content")
	}
	if len(resp.Header) == 0 {
		t.Error("http head received empty header")
	}
	t.Log(resp.Header)
}

func TestClientDoHead(t *testing.T) {
	resp, body, err := http.Do(http.HEAD, "www.baidu.com")
	if err != nil {
		t.Error(resp)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) > 0 {
		t.Error("http head received body content")
	}
	if len(resp.Header) == 0 {
		t.Error("http head received empty header")
	}
	t.Log(resp.Header)
}

func TestClientParam(t *testing.T) {
	resp, body, err := http.Get("https://top.baidu.com/board",
		http.WithParam("platform", "pc"),
		http.WithParam("sa", "pcindex_entry"),
	)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http received empty body")
	}
	t.Log(string(body))
}

func TestClientForm(t *testing.T) {
	resp, body, err := http.NewClient().WithForm(url.Values{
		"platform": []string{"pc"},
		"sa":       []string{"pcindex_entry"},
	}).Get("https://top.baidu.com/board")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http received empty body")
	}
	t.Log(string(body))
}

func TestClientProtocol(t *testing.T) {
	resp, body, err := http.NewClient().WithProtocol("https").Get("top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http get received empty body")
	}
	t.Log(string(body))
}

func TestClientTimeout(t *testing.T) {
	_, _, err := http.NewClient().WithTimeout(time.Duration(1) * time.Microsecond).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), "context deadline exceeded") < 0 {
		t.Error("http timeout does not work")
	}
}

func TestContentType(t *testing.T) {
	resp, body, err := http.NewClient().WithContentType("exception").WithRetry(3).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http received empty body")
	}
	if resp.Request.Header["Content-Type"][0] != "exception" {
		t.Error("http content type does not work")
	}
}

func TestBody(t *testing.T) {
	http.Serve("localhost:1989",
		http.WithBackground(true),
		http.WithRoute("/call", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprint(w, body)
		}))
	text := "this is a body text"
	resp, body, err := http.NewClient().WithBody(text).Get("localhost:1989/call")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("http received response status not 200")
	}
	if len(body) == 0 {
		t.Error("http received empty body")
	}
	if string(body) != text {
		t.Error("http body does not work")
	}
}

func TestRetry(t *testing.T) {
	_, _, err := http.NewClient().WithTimeout(time.Duration(1) * time.Microsecond).WithRetry(3).Get("https://top.baidu.com/board?platform=pc&sa=pcindex_entry")
	if strings.Index(err.Error(), "context deadline exceeded") < 0 {
		t.Error("http timeout does not work")
	}
}
