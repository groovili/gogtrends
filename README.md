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

* `dT, err := gogtrends.Daily(ctx, "US")` - returns `TrendingSearch` structs descending ordered by days and articles corresponding to it.

* `rT, err := gogtrends.Realtime(ctx, "US", "all")` - returns `TrendingStory` structs, which represent realtime trends.

* `c := TrendsCategories()` - list of available categories for Realtime trends.

* `l := TrendsLocations()` - list of available locations (geo).

* `c, err :=  gogtrends.ExploreCategories(ctx)` - tree of categories for explore and comparison. Called only once, then returned from client cache.

### Parameters 

* `loc` - string, uppercase location (geo) country code, example "US" - United States

* `cat` - string, lowercase category for real time trends, example "all" - all categories

### Licence
 
Package is distributed under [MIT Licence](https://opensource.org/licenses/MIT).