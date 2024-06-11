package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kanatsanan6/tldr/cache"
)

const (
	remoteUrl = "https://tldr.sh/assets/tldr.zip"
	ttl       = 24 * time.Hour
)

func printPage(page string) {
	_, err := cache.NewRepository(remoteUrl, ttl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(page)
}

func main() {
	flag.Parse()

	printPage(flag.Arg(0))
}
