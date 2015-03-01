package waypoint

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var (
	oziExplorerHeaderRegexps = []*regexp.Regexp{
		regexp.MustCompile(`\AOziExplorer Waypoint File Version 1\.\d+\s*\z`),
		regexp.MustCompile(`\AWGS 84\s*\z`),
		regexp.MustCompile(`\AReserved 2\s*\z`),
		regexp.MustCompile(`\AReserved 3\s*\z`),
	}
	oziExplorerCommaRegexp = regexp.MustCompile(`\s*,\s*`)
)

type OziExplorerFormat struct{}

func NewOziExplorerFormat() *OziExplorerFormat {
	return &OziExplorerFormat{}
}

func (*OziExplorerFormat) Read(r io.Reader) (Collection, error) {
	var wc Collection
	scanner := bufio.NewScanner(r)
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Text()
		switch {
		case lineno <= 4:
			if oziExplorerHeaderRegexps[lineno-1].FindString(line) == "" {
				return nil, ErrSyntax{LineNo: lineno, Line: line}
			}
		default:
			ss := oziExplorerCommaRegexp.Split(line, -1)
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
			w := &T{
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

func (*OziExplorerFormat) Write(w io.Writer, wc Collection) error {
	for _, s := range []string{
		"OziExplorer Waypoint File Version 1.0\r\n",
		"WGS 84\r\n",
		"Reserved 2\r\n",
		"Reserved 3\r\n",
	} {
		if _, err := fmt.Fprintf(w, s); err != nil {
			return err
		}
	}
	for i, wp := range wc {
		if _, err := fmt.Fprintf(w, "%4d,%s,%11.6f,%11.6f,40652.2883218,0, 1, 3, 0, 65535,%-40s, 0, 0, 0, %d\r\n",
			i+1, wp.Id, wp.Latitude, wp.Longitude, wp.Description, int(wp.Altitude/0.3048)); err != nil {
			return err
		}
	}
	return nil
}
