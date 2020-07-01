package photonanalytics

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const baseURL = "https://counter.photonengine.com/Counter/api"

type APIClient struct {
	HTTPClient *http.Client
	Logger     Logger
	URL        *url.URL
	token      string
}

func New(token string) (*APIClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &APIClient{
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "", log.LstdFlags),
		URL:        u,
		token:      token,
	}, nil
}

func NewWithClient(token string, httpClient *http.Client) (*APIClient, error) {
	c, err := New(token)
	if err != nil {
		return nil, err
	}
	c.HTTPClient = httpClient
	return c, nil
}

type Params struct {
	headers map[string]string
	queries map[string]string
	body    io.Reader
}

func (c *APIClient) newRequest(
	ctx context.Context,
	method string,
	reqpath string,
	params *Params,
) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, reqpath)
	q := u.Query()
	for k, v := range params.queries {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, method, u.String(), params.body)
	for k, v := range params.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", "go-photonanalytics")
	req.Header.Set("Authorization", "Bearer "+c.token)
	return req, err
}
