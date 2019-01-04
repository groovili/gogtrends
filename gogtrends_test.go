package gogtrends

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locUS  = "US"
	catAll = "all"
	langEN = "EN"
)

func TestRequests(t *testing.T) {
	ctx := context.Background()

	u, err := url.Parse(gAPI + gDaily)
	assert.NoError(t, err)

	p := client.defParams
	p.Set("geo", locUS)
	u.RawQuery = p.Encode()

	data, err := client.do(ctx, u)
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	u, err = url.Parse(gAPI + gRealtime)
	assert.NoError(t, err)
	u.RawQuery = p.Encode()

	data, err = client.do(ctx, u)
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
}

func TestDailyTrending(t *testing.T) {
	locations := TrendsLocations()
	_, ok := locations[locUS]
	assert.True(t, ok)

	resp, err := Daily(context.Background(), langEN, locUS)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title.Query) > 0)
}

func TestRealtimeTrending(t *testing.T) {
	categories := TrendsCategories()
	assert.True(t, len(categories) > 0)
	_, ok := categories[catAll]
	assert.True(t, ok)

	resp, err := Realtime(context.Background(), langEN, locUS, catAll)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title) > 0)
}

func TestExploreCategories(t *testing.T) {
	exploreCats, err := ExploreCategories(context.Background())
	assert.NoError(t, err)
	assert.True(t, len(exploreCats.Children) > 0)
}

func TestExplore(t *testing.T) {
	explore, err := Explore(context.Background(), &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Geo:     locUS,
				Time:    "today+12-m",
			},
		},
		Category: 31, // Programming category
		Property: "",
	}, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)
}

func TestInterestOverTime(t *testing.T) {
	explore, err := Explore(context.Background(), &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Geo:     locUS,
				Time:    "today+12-m",
			},
		},
		Category: 31, // Programming category
		Property: "",
	}, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	overTime, err := InterestOverTime(context.Background(), explore[0], langEN)
	assert.NoError(t, err)
	assert.True(t, len(overTime) > 0)
}
