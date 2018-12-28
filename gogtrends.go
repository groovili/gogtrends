package gogtrends

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	gAPI      = "https://trends.google.com/trends/api"
	gDaily    = "/dailytrends"
	gRealtime = "/realtimetrends"

	paramCat = "cat"
	paramGeo = "geo"

	errParsing         = "failed to parse json"
	errRequestFailed   = "failed to perform http request to API"
	errReqDataF        = "request data: code = %d, status = %s, body = %s"
	errInvalidCategory = "invalid category param"
	errInvalidLocation = "invalid location param"
)

type gClient struct {
	c          *http.Client
	defParams  url.Values
	categories map[string]string
	locations  map[string]string
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
	Title            *SearchTitle     `json:"title"`
	FormattedTraffic string           `json:"formattedTraffic"`
	Image            *SearchImage     `json:"image"`
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
	Title   string       `json:"title"`
	TimeAgo string       `json:"timeAgo"`
	Source  string       `json:"source"`
	Image   *SearchImage `json:"image"`
	URL     string       `json:"url"`
	Snippet string       `json:"snippet"`
}

type realtimeOut struct {
	StorySummaries *storySummary `json:"storySummaries"`
}

type storySummary struct {
	TrendingStories []*TrendingStory `json:"trendingStories"`
}

type TrendingStory struct {
	Title    string             `json:"title"`
	Image    *SearchImage       `json:"image"`
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

func getAvailableLocations() map[string]string {
	return map[string]string{
		"AU": "Australia",
		"AT": "Austria",
		"AR": "Argentina",
		"BE": "Belgium",
		"BR": "Brazil",
		"GB": "United Kingdom",
		"HU": "Hungary",
		"VN": "Vietnam",
		"DE": "Germany",
		"HK": "Hong Kong",
		"GR": "Greece",
		"DK": "Denmark",
		"EG": "Egypt",
		"IL": "Israel",
		"IN": "India",
		"IE": "Ireland",
		"IT": "Italy",
		"CA": "Canada",
		"KE": "Kenia",
		"CO": "Columbia",
		"MY": "Malaysia",
		"MX": "Mexico",
		"NG": "Nigeria",
		"NL": "Netherlands",
		"NZ": "New Zeland",
		"NO": "Norway",
		"PL": "Poland",
		"PT": "Portugal",
		"KR": "Korean Republic",
		"RU": "Russia",
		"RO": "Romania",
		"SA": "Saudi Arabia",
		"SG": "Singapore",
		"US": "United States",
		"TH": "Thailand",
		"TW": "Taiwan",
		"TR": "Turkey",
		"UA": "Ukraine",
		"PH": "Philippines",
		"FI": "Finland",
		"FR": "France",
		"CZ": "Czech Republic",
		"CL": "Chili",
		"CH": "Switzerland",
		"SE": "Sweden",
		"ZA": "Republic of South Africa",
		"JP": "Japan",
	}
}

func getAvailableCategories() map[string]string {
	return map[string]string{
		"all": "all",
		"b":   "business",
		"h":   "main news",
		"m":   "health",
		"t":   "science and technics",
		"e":   "entertainment",
		"s":   "sport",
	}
}

func newGClient() *gClient {
	p := url.Values{}
	mParam := getDefParams()
	for k, v := range mParam {
		p.Add(k, v)
	}

	return &gClient{
		c:          new(http.Client),
		defParams:  p,
		categories: getAvailableCategories(),
		locations:  getAvailableLocations(),
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

func (c *gClient) trends(ctx context.Context, path, loc string, args ...map[string]string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	// required param for all methods
	p := client.defParams
	p.Set(paramGeo, loc)

	// additional params
	if len(args) > 0 {
		for _, arg := range args {
			for n, v := range arg {
				p.Set(n, v)
			}
		}
	}

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

func (c *gClient) validateCategory(cat string) bool {
	for c := range client.categories {
		if c == cat {
			return true
		}
	}

	return false
}

func (c *gClient) validateLocation(loc string) bool {
	for l := range client.locations {
		if loc == l {
			return true
		}
	}

	return false
}

var client = newGClient()

// AvailableLocations returns general list of locations as map[param]name
func AvailableLocations() map[string]string {
	return client.locations
}

// Daily gets daily trends for region by location param
func Daily(ctx context.Context, loc string) ([]*TrendingSearch, error) {
	if !client.validateLocation(loc) {
		return nil, errors.New(errInvalidLocation)
	}

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

// RealtimeAvailableCategories return list of available categories for Realtime method as [param]description map
func RealtimeAvailableCategories() map[string]string {
	return client.categories
}

// Realtime gets current trends for location and category, available for limited list of locations
func Realtime(ctx context.Context, loc, cat string) ([]*TrendingStory, error) {
	if !client.validateLocation(loc) {
		return nil, errors.New(errInvalidLocation)
	}

	if !client.validateCategory(cat) {
		return nil, errors.New(errInvalidCategory)
	}

	data, err := client.trends(ctx, gAPI+gRealtime, loc, map[string]string{paramCat: cat})
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
