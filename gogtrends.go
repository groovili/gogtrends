package gogtrends

import (
	"context"
	"net/http"
	"net/url"
)

const (
	gAPI      = "https://trends.google.com/trends/api"
	gDaily    = "/dailytrends"
	gRealtime = "/realtimetrends"
)

type gClient struct {
	c         *http.Client
	defParams url.Values
}

func newGClient() *gClient {
	p := url.Values{}
	mParam := getDefParams()
	for k, v := range mParam {
		p.Add(k, v)
	}

	return &gClient{
		c:         new(http.Client),
		defParams: p,
	}
}

var client = newGClient()

func (c *gClient) do(ctx context.Context, url *url.URL) (*http.Response, error) {
	r := &http.Request{
		URL:    url,
		Method: http.MethodGet,
	}
	r.WithContext(ctx)

	resp, err := c.c.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getDefParams() map[string]string {
	return map[string]string{
		"tz":  "0",
		"cat": "all",
		"fi":  "0",
		"fs":  "0",
		"geo": "US",
		"hl":  "EN",
		"ri":  "10",
		"rs":  "10",
	}
}

func (c *gClient) trends(ctx context.Context, path, loc string) (*http.Response, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	p := client.defParams
	p.Set("geo", loc)

	u.RawQuery = p.Encode()

	resp, err := client.do(ctx, u)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Daily gets daily trends for region specified by loc param
func Daily(ctx context.Context, loc string) (*http.Response, error) {
	return client.trends(ctx, gAPI+gDaily, loc)
}

// Realtime gets recent trends for region specified by loc param, available for limited list of regions
func Realtime(ctx context.Context, loc string) (*http.Response, error) {
	return client.trends(ctx, gAPI+gRealtime, loc)
}
