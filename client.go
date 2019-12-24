package gogtrends

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type gClient struct {
	c           *http.Client
	defParams   url.Values
	trendsCats  map[string]string
	exploreCats *ExploreCatTree
	exploreLocs *ExploreLocTree
	cookie      string
	debug       bool
}

func newGClient() *gClient {
	// default request params
	p := make(url.Values)
	for k, v := range defaultParams {
		p.Add(k, v)
	}

	return &gClient{
		c:          http.DefaultClient,
		defParams:  p,
		trendsCats: trendsCategories,
		cookie:     "",
		debug:      false,
	}
}

func (c *gClient) do(ctx context.Context, u *url.URL) ([]byte, error) {
	p := strings.Replace(u.String(), "%2B", "+", -1)
	u, _ = u.Parse(p)

	r := &http.Request{
		URL:    u,
		Method: http.MethodGet,
	}
	r = r.WithContext(ctx)

	r.Header = make(http.Header)
	if client.cookie != "" {
		r.Header.Add("Cookie", client.cookie)
	}

	r.Header.Add("Accept", "application/json")

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
		cookie := strings.Split(resp.Header.Get("Set-Cookie"), ";")
		if len(cookie) > 0 {
			client.cookie = cookie[0]
			r.Header.Set("Cookie", cookie[0])

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

func (c *gClient) unmarshal(str string, dest interface{}) error {
	if err := jsoniter.UnmarshalFromString(str, dest); err != nil {
		return errors.Wrap(err, errParsing)
	}

	return nil
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
	_, ok := client.trendsCats[cat]

	return ok
}
