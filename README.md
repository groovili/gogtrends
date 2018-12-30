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

* `dailyT, err := gogtrends.Daily(ctx, "US")` - returns `TrendingSearch` structs descending ordered by days and articles corresponding to it.

* `realtimeT, err := gogtrends.Realtime(ctx, "US", "all")` - returns `TrendingStory` structs, which represent realtime trends.

* `categories := RealtimeAvailableCategories()` - list of available categories.

* `locations := AvailableLocations()` - list of available locations (geo).

### Parameters 

* `loc` - string, uppercase location (geo) country code, example "US" - United States

* `cat` - string, lowercase category for real time trends, example "all" - all categories

### Licence
 
Package is distributed under [MIT Licence](https://opensource.org/licenses/MIT).