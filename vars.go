package gogtrends

import (
	"sort"
	"strconv"
	"strings"
)

const (
	gAPI = "https://trends.google.com/trends/api"

	gDaily    = "/dailytrends"
	gRealtime = "/realtimetrends"

	gSExplore      = "/explore"
	gSCategories   = "/explore/pickers/category"
	gSGeo          = "/explore/pickers/geo"
	gSRelated      = "/widgetdata/relatedsearches"
	gSIntOverTime  = "/widgetdata/multiline"
	gSIntOverReg   = "/widgetdata/comparedgeo"
	gSAutocomplete = "/autocomplete"

	paramHl    = "hl"
	paramCat   = "cat"
	paramGeo   = "geo"
	paramReq   = "req"
	paramTZ    = "tz"
	paramToken = "token"

	compareDataMode = "PERCENTAGES"
)

type WidgetType string

const (
	IntOverTimeWidgetID WidgetType = "TIMESERIES"
	IntOverRegionID     WidgetType = "GEO_MAP"
	RelatedQueriesID    WidgetType = "RELATED_QUERIES"
	RelatedTopicsID     WidgetType = "RELATED_TOPICS"
)

var (
	defaultParams = map[string]string{
		paramTZ:  "0",
		paramCat: "all",
		"fi":     "0",
		"fs":     "0",
		paramHl:  "EN",
		"ri":     "300",
		"rs":     "20",
	}
	trendsCategories = map[string]string{
		"all": "all",
		"b":   "business",
		"h":   "main news",
		"m":   "health",
		"t":   "science and technics",
		"e":   "entertainment",
		"s":   "sport",
	}
)

type dailyOut struct {
	Default *trendingSearchesDays `json:"default" bson:"default"`
}

type trendingSearchesDays struct {
	Searches []*trendingSearchDays `json:"trendingSearchesDays" bson:"trending_search_days"`
}

type trendingSearchDays struct {
	FormattedDate string            `json:"formattedDate" bson:"formatted_date"`
	Searches      []*TrendingSearch `json:"trendingSearches" bson:"searches"`
}

// TrendingSearch is a representation trending search in period of 24 hours
type TrendingSearch struct {
	Title            *SearchTitle     `json:"title" bson:"title"`
	FormattedTraffic string           `json:"formattedTraffic" bson:"formatted_traffic"`
	Image            *SearchImage     `json:"image" bson:"image"`
	Articles         []*SearchArticle `json:"articles" bson:"articles"`
}

// SearchTitle is a user query string for daily trending search
type SearchTitle struct {
	Query string `json:"query" bson:"query"`
}

// SearchImage is a picture of trending search
type SearchImage struct {
	NewsURL  string `json:"newsUrl" bson:"news_url"`
	Source   string `json:"source" bson:"source"`
	ImageURL string `json:"imageUrl" bson:"image_url"`
}

// SearchArticle is a news relative to trending search
type SearchArticle struct {
	Title   string       `json:"title" bson:"title"`
	TimeAgo string       `json:"timeAgo" bson:"time_ago"`
	Source  string       `json:"source" bson:"source"`
	Image   *SearchImage `json:"image" bson:"image"`
	URL     string       `json:"url" bson:"url"`
	Snippet string       `json:"snippet" bson:"snippet"`
}

type realtimeOut struct {
	StorySummaries *storySummary `json:"storySummaries" bson:"story_summaries"`
}

type storySummary struct {
	TrendingStories []*TrendingStory `json:"trendingStories" bson:"trending_stories"`
}

// TrendingStory is a representation of realtime trend
type TrendingStory struct {
	Title    string             `json:"title" bson:"title"`
	Image    *SearchImage       `json:"image" bson:"image"`
	Articles []*TrendingArticle `json:"articles" bson:"articles"`
}

// TrendingArticle is an article relative to trending story
type TrendingArticle struct {
	Title   string `json:"articleTitle" bson:"title"`
	URL     string `json:"url" bson:"url"`
	Source  string `json:"source" bson:"source"`
	Time    string `json:"time" bson:"time"`
	Snippet string `json:"snippet" bson:"snippet"`
}

