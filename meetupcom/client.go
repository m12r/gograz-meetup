package meetupcom

import (
	"context"
	"net/http"
	"net/url"
)

type Client struct {
	opts ClientOptions
}

type ClientOptions struct {
	APIKey string
}

func NewClient(opts ClientOptions) *Client {
	c := Client{opts: opts}
	return &c
}

func (c *Client) executeGet(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	query.Set("key", c.opts.APIKey)
	u := url.URL{Scheme: "https", Host: "api.meetup.com", Path: path, RawQuery: query.Encode()}
	req := http.Request{
		Method: http.MethodGet,
		URL:    &u,
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}
