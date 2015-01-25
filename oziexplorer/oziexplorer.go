package oziexplorer

import (
	"bitbucket.org/twpayne/waypoint"

	"bufio"
	"io"
	"regexp"
	"strconv"
)

var (
	headerRegexps = []*regexp.Regexp{
		regexp.MustCompile(`\AOziExplorer Waypoint File Version 1\.\d+\s*\z`),
		regexp.MustCompile(`\AWGS 84\s*\z`),
		regexp.MustCompile(`\AReserved 2\s*\z`),
		regexp.MustCompile(`\AReserved 3\s*\z`),
	}
	commaRegexp = regexp.MustCompile(`\s*,\s*`)
)

type T struct{}

func New() *T {
	return &T{}
}

func (*T) Read(r io.Reader) (waypoint.Collection, error) {
	var wc waypoint.Collection
	scanner := bufio.NewScanner(r)
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Text()
		switch {
		case lineno <= 4:
			if headerRegexps[lineno-1].FindString(line) == "" {
				return nil, waypoint.ErrSyntax{LineNo: lineno, Line: line}
			}
		default:
			ss := commaRegexp.Split(line, -1)
			if len(ss) < 15 {
				continue
			}
			id := ss[1]
			lat, err := strconv.ParseFloat(ss[2], 64)
			if err != nil {
				continue
			}
			lng, err := strconv.ParseFloat(ss[3], 64)
			if err != nil {
				continue
			}
			description := ss[10]
			alt, err := strconv.ParseInt(ss[14], 10, 64)
			if err != nil {
				continue
			}
			w := &waypoint.T{
				Id:          id,
				Latitude:    lat,
				Longitude:   lng,
				Altitude:    0.3048 * float64(alt),
				Description: description,
			}
			wc = append(wc, w)
		}
	}
	return wc, scanner.Err()
}
