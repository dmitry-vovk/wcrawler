package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmitry-vovk/wcrowler/crawler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configFile := "config.json"
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else if cfgPath := os.Getenv("CRAWLER_CONFIG"); cfgPath != "" {
		configFile = cfgPath
	}
	cfg, err := crawler.Read(configFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading config file %q: %s\n", configFile, err)
		os.Exit(1)
	}
	start := time.Now()
	if err = crawler.New(cfg).Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error running crawler: %s\n", err)
	} else {
		fmt.Printf("Crawler finished in %s\n", time.Since(start))
	}
}
