package gogtrends

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type gClient struct {
	c          *http.Client
	defParams  url.Values
	categories map[string]string
	locations  map[string]string
}

func newGClient() *gClient {
	p := make(url.Values)
	for k, v := range defaultParams {
		p.Add(k, v)
	}

	return &gClient{
		c:          http.DefaultClient,
		defParams:  p,
		categories: availableCategories,
		locations:  availableLocations,
	}
}

func (c *gClient) do(ctx context.Context, url *url.URL) (*http.Response, error) {
	r := &http.Request{
		URL:    url,
		Method: http.MethodGet,
	}
	r = r.WithContext(ctx)

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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
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
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return data, nil
}

func (c *gClient) validateCategory(cat string) bool {
	_, ok := client.categories[cat]

	return ok
}

func (c *gClient) validateLocation(loc string) bool {
	_, ok := client.locations[loc]

	return ok
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
