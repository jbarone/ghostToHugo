package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	ghostToHugo "github.com/jbarone/ghostToHugo/lib"
)

// Print usage information
func usage() {
	fmt.Printf("Usage: %s [OPTIONS] <Ghost Export>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	var path, loc, format string

	flag.Usage = usage

	flag.StringVarP(&path, "hugo", "p", ".", "Path to hugo project")
	flag.StringVarP(&loc, "location", "l", "",
		"Location to use for time conversions (default: local)")
	flag.StringVarP(&format, "dateformat", "f", "",
		"Date format string to use for time conversions (default: RFC3339)")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	opts := []func(*ghostToHugo.GhostToHugo){
		ghostToHugo.WithHugoPath(path),
	}
	if loc != "" {
		location, err := time.LoadLocation(loc)
		if err != nil {
			log.Fatalf("Error loading location %s: %v", loc, err)
		}
		opts = append(opts, ghostToHugo.WithLocation(location))
	}

	if format != "" {
		opts = append(opts, ghostToHugo.WithDateFormat(format))
	}

	gth, err := ghostToHugo.NewGhostToHugo(opts...)
	if err != nil {
		log.Fatalf("Error initializing converter (%v)", err)
	}

	gth.Export(flag.Arg(0))
}
