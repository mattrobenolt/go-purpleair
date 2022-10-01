package purpleair

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	key string
}

var (
	clientOnce    sync.Once
	defaultClient *http.Client
)

const (
	baseURL      = "https://api.purpleair.com/v1"
	apiHeaderKey = "X-Api-Key"
)

func New(key string) *Client {
	return &Client{
		key: key,
	}
}

func ensurePrefix(s string) string {
	if len(s) == 0 {
		return "/"
	}
	if s[0] != '/' {
		return "/" + s
	}
	return s
}

func (c *Client) NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	url = baseURL + ensurePrefix(url)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	clientOnce.Do(func() {
		defaultClient = cleanhttp.DefaultPooledClient()
	})
	req.Header[apiHeaderKey] = []string{c.key}
	return defaultClient.Do(req)
}

func (c *Client) get(ctx context.Context, url string, expected int) ([]byte, *http.Response, error) {
	req, err := c.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return b, resp, err
	}
	if resp.StatusCode != expected {
		errx := &ResponseError{}
		if err := decodeBody(b, errx); err != nil {
			return b, resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return b, resp, errx
	}
	return b, resp, nil
}

func decodeBody[T any](b []byte, v T) error {
	return json.NewDecoder(bytes.NewReader(b)).Decode(v)
}

func Pointer[T any](d T) *T {
	return &d
}

func itoa[N int | uint64 | uint16](i N) string {
	return strconv.Itoa(int(i))
}

type ResponseError struct {
	ApiVersion  string `json:"api_version"`
	TimeStamp   int64  `json:"time_stamp"`
	ErrorX      string `json:"error"`
	Description string `json:"description"`
}

func (r *ResponseError) Error() string {
	return r.ErrorX + ": " + r.Description
}
