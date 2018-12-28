package gogtrends

import (
	"bytes"
	"context"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	gAPI             = "https://trends.google.com/trends/api"
	gDaily           = "/dailytrends"
	gRealtime        = "/realtimetrends"
	errParsing       = "failed to parse json"
	errRequestFailed = "failed to perform http request to API"
	errReqDataF      = "request data: code = %d, status = %s, body = %s"
)

type gClient struct {
	c         *http.Client
	defParams url.Values
}

type dailyOut struct {
	Default trendingSearchesDays `json:"default"`
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

type realtimeOut struct {
	StorySummaries storySummary `json:"storySummaries"`
}

type storySummary struct {
	TrendingStories []*TrendingStory `json:"trendingStories"`
}

type TrendingStory struct {
	Title    string             `json:"title"`
	Image    SearchImage        `json:"image"`
	Articles []*TrendingArticle `json:"articles"`
}

type TrendingArticle struct {
	Title   string `json:"articleTitle"`
	URL     string `json:"url"`
	Source  string `json:"source"`
	Time    string `json:"time"`
	Snippet string `json:"snippet"`
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errors.New(errRequestFailed), errReqDataF, resp.StatusCode, resp.Status, resp.Body)
	}

	return resp, nil
}

func (c *gClient) getRespData(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *gClient) trends(ctx context.Context, path, loc string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	p := client.defParams
	p.Set("geo", loc)

	u.RawQuery = p.Encode()

	resp, err := client.do(ctx, u)
	if err != nil {
		return "", err
	}

	data, err := c.getRespData(resp)
	if err != nil {
		return "", err
	}

	return data, nil
}

var client = newGClient()

// TODO: validation for location parameter
// Daily gets daily trends for region specified by loc param
func Daily(ctx context.Context, loc string) ([]*TrendingSearch, error) {
	data, err := client.trends(ctx, gAPI+gDaily, loc)
	if err != nil {
		return nil, err
	}

	out := new(dailyOut)
	// google api returns not valid json :(
	str := strings.Replace(data, ")]}',", "", 1)
	if err := jsoniter.UnmarshalFromString(str, out); err != nil {
		return nil, errors.Wrap(err, errParsing)
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
func Realtime(ctx context.Context, loc string) ([]*TrendingStory, error) {
	data, err := client.trends(ctx, gAPI+gRealtime, loc)
	if err != nil {
		return nil, err
	}

	out := new(realtimeOut)
	// google api returns not valid json :(
	str := strings.Replace(data, ")]}'", "", 1)
	if err := jsoniter.UnmarshalFromString(str, out); err != nil {
		return nil, errors.Wrap(err, errParsing)
	}

	trends := make([]*TrendingStory, 0)
	for _, v := range out.StorySummaries.TrendingStories {
		trends = append(trends, v)
	}

	return trends, nil
}
