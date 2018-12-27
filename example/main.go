package main

import (
	"context"
	"log"

	"github.com/groovili/gogtrends"
)

func main() {
	_, err := gogtrends.Daily(context.Background(), "US")
	if err != nil {
		log.Fatal(err)
	}
}
