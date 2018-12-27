package gogtrends

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const locUS = "US"

func TestRequest(t *testing.T) {
	ctx := context.Background()
	resp, err := client.trends(ctx, gAPI+gDaily, locUS)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK)

	resp, err = Realtime(ctx, locUS)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK)
}

func TestDailyTrending(t *testing.T) {
	resp, err := Daily(context.Background(), locUS)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title.Query) > 0)
}

func TestRealtimeTrending(t *testing.T) {
	resp, err := Realtime(context.Background(), locUS)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK)
}
