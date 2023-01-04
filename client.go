package gogtrends

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"log"

	"github.com/pkg/errors"
)

const (
	headerKeyAccept    = "Accept"
	headerKeyCookie    = "Cookie"
	headerKeySetCookie = "Set-Cookie"
	contentTypeJSON    = "application/json"
)

type gClient struct {
	c         *http.Client
	defParams url.Values

	tcm        *sync.RWMutex
	trendsCats map[string]string

	cm          *sync.RWMutex
	exploreCats *ExploreCatTree

	lm          *sync.RWMutex
	exploreLocs *ExploreLocTree

	cookie string
	debug  bool
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
		tcm:        new(sync.RWMutex),
		trendsCats: trendsCategories,
		cm:         new(sync.RWMutex),
		lm:         new(sync.RWMutex),
	}
}

func (c *gClient) defaultParams() url.Values {
	out := make(map[string][]string, len(c.defParams))
	for i, v := range c.defParams {
		out[i] = make([]string, len(v))
		copy(out[i], v)
	}

	return out
}

func (c *gClient) getCategories() *ExploreCatTree {
	c.cm.RLock()
	defer c.cm.RUnlock()
	return c.exploreCats
}

func (c *gClient) setCategories(cats *ExploreCatTree) {
	c.cm.Lock()
	defer c.cm.Unlock()
	c.exploreCats = cats
}

func (c *gClient) getLocations() *ExploreLocTree {
	c.lm.RLock()
	defer c.lm.RUnlock()
	return c.exploreLocs
}

func (c *gClient) setLocations(locs *ExploreLocTree) {
	c.lm.Lock()
	defer c.lm.Unlock()
	c.exploreLocs = locs
}

func (c *gClient) do(ctx context.Context, u *url.URL) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, errCreateRequest)
	}

	r.Header.Add(headerKeyAccept, contentTypeJSON)

	geo := u.Query().Get(paramGeo)
	if len(geo) == 0 {
		geo = "US"
	}
	if len(client.cookie) == 0 {
		if err = c.getCookie(ctx, geo); err != nil {
			log.Println("[Debug] Get Cookie failed: ", err.Error())
		}
	}
	r.Header.Add(headerKeyCookie, client.cookie)

	if client.debug {
		log.Println("[Debug] Request with params: ", r.URL)
	}

	resp, err := c.c.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, errDoRequest)
	}
	defer resp.Body.Close()

	if client.debug {
		log.Println("[Debug] Response: ", resp)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		if err = c.getCookie(ctx, geo); err != nil {
			log.Println("[Debug] Get Cookie failed: ", err.Error())
			return nil, errors.Wrap(err, errReqDataF)
		}
		r.Header.Set(headerKeyCookie, client.cookie)
		resp, err = c.c.Do(r)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(ErrRequestFailed, errReqDataF, resp.StatusCode, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *gClient) getCookie(ctx context.Context, geo string) (err error) {
	u, _ := url.Parse(gCookieAPI)

	p := make(url.Values)
	p.Set(paramGeo, geo)

	u.RawQuery = p.Encode()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return errors.Wrap(err, errCreateRequest)
	}

	resp, err := c.c.Do(r)
	if err != nil {
		return errors.Wrap(err, errDoRequest)
	}
	defer resp.Body.Close()

	for _, v := range strings.Split(resp.Header.Get(headerKeySetCookie), ";") {
		if strings.HasPrefix(v, "NID=") {
			client.cookie = v
		}
	}
	return
}

func (c *gClient) unmarshal(str string, dest interface{}) error {
	if err := jsoniter.UnmarshalFromString(str, dest); err != nil {
		return errors.Wrap(err, errParsing)
	}

	return nil
}

func (c *gClient) trends(ctx context.Context, path, hl, loc string, args ...map[string]string) (string, error) {
	u, _ := url.Parse(path)

	// required params
	p := client.defaultParams()
	if len(loc) > 0 {
		p.Set(paramGeo, loc)
	}
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
	c.tcm.RLock()
	_, ok := client.trendsCats[cat]
	c.tcm.RUnlock()

	return ok
}
