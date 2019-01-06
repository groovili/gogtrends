package main

import (
	"context"
	"reflect"
	"sync"

	"github.com/groovili/gogtrends"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	locUS  = "US"
	catAll = "all"
	langEn = "EN"
)

var sg = new(sync.WaitGroup)

func main() {
	//Enable debug to see request-response
	//gogtrends.Debug(true)

	ctx := context.Background()

	log.Println("Daily trending searches:")
	dailySearches, err := gogtrends.Daily(ctx, langEn, locUS)
	handleError(err, "Failed to get daily searches")
	printItems(dailySearches)

	log.Println("Realtime trends:")
	realtime, err := gogtrends.Realtime(ctx, langEn, locUS, catAll)
	handleError(err, "Failed to get realtime trends")
	printItems(realtime)

	log.Info("Available explore categories:")
	cats, err := gogtrends.ExploreCategories(ctx)
	handleError(err, "Failed to explore categories")

	// recursive print of categories tree
	// do it concurrent to increase execution speed
	for _, v := range cats.Children {
		log.Println(v.Name, v.ID)
		sg.Add(1)
		go printNestedItems(v.Children)
	}
	sg.Wait()

	log.Info("Explore trends:")
	// get widgets for Golang keyword in programming category
	explore, err := gogtrends.Explore(ctx, &gogtrends.ExploreRequest{
		ComparisonItems: []*gogtrends.ComparisonItem{
			{
				Keyword: "Go",
				Geo:     locUS,
				Time:    "today+12-m",
			},
		},
		Category: 31, // Programming category
		Property: "",
	}, langEn)
	handleError(err, "Failed to explore widgets")
	printItems(explore)

	log.Info("Interest over time:")
	overTime, err := gogtrends.InterestOverTime(ctx, explore[0], langEn)
	handleError(err, "Failed in call interest over time")
	printItems(overTime)

	log.Info("Interest by location:")
	overReg, err := gogtrends.InterestByLocation(ctx, explore[1], langEn)
	handleError(err, "Failed in call interest by location")
	printItems(overReg)

	log.Info("Related topics:")
	relT, err := gogtrends.Related(ctx, explore[2], langEn)
	handleError(err, "Failed to get related topics")
	printItems(relT)

	log.Info("Related queries:")
	relQ, err := gogtrends.Related(ctx, explore[3], langEn)
	handleError(err, "Failed to get related queries")
	printItems(relQ)
}

func handleError(err error, errMsg string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, errMsg))
	}
}

func printItems(items interface{}) {
	ref := reflect.ValueOf(items)

	if ref.Kind() != reflect.Slice {
		log.Fatalf("Failed to print %s. It's not a slice type.", ref.Kind())
	}

	for i := 0; i < ref.Len(); i++ {
		log.Println(ref.Index(i).Interface())
	}
}

func printNestedItems(cats []*gogtrends.ExploreCatTree) {
	defer sg.Done()
	for _, v := range cats {
		log.Println(v.Name, v.ID)
		if len(v.Children) > 0 {
			sg.Add(1)
			go printNestedItems(v.Children)
		}
	}
}
