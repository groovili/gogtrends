# Google Trends API for Go

###### Unofficial Google Trends API for Golang

**gogtrends** is API wrapper which allows to get reports from Google Trends.

All contributions, updates and issues are warmly welcome.

### Installation 

Go modules support is required to use this package, also all dependencies can be found in `go.sum` file.

```bash 
export GO111MODULE=on
```

Add `github.com/groovili/gogtrends` as import and run `go build` or manually require in go.mod file.

### Available methods

* `dT, err := gogtrends.Daily(ctx, hl, loc)` - `TrendingSearch` structs descending ordered by days and articles corresponding to it.

* `rT, err := gogtrends.Realtime(ctx, hl, loc, cat)` - `TrendingStory` structs, represents realtime trends with included articles and sources.

* `e, err := gogtrends.Explore(ctx, exploreReq, hl)` - list of widgets with **tokens**. Every widget is related to specific method (`InterestOverTime`, `InterestOverLoc`, `RelatedSearches`, `Suggestions`) and contains required **token** and request information.

* `iOT, err :=  gogtrends.InterestOverTime(ctx, widget, hl)` - interest over time, as list of `Timeline` dots for chart. 

* `c := gogtrends.TrendsCategories()` - available categories for Realtime trends.

* `l := gogtrends.TrendsLocations()` - available locations (geo).

* `c, err :=  gogtrends.ExploreCategories(ctx)` - tree of categories for explore and comparison. Called only once, then returned from client cache.

### Parameters 

* `hl` -  string, user interface language

* `loc` - string, uppercase location (geo) country code, example "US" - United States

* `cat` - string, lowercase category for real time trends, example "all" - all categories

* `exploreReq` - `ExploreRequest` struct, represents search or comparison items.

* `widget` - `ExploreWidget` struct, specific for every method, can be received by `Explore` method.

### Licence
 
Package is distributed under [MIT Licence](https://opensource.org/licenses/MIT).