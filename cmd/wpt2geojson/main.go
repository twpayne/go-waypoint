package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/twpayne/go-waypoint"
)

func readWaypoints(filename string) (waypoint.Collection, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	wc, _, err := waypoint.Read(f)
	return wc, err
}

func run() error {
	flag.Parse()

	var wcs []*waypoint.T
	for _, arg := range flag.Args() {
		wc, err := readWaypoints(arg)
		if err != nil {
			return err
		}
		wcs = append(wcs, wc...)
	}

	return waypoint.NewGeoJSONFormat().Write(os.Stdout, wcs)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
