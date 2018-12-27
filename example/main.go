package main

import (
	"context"
	"github.com/groovili/gogtrends"
	log "github.com/sirupsen/logrus"
)

func main() {
	dailySearches, err := gogtrends.Daily(context.Background(), "US")
	if err != nil {
		log.Fatal("Failed to get daily searches", err)
	}

	log.Println("Daily trending searches:")
	for _, v := range dailySearches {
		log.Info(v)
	}
}
