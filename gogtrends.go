package gogtrends

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

var client = newGClient()

func Debug(debug bool) {
	client.debug = debug
}

// TrendsLocations returns general list of locations as map[param]name
func TrendsLocations() map[string]string {
	return client.locations
}

// TrendsCategories return list of available categories for Realtime method as [param]description map
func TrendsCategories() map[string]string {
	return client.categories
}

// Daily gets daily trends descending ordered by days and articles corresponding to it.
func Daily(ctx context.Context, hl, loc string) ([]*TrendingSearch, error) {
	if !client.validateLocation(loc) {
		return nil, errors.New(errInvalidLocation)
	}

	data, err := client.trends(ctx, gAPI+gDaily, hl, loc)
	if err != nil {
		return nil, err
	}

	// google api returns not valid json :(
	str := strings.Replace(data, ")]}',", "", 1)

	out := new(dailyOut)
	if err := client.unmarshal(str, out); err != nil {
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

// Realtime represents realtime trends with included articles and sources.
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

	// google api returns not valid json :(
	str := strings.Replace(data, ")]}'", "", 1)

	out := new(realtimeOut)
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	return out.StorySummaries.TrendingStories, nil
}

// ExploreCategories gets available categories for explore and comparison and caches it in client
func ExploreCategories(ctx context.Context) (*ExploreCatTree, error) {
	if client.exploreCats != nil {
		return client.exploreCats, nil
	}

	u, err := url.Parse(gAPI + gSCategories)
	if err != nil {
		return nil, err
	}

	b, err := client.do(ctx, u)
	str := strings.Replace(string(b), ")]}'", "", 1)

	out := new(ExploreCatTree)
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	client.exploreCats = out

	return out, nil
}

// ExploreLocations gets available locations for explore and comparison and caches it in client
func ExploreLocations(ctx context.Context) (*ExploreLocTree, error) {
	if client.exploreLocs != nil {
		return client.exploreLocs, nil
	}

	u, err := url.Parse(gAPI + gSGeo)
	if err != nil {
		return nil, err
	}

	b, err := client.do(ctx, u)
	str := strings.Replace(string(b), ")]}'", "", 1)

	out := new(ExploreLocTree)
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	client.exploreLocs = out

	return out, nil
}

// Explore list of widgets with tokens. Every widget
// is related to specific method (`InterestOverTime`, `InterestOverLoc`, `RelatedSearches`, `Suggestions`)
// and contains required token and request information.
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
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	return out.Widgets, nil
}

// InterestOverTime as list of `Timeline` dots for chart.
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
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	return out.Default.TimelineData, nil
}

// InterestByLocation as list of `GeoMap`, with geo codes and interest values.
func InterestByLocation(ctx context.Context, w *ExploreWidget, hl string) ([]*GeoMap, error) {
	if w.ID != intOverRegionID {
		return nil, errors.New(errInvalidWidgetType)
	}

	u, err := url.Parse(gAPI + gSIntOverReg)
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

	out := new(GeoOut)
	if err := client.unmarshal(str, out); err != nil {
		return nil, err
	}

	return out.Default.GeoMapData, nil
}
