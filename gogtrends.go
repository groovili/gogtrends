package gogtrends

import (
	"bytes"
	"context"
	"github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type dailyOut struct {
	Default *trendingSearchesDays `json:"default"`
}

type trendingSearchesDays struct {
	Searches []*trendingSearchDays `json:"trendingSearchesDays"`
}

type trendingSearchDays struct {
	FormattedDate string            `json:"formattedDate"`
	Searches      []*TrendingSearch `json:"trendingSearches"`
}

type TrendingSearch struct {
	Title            SearchTitle      `json:"title"`
	FormattedTraffic string           `json:"formattedTraffic"`
	Image            SearchImage      `json:"image"`
	Articles         []*SearchArticle `json:"articles"`
}

type SearchTitle struct {
	Query string `json:"query"`
}

type SearchImage struct {
	NewsURL  string `json:"newsUrl"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
}

type SearchArticle struct {
	Title   string      `json:"title"`
	TimeAgo string      `json:"timeAgo"`
	Source  string      `json:"source"`
	Image   SearchImage `json:"image"`
	URL     string      `json:"url"`
	Snippet string      `json:"snippet"`
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

var client = newGClient()

// Daily gets daily trends for region specified by loc param
func Daily(ctx context.Context, loc string) ([]*TrendingSearch, error) {
	resp, err := client.trends(ctx, gAPI+gDaily, loc)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	out := new(dailyOut)
	// google api returns not valid json :(
	str := strings.Replace(buf.String(), ")]}',", "", 1)
	err = jsoniter.UnmarshalFromString(str, out)
	if err != nil {
		return nil, err
	}

	searches := make([]*TrendingSearch, 0)
	for _, v := range out.Default.Searches {
		for _, k := range v.Searches {
			searches = append(searches, k)
		}
	}

	return searches, nil
}

// Realtime gets recent trends for region specified by loc param, available for limited list of regions
func Realtime(ctx context.Context, loc string) (*http.Response, error) {
	return client.trends(ctx, gAPI+gRealtime, loc)
}
