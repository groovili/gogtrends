package gogtrends

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locUS  = "US"
	catAll = "all"
	langEN = "EN"
)

var exploreReq = &ExploreRequest{
	ComparisonItems: []*ComparisonItem{
		{
			Keyword: "Golang",
			Geo:     locUS,
			Time:    "today+12-m",
		},
	},
	Category: 31, // Programming category
	Property: "",
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

func TestExploreLocations(t *testing.T) {
	exploreLocs, err := ExploreLocations(context.Background())
	assert.NoError(t, err)
	assert.True(t, len(exploreLocs.Children) > 0)
}

func TestExplore(t *testing.T) {
	explore, err := Explore(context.Background(), exploreReq, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)
}

func TestInterestOverTime(t *testing.T) {
	explore, err := Explore(context.Background(), exploreReq, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	overTime, err := InterestOverTime(context.Background(), explore[0], langEN)
	assert.NoError(t, err)
	assert.True(t, len(overTime) > 0)
}

func TestInterestByLocation(t *testing.T) {
	explore, err := Explore(context.Background(), exploreReq, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	overTime, err := InterestByLocation(context.Background(), explore[1], langEN)
	assert.NoError(t, err)
	assert.True(t, len(overTime) > 0)
}

func TestRelated(t *testing.T) {
	explore, err := Explore(context.Background(), exploreReq, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	relatedTopics, err := Related(context.Background(), explore[2], langEN)
	assert.NoError(t, err)
	assert.True(t, len(relatedTopics) > 0)

	relatedQueries, err := Related(context.Background(), explore[3], langEN)
	assert.NoError(t, err)
	assert.True(t, len(relatedQueries) > 0)
}
