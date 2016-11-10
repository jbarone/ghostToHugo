package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/jbarone/ghostToHugo/lib/ghost"
	"github.com/jbarone/ghostToHugo/lib/hugo"
)

// Print usage information
func usage() {
	fmt.Printf("Usage: %s [OPTIONS] <Ghost Export>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	var c hugo.Config
	var l, f string

	flag.Usage = usage

	flag.StringVarP(&c.Path, "hugo", "p", ".", "Path to hugo project")
	flag.StringVarP(&l, "location", "l", "",
		"Location to use for time conversions (default: local)")
	flag.StringVarP(&f, "dateformat", "f", "",
		"Date format string to use for time conversions (default: RFC3339)")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if err := hugo.Init(c); err != nil {
		log.Fatalf("Error initializing Hugo Config (%v)", err)
	}

	if l != "" {
		location, err := time.LoadLocation(l)
		if err != nil {
			log.Fatalf("Error loading location %s: %v", l, err)
		}
		ghost.SetLocation(location)
	}

	if f != "" {
		ghost.SetDateFormat(f)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error opening export: %v", err)
	}
	defer file.Close()

	reader := ghost.ExportReader{file}

	entries, err := ghost.Process(reader)
	if err != nil {
		log.Fatalf("Error processing Ghost export: %v", err)
	}

	var wg sync.WaitGroup
	for _, entry := range entries {
		wg.Add(1)
		go func(data ghost.ExportData) {
			defer wg.Done()
			hugo.ExportGhost(&data)
		}(entry.Data)
	}

	wg.Wait()
}
