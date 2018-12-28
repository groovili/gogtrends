package main

import (
	"context"
	"github.com/groovili/gogtrends"
	log "github.com/sirupsen/logrus"
)

const locUS = "US"

func main() {
	ctx := context.Background()

	dailySearches, err := gogtrends.Daily(ctx, locUS)
	if err != nil {
		log.Fatal("Failed to get daily searches", err)
	}

	log.Println("Daily trending searches:")
	for _, v := range dailySearches {
		log.Info(v)
	}

	log.Println("Realtime trends:")
	realtime, err := gogtrends.Realtime(ctx, locUS)
	if err != nil {
		log.Fatal("Failed to get realtime trends", err)
	}

	for _, v := range realtime {
		log.Info(v)
	}
}