// ExploreRequest it's an input which can contain multiple items (keywords) to discover
// category can be found in ExploreCategories output
type ExploreRequest struct {
	ComparisonItems []*ComparisonItem `json:"comparisonItem" bson:"comparison_items"`
	Category        int               `json:"category" bson:"category"`
	Property        string            `json:"property" bson:"property"`
}

// ComparisonItem it's concrete search keyword
// with Geo (can be found with ExploreLocations method) locality and Time period
type ComparisonItem struct {
	Keyword                string `json:"keyword" bson:"keyword"`
	Geo                    string `json:"geo,omitempty" bson:"geo"`
	Time                   string `json:"time" bson:"time"`
	GranularTimeResolution bool   `json:"granularTimeResolution" bson:"granular_time_resolution"`
	StartTime              string `json:"startTime" bson:"start_time"`
	EndTime                string `json:"endTime" bson:"end_time"`
}

// ExploreCatTree - available categories list tree
type ExploreCatTree struct {
	Name     string            `json:"name" bson:"name"`
	ID       int               `json:"id" bson:"id"`
	Children []*ExploreCatTree `json:"children" bson:"children"`
}

// ExploreLocTree - available locations list tree
type ExploreLocTree struct {
	Name     string            `json:"name" bson:"name"`
	ID       string            `json:"id" bson:"id"`
	Children []*ExploreLocTree `json:"children" bson:"children"`
}

type exploreOut struct {
	Widgets []*ExploreWidget `json:"widgets" bson:"widgets"`
}

// ExploreWidget - output of Explore method, required for InterestOverTime, InterestByLocation and Related methods.
// Globally it's a structure related to Google Trends UI and contains mostly system info
type ExploreWidget struct {
	Token   string          `json:"token" bson:"token"`
	Type    string          `json:"type" bson:"type"`
	Title   string          `json:"title" bson:"title"`
	ID      string          `json:"id" bson:"id"`
	Request *WidgetResponse `json:"request" bson:"request"`
}

type ExploreResponse []*ExploreWidget

func (e ExploreResponse) Sort() {
	sort.Sort(e)
}

func (e ExploreResponse) Len() int {
	return len(e)
}

func (e ExploreResponse) Less(i, j int) bool {
	numI := strings.LastIndex(e[i].ID, "_")
	if numI < 0 {
		return true
	}

	numJ := strings.LastIndex(e[j].ID, "_")
	if numJ < 0 {
		return false
	}

	valI, err := strconv.ParseInt(e[i].ID[numI+1:], 10, 32)
	if err != nil {
		return true
	}

	valJ, err := strconv.ParseInt(e[j].ID[numJ+1:], 10, 32)
	if err != nil {
		return false
	}

	return valI < valJ
}

func (e ExploreResponse) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ExploreResponse) GetWidgetsByOrder(i int) ExploreResponse {
	out := make(ExploreResponse, 0)
	for _, v := range e {
		if v.ID == string(IntOverTimeWidgetID) || v.ID == string(IntOverRegionID) {
			continue
		}

		ind := strings.LastIndex(v.ID, "_")
		val, err := strconv.ParseInt(v.ID[ind+1:], 10, 32)
		if err != nil {
			return out
		}

		if int(val) == i {
			out = append(out, v)
		}
	}

	return out
}

func (e ExploreResponse) GetWidgetsByType(t WidgetType) ExploreResponse {
	out := make(ExploreResponse, 0)
	for _, v := range e {
		if strings.Contains(v.ID, string(t)) {
			out = append(out, v)
		}
	}

	return out
}

// WidgetResponse - system info for every available trends search mode
type WidgetResponse struct {
	Geo                interface{}             `json:"geo,omitempty" bson:"geo"`
	Time               string                  `json:"time,omitempty" bson:"time"`
	Resolution         string                  `json:"resolution,omitempty" bson:"resolution"`
	Locale             string                  `json:"locale,omitempty" bson:"locale"`
	Restriction        WidgetComparisonItem    `json:"restriction" bson:"restriction"`
	CompItem           []*WidgetComparisonItem `json:"comparisonItem" bson:"comparison_item"`
	RequestOpt         RequestOptions          `json:"requestOptions" bson:"request_option"`
	KeywordType        string                  `json:"keywordType" bson:"keyword_type"`
	Metric             []string                `json:"metric" bson:"metric"`
	Language           string                  `json:"language" bson:"language"`
	TrendinessSettings map[string]string       `json:"trendinessSettings" bson:"trendiness_settings"`
	DataMode           string                  `json:"dataMode,omitempty" bson:"data_mode"`
	UserConfig         map[string]string       `json:"userConfig,omitempty" bson:"user_config"`
	UserCountryCode    string                  `json:"userCountryCode,omitempty" bson:"user_country_code"`
}

