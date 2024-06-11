package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kanatsanan6/tldr/cache"
	"github.com/kanatsanan6/tldr/utils"
)

const (
	remoteUrl = "https://tldr.sh/assets/tldr.zip"
	ttl       = 24 * time.Hour
	platform  = "osx" // only support macos for now
)

func printPage(page string) {
	repo, err := cache.NewRepository(remoteUrl, ttl)
	if err != nil {
		log.Fatal(err)
	}

	markdown, err := repo.MarkDown(platform, page)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("The page doesn't exist")
		} else {
			log.Fatal(err)
		}
		return
	}
	defer markdown.Close()

	render, err := utils.Render(markdown)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(render)
}

func main() {
	flag.Parse()

	printPage(flag.Arg(0))
}
