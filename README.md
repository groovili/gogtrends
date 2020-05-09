# Google Trends API for Go

###### Unofficial Google Trends API for Golang

[![Coverage Status](https://coveralls.io/repos/github/groovili/gogtrends/badge.svg)](https://coveralls.io/github/groovili/gogtrends) [![Go Report Card](https://goreportcard.com/badge/github.com/groovili/gogtrends)](https://goreportcard.com/report/github.com/groovili/gogtrends) [![License](https://img.shields.io/badge/licence-MIT-blue.svg)](https://github.com/groovili/gogtrends/blob/master/LICENSE)

**gogtrends** is API wrapper which allows to get reports from Google Trends.

All contributions, updates and issues are warmly welcome.

### Installation 

``go get -u github.com/groovili/gogtrends``

#### Debug

To see request-response details use `gogtrends.Debug(true)`

#### Usage

**Daily** and **Realtime** trends used as it is. For both methods user interface language are required. For **Realtime** trends category is required param, list of available categories -  **TrendsCategories**.

Please notice that **Realtime** trends are available only for limited list of locations.


For **InterestOverTime**, **InterestByLocation** and **Related** - widget and user interface language are required.

To get widget you should call **Explore** methods first, it will return constant list of available widgets, every widget corresponds to methods above.

Widget includes request params and unique token for every method.

Also **Explore** method supports single and multiple items for comparision. Please take a look at **ExploreRequest** input.
It supports search by multiple categories and locations which you can get as tree structure by **ExploreCategories** and **ExploreLocations**.

Please notice, when you call **Explore** method for keywords comparison, two first widgets would be for all of compared items, next widgets would be for each of individual items.

### Available methods

* `Daily(ctx context.Context, hl, loc string) ([]*TrendingSearch, error)` - daily trends descending ordered by days and articles corresponding to it.

* `Realtime(ctx context.Context, hl, loc, cat string) ([]*TrendingStory, error)` - represents realtime trends with included articles and sources.

* `Explore(ctx context.Context, r *ExploreRequest, hl string) ([]*ExploreWidget, error)` - widgets with **tokens**. Every widget is related to specific method (`InterestOverTime`, `InterestByLocation`, `Related`) and contains required **token** and request information.

* `InterestOverTime(ctx context.Context, w *ExploreWidget, hl string) ([]*Timeline, error)` - interest over time, dots for chart. 

* `InterestByLocation(ctx context.Context, w *ExploreWidget, hl string) ([]*GeoMap, error)` - interest by location, list for map with geo codes and interest values.

* `Related(ctx context.Context, w *ExploreWidget, hl string) ([]*RankedKeyword, error)` - related topics or queries, supports two types of widgets.

* `TrendsCategories() map[string]string` - available categories for `Realtime` trends.

* `ExploreCategories(ctx context.Context) (*ExploreCatTree, error)` - tree of categories for explore and comparison. Called once, then returned from cache.

* `ExploreLocations(ctx context.Context) (*ExploreLocTree, error)` - tree of locations for explore and comparison. Called once, then returned from cache.

#### Parameters 

* `hl` -  string, user interface language

* `loc` - string, uppercase location (geo) country code, example "US" - United States

* `cat` - string, lowercase category for real time trends, example "all" - all categories

* `exploreReq` - `ExploreRequest` struct, represents search or comparison items.

* `widget` - `ExploreWidget` struct, specific for every method, can be received by `Explore` method.

### Examples

Working detailed examples for all methods and cases can be found in ***example*** folder. Short version below.

```go
// Daily trends
ctx := context.Background()
dailySearches, err := gogtrends.Daily(ctx, "EN", "US")
```

```go
// Get available trends categories and realtime trends
cats := gogtrends.TrendsCategories()
realtime, err := gogtrends.Realtime(ctx, "EN", "US", "all")
```


```go
// Explore available widgets for keywords and get all available stats for it
explore, err := gogtrends.Explore(ctx, 
	    &gogtrends.ExploreRequest{
            ComparisonItems: []*gogtrends.ComparisonItem{
                {
                    Keyword: "Go",
                    Geo:     "US",
                    Time:    "today 12-m",
                },
            },
            Category: 31, // Programming category
            Property: "",
        }, "EN")

// Interest over time
overTime, err := gogtrends.InterestOverTime(ctx, explore[0], "EN")

// Interest by location
byLoc, err := gogtrends.InterestByLocation(ctx, explore[1], "EN")

// Related topics for keyword
relT, err := gogtrends.Related(ctx, explore[2], "EN")

// Related queries for keyword
relQ, err := gogtrends.Related(ctx, explore[3], "EN")

// Compare keywords interest
compare, err := gogtrends.Explore(ctx, 
	    &gogtrends.ExploreRequest{
            ComparisonItems: []*gogtrends.ComparisonItem{
                {
                    Keyword: "Go",
                    Geo:     "US",
                    Time:    "today 12-m",
                },
                {
                    Keyword: "Python",
                    Geo:     "US",
                    Time:    "today 12-m",
                },
                {
                    Keyword: "PHP",
                    Geo:     "US",
                    Time:    "today 12-m",
                },                               
            },
            Category: 31, // Programming category
            Property: "",
        }, "EN")

```

### Licence
 
Package is distributed under [MIT Licence](https://opensource.org/licenses/MIT).
