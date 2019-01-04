package gogtrends

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type gClient struct {
	c           *http.Client
	defParams   url.Values
	categories  map[string]string
	locations   map[string]string
	exploreCats *ExploreCategoriesTree
	cookie      string
	debug       bool
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
		cookie:     "",
		debug:      false,
	}
}

func (c *gClient) do(ctx context.Context, u *url.URL) ([]byte, error) {
	p, _ := url.PathUnescape(u.String())
	u, _ = u.Parse(p)

	r := &http.Request{
		URL:    u,
		Method: http.MethodGet,
	}
	r = r.WithContext(ctx)

	r.Header = make(http.Header)
	if client.cookie != "" {
		r.Header.Add("cookie", client.cookie)
	}

	if client.debug {
		log.Info("[Debug] Request with params: ", r.URL)
	}

	resp, err := c.c.Do(r)
	if err != nil {
		return nil, err
	}

	if client.debug {
		log.Info("[Debug] Response: ", resp)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		cookie := strings.Split(resp.Header.Get("set-cookie"), ";")
		if len(cookie) > 0 {
			client.cookie = cookie[0]
			r.Header.Add("cookie", cookie[0])

			resp, err = c.c.Do(r)
			if err != nil {
				return nil, err
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errors.New(errRequestFailed), errReqDataF, resp.StatusCode, resp.Status, resp.Body)
	}

	data, err := c.getRespData(resp)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *gClient) getRespData(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *gClient) trends(ctx context.Context, path, hl, loc string, args ...map[string]string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	// required params
	p := client.defParams
	p.Set(paramGeo, loc)
	p.Set(paramHl, hl)

	// additional params
	if len(args) > 0 {
		for _, arg := range args {
			for n, v := range arg {
				p.Set(n, v)
			}
		}
	}

	u.RawQuery = p.Encode()

	data, err := client.do(ctx, u)
	if err != nil {
		return "", err
	}

	return string(data), nil
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

func Debug(debug bool) {
	client.debug = debug
}

func FormatTime(t time.Time) string {
	return t.Format(timeLayoutFull)
}

// TrendsLocations returns general list of locations as map[param]name
func TrendsLocations() map[string]string {
	return client.locations
}

// Daily gets daily trends for region for language and location param
func Daily(ctx context.Context, hl, loc string) ([]*TrendingSearch, error) {
	if !client.validateLocation(loc) {
		return nil, errors.New(errInvalidLocation)
	}

	data, err := client.trends(ctx, gAPI+gDaily, hl, loc)
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

// TrendsCategories return list of available categories for Realtime method as [param]description map
func TrendsCategories() map[string]string {
	return client.categories
}

// Realtime gets current trends for language, location and category, available for limited list of locations
func Realtime(ctx context.Context, hl, loc, cat string) ([]*TrendingStory, error) {
	if !client.validateLocation(loc) {
		return nil, errors.New(errInvalidLocation)
	}

	if !client.validateCategory(cat) {
		return nil, errors.New(errInvalidCategory)
	}

	data, err := client.trends(ctx, gAPI+gRealtime, hl, loc, map[string]string{paramCat: cat})
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

// ExploreCategories gets available categories for explore and comparison and caches it in client
func ExploreCategories(ctx context.Context) (*ExploreCategoriesTree, error) {
	if client.exploreCats != nil {
		return client.exploreCats, nil
	}

	u, err := url.Parse(gAPI + gSCategories)
	if err != nil {
		return nil, err
	}

	b, err := client.do(ctx, u)
	str := strings.Replace(string(b), ")]}'", "", 1)

	out := new(ExploreCategoriesTree)
	if err := jsoniter.UnmarshalFromString(str, out); err != nil {
		return nil, errors.Wrap(err, errParsing)
	}

	client.exploreCats = out

	return out, nil
}

func Explore(ctx context.Context, r *ExploreRequest, hl string) ([]*ExploreWidget, error) {
	u, err := url.Parse(gAPI + gSExplore)
	if err != nil {
		return nil, err
	}

	p := make(url.Values)
	p.Set(paramTZ, "0")
	p.Set(paramHl, hl)

	mReq, err := jsoniter.MarshalToString(r)
	if err != nil {
		return nil, errors.Wrapf(err, errInvalidRequest)
	}

	p.Set(paramReq, mReq)
	u.RawQuery = p.Encode()

	b, err := client.do(ctx, u)
	if err != nil {
		return nil, err
	}

	str := strings.Replace(string(b), ")]}'", "", 1)

	out := new(ExploreOut)
	if err := jsoniter.UnmarshalFromString(str, out); err != nil {
		return nil, errors.Wrap(err, errParsing)
	}

	return out.Widgets, nil
}

func InterestOverTime(ctx context.Context, w *ExploreWidget, hl string) ([]*Timeline, error) {
	if w.ID != intOverTimeWidgetID {
		return nil, errors.New(errInvalidWidgetType)
	}

	u, err := url.Parse(gAPI + gSIntOverTime)
	if err != nil {
		return nil, err
	}

	p := make(url.Values)
	p.Set(paramTZ, "0")
	p.Set(paramHl, hl)
	p.Set(paramToken, w.Token)

	mReq, err := jsoniter.MarshalToString(w.Request)
	if err != nil {
		return nil, errors.Wrapf(err, errInvalidRequest)
	}

	p.Set(paramReq, mReq)
	u.RawQuery = p.Encode()

	b, err := client.do(ctx, u)
	if err != nil {
		return nil, err
	}

	str := strings.Replace(string(b), ")]}',", "", 1)
	out := new(MultilineOut)
	if err := jsoniter.UnmarshalFromString(str, out); err != nil {
		return nil, errors.Wrap(err, errParsing)
	}

	return out.Default.TimelineData, nil
}
