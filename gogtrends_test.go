package gogtrends

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locUS  = "US"
	catAll = "all"
)

func TestRequest(t *testing.T) {
	ctx := context.Background()

	u, err := url.Parse(gAPI + gDaily)
	assert.NoError(t, err)

	p := client.defParams
	p.Set("geo", locUS)
	u.RawQuery = p.Encode()

	resp, err := client.do(ctx, u)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK)

	u, err = url.Parse(gAPI + gRealtime)
	assert.NoError(t, err)
	u.RawQuery = p.Encode()

	resp, err = client.do(ctx, u)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK)
}

func TestDailyTrending(t *testing.T) {
	locations := AvailableLocations()
	_, ok := locations[locUS]
	assert.True(t, ok)

	resp, err := Daily(context.Background(), locUS)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title.Query) > 0)
}

func TestRealtimeTrending(t *testing.T) {
	categories := RealtimeAvailableCategories()
	assert.True(t, len(categories) > 0)
	_, ok := categories[catAll]
	assert.True(t, ok)

	resp, err := Realtime(context.Background(), locUS, catAll)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title) > 0)
}
