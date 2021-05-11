package dhttp

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-starter-project/pkg/threadlocal"
)

type DiancaiClient struct {
	*http.Client
}

func (c *DiancaiClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Diancai-Correlation-Id", threadlocal.GetCorrelationID())
	return c.Client.Do(req)
}

func (c *DiancaiClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *DiancaiClient) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *DiancaiClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *DiancaiClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func WrapClient(c *http.Client) *DiancaiClient {
	return &DiancaiClient{
		Client: c,
	}
}
