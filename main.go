package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kanatsanan6/tldr/cache"
)

const remoteUrl = "https://tldr.sh/assets/tldr.zip"

func printPage(page string) {
	_, err := cache.NewRepository(remoteUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(page)
}

func main() {
	flag.Parse()

	printPage(flag.Arg(0))
}
