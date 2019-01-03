package main

import (
	"context"
	"github.com/groovili/gogtrends"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	locUS  = "US"
	catAll = "all"
	langEn = "EN"
)

var sg = new(sync.WaitGroup)

func main() {
	ctx := context.Background()

	dailySearches, err := gogtrends.Daily(ctx, langEn, locUS)
	if err != nil {
		log.Fatal("Failed to get daily searches", err)
	}

	log.Println("Daily trending searches:")
	for _, v := range dailySearches {
		log.Info(v)
	}

	log.Println("Realtime trends:")
	realtime, err := gogtrends.Realtime(ctx, langEn, locUS, catAll)
	if err != nil {
		log.Fatal("Failed to get realtime trends", err)
	}

	for _, v := range realtime {
		log.Info(v)
	}

	cats, err := gogtrends.ExploreCategories(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Available explore categories:")
	for _, v := range cats.Children {
		log.Println(v.Name, v.ID)
		sg.Add(1)
		go printNestedItems(v.Children)
	}
	sg.Wait()

	log.Info("Explore trends:")

	gogtrends.Debug(true)

	explore, err := gogtrends.Explore(ctx, &gogtrends.ExploreRequest{
		ComparisonItems: []*gogtrends.ComparisonItem{
			{
				Keyword: "Golang",
				Geo:     locUS,
				Time:    "today+12-m",
			},
		},
		Category: 31, // Programming category
		Property: "",
	}, langEn)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range explore {
		log.Info(v)
	}
}

func printNestedItems(cats []*gogtrends.ExploreCategoriesTree) {
	defer sg.Done()
	for _, v := range cats {
		log.Println(v.Name, v.ID)
		if len(v.Children) > 0 {
			sg.Add(1)
			go printNestedItems(v.Children)
		}
	}
}
