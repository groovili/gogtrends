package gogtrends

const (
	gAPI = "https://trends.google.com/trends/api"

	gDaily    = "/dailytrends"
	gRealtime = "/realtimetrends"

	gSExplore     = "/explore"
	gSCategories  = "/explore/pickers/category"
	gSGeo         = "/explore/pickers/geo"
	gSRelated     = "/widgetdata/relatedsearches"
	gSIntOverTime = "/widgetdata/multiline"
	gSIntOverReg  = "/widgetdata/comparedgeo"

	paramHl    = "hl"
	paramCat   = "cat"
	paramGeo   = "geo"
	paramReq   = "req"
	paramTZ    = "tz"
	paramToken = "token"

	intOverTimeWidgetID = "TIMESERIES"
	intOverRegionID     = "GEO_MAP"
	relatedQueriesID    = "RELATED_QUERIES"
	relatedTopicsID     = "RELATED_TOPICS"

	compareDataMode = "PERCENTAGES"
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

type TrendingSearch struct {
	Title            *SearchTitle     `json:"title" bson:"title"`
	FormattedTraffic string           `json:"formattedTraffic" bson:"formatted_traffic"`
	Image            *SearchImage     `json:"image" bson:"image"`
	Articles         []*SearchArticle `json:"articles" bson:"articles"`
}

type SearchTitle struct {
	Query string `json:"query" bson:"query"`
}

type SearchImage struct {
	NewsURL  string `json:"newsUrl" bson:"news_url"`
	Source   string `json:"source" bson:"source"`
	ImageURL string `json:"imageUrl" bson:"image_url"`
}

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

type TrendingStory struct {
	Title    string             `json:"title" bson:"title"`
	Image    *SearchImage       `json:"image" bson:"image"`
	Articles []*TrendingArticle `json:"articles" bson:"articles"`
}

type TrendingArticle struct {
	Title   string `json:"articleTitle" bson:"title"`
	URL     string `json:"url" bson:"url"`
	Source  string `json:"source" bson:"source"`
	Time    string `json:"time" bson:"time"`
	Snippet string `json:"snippet" bson:"snippet"`
}

type ExploreRequest struct {
	ComparisonItems []*ComparisonItem `json:"comparisonItem" bson:"comparison_items"`
	Category        int               `json:"category" bson:"category"`
	Property        string            `json:"property" bson:"property"`
}

type ComparisonItem struct {
	Keyword string `json:"keyword" bson:"keyword"`
	Geo     string `json:"geo,omitempty" bson:"geo"`
	Time    string `json:"time" bson:"time"`
}

type ExploreCatTree struct {
	Name     string            `json:"name" bson:"name"`
	ID       int               `json:"id" bson:"id"`
	Children []*ExploreCatTree `json:"children" bson:"children"`
}

type ExploreLocTree struct {
	Name     string            `json:"name" bson:"name"`
	ID       string            `json:"id" bson:"id"`
	Children []*ExploreLocTree `json:"children" bson:"children"`
}

type exploreOut struct {
	Widgets []*ExploreWidget `json:"widgets" bson:"widgets"`
}

type ExploreWidget struct {
	Token   string          `json:"token" bson:"token"`
	Type    string          `json:"type" bson:"type"`
	Title   string          `json:"title" bson:"title"`
	ID      string          `json:"id" bson:"id"`
	Request *WidgetResponse `json:"request" bson:"request"`
}

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
	UserCountryCode    string                  `json:"userCountryCode,omitempty" bson:"user_country_code"`
}

type WidgetComparisonItem struct {
	Geo                            map[string]string   `json:"geo,omitempty" bson:"geo"`
	Time                           string              `json:"time,omitempty" bson:"time"`
	ComplexKeywordsRestriction     KeywordsRestriction `json:"complexKeywordsRestriction,omitempty" bson:"complex_keywords_restriction"`
	OriginalTimeRangeForExploreURL string              `json:"originalTimeRangeForExploreUrl,omitempty" bson:"original_time_range_for_explore_url"`
}

type KeywordsRestriction struct {
	Keyword []*KeywordRestriction `json:"keyword" bson:"keyword"`
}

type KeywordRestriction struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

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

type RankedKeyword struct {
	Query          string       `json:"query,omitempty" bson:"query"`
	Topic          KeywordTopic `json:"topic,omitempty" bson:"topic"`
	Value          int          `json:"value" bson:"value"`
	FormattedValue string       `json:"formattedValue" bson:"formatted_value"`
	HasData        bool         `json:"hasData" bson:"has_data"`
	Link           string       `json:"link" bson:"link"`
}

type KeywordTopic struct {
	Mid   string `json:"mid" bson:"mid"`
	Title string `json:"title" bson:"title"`
	Type  string `json:"type" bson:"type"`
}