// WidgetComparisonItem - system info for comparison item part of WidgetResponse
type WidgetComparisonItem struct {
	Geo                            map[string]string   `json:"geo,omitempty" bson:"geo"`
	Time                           string              `json:"time,omitempty" bson:"time"`
	ComplexKeywordsRestriction     KeywordsRestriction `json:"complexKeywordsRestriction,omitempty" bson:"complex_keywords_restriction"`
	OriginalTimeRangeForExploreURL string              `json:"originalTimeRangeForExploreUrl,omitempty" bson:"original_time_range_for_explore_url"`
}

// KeywordsRestriction - system info for keywords limitations, not used. part of WidgetResponse
type KeywordsRestriction struct {
	Keyword []*KeywordRestriction `json:"keyword" bson:"keyword"`
}

// KeywordRestriction - specific keyword limitation. Part of KeywordsRestriction
type KeywordRestriction struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

// RequestOptions - part of WidgetResponse
type RequestOptions struct {
	Property string `json:"property" bson:"property"`
	Backend  string `json:"backend" bson:"backend"`
	Category int    `json:"category" bson:"category"`
}

type multilineOut struct {
	Default multiline `json:"default" bson:"default"`
}

type multiline struct {
	TimelineData []*Timeline `json:"timelineData" bson:"timeline_data"`
}

// Timeline - it's representation of interest to trend trough period timeline. Mostly used for charts
type Timeline struct {
	Time              string   `json:"time" bson:"time"`
	FormattedTime     string   `json:"formattedTime" bson:"formatted_time"`
	FormattedAxisTime string   `json:"formattedAxisTime" bson:"formatted_axis_time"`
	Value             []int    `json:"value" bson:"value"`
	HasData           []bool   `json:"hasData" bson:"has_data"`
	FormattedValue    []string `json:"formattedValue" bson:"formatted_value"`
}

type geoOut struct {
	Default geo `json:"default" bson:"default"`
}

type geo struct {
	GeoMapData []*GeoMap `json:"geoMapData" bson:"geomap_data"`
}

// GeoMap - it's representation of interest by location. Mostly used for maps
type GeoMap struct {
	GeoCode        string   `json:"geoCode" bson:"geo_code"`
	GeoName        string   `json:"geoName" bson:"geo_name"`
	Value          []int    `json:"value" bson:"value"`
	FormattedValue []string `json:"formattedValue" bson:"formatted_value"`
	MaxValueIndex  int      `json:"maxValueIndex" bson:"max_value_index"`
	HasData        []bool   `json:"hasData" bson:"has_data"`
}

type relatedOut struct {
	Default relatedList `json:"default" bson:"default"`
}

type relatedList struct {
	Ranked []*rankedList `json:"rankedList" bson:"ranked"`
}

type rankedList struct {
	Keywords []*RankedKeyword `json:"rankedKeyword" bson:"keywords"`
}

type searchOut struct {
	Default searchList `json:"default" bson:"default"`
}

type searchList struct {
	Keywords []*KeywordTopic `json:"topics" bson:"keywords"`
}

// RankedKeyword - it's representation of related to search items
type RankedKeyword struct {
	Query          string       `json:"query,omitempty" bson:"query"`
	Topic          KeywordTopic `json:"topic,omitempty" bson:"topic"`
	Value          int          `json:"value" bson:"value"`
	FormattedValue string       `json:"formattedValue" bson:"formatted_value"`
	HasData        bool         `json:"hasData" bson:"has_data"`
	Link           string       `json:"link" bson:"link"`
}

// KeywordTopic - is a part of RankedKeyword
type KeywordTopic struct {
	Mid   string `json:"mid" bson:"mid"`
	Title string `json:"title" bson:"title"`
	Type  string `json:"type" bson:"type"`
}
