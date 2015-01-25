package compegps

import (
	"bitbucket.org/twpayne/waypoint"

	"bufio"
	"image/color"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	gRegExp  = regexp.MustCompile(`\AG\s+WGS\s+84\s*\z`)
	uRegExp  = regexp.MustCompile(`\AU\s+1\s*\z`)
	wRegExp  = regexp.MustCompile(`\AW\s+(.{6})\s+A\s+(\d+(?:\.\d*)?).?([NS])\s+(\d+(?:\.\d*)?).?([EW])\s+\S+\s+\S+\s+(\d+(?:\.\d*)?)(.*)\z`)
	w2RegExp = regexp.MustCompile(`\Aw\s+[^,]*,[^,]*,[^,]*,(\d*),[^,]*,[^,]*,[^,]*,[^,]*,[^,]*\s*\z`)
)

type T struct{}

func New() *T {
	return &T{}
}

func (*T) Read(r io.Reader) (waypoint.Collection, error) {
	var wc waypoint.Collection
	scanner := bufio.NewScanner(r)
	lineno := 0
	var w *waypoint.T
	for scanner.Scan() {
		lineno++
		line := scanner.Text()
		switch {
		case lineno == 1:
			if gRegExp.FindString(line) == "" {
				return nil, waypoint.ErrSyntax{LineNo: lineno, Line: line}
			}
		case lineno == 2:
			if uRegExp.FindString(line) == "" {
				return nil, waypoint.ErrSyntax{LineNo: lineno, Line: line}
			}
		case w == nil:
			ss := wRegExp.FindStringSubmatch(line)
			if ss == nil {
				continue
			}
			id := ss[1]
			lat, _ := strconv.ParseFloat(ss[2], 64)
			if ss[3] == "S" {
				lat = -lat
			}
			lng, _ := strconv.ParseFloat(ss[4], 64)
			if ss[5] == "W" {
				lng = -lng
			}
			alt, _ := strconv.ParseFloat(ss[6], 64)
			description := strings.TrimSpace(ss[7])
			w = &waypoint.T{
				Id:          id,
				Latitude:    lat,
				Longitude:   lng,
				Altitude:    alt,
				Description: description,
			}
		default:
			ss := w2RegExp.FindStringSubmatch(line)
			if ss != nil {
				rgb, _ := strconv.ParseInt(ss[1], 10, 64)
				w.Color = color.RGBA{
					R: uint8(rgb >> 16),
					G: uint8(rgb >> 8),
					B: uint8(rgb),
				}
			}
			wc = append(wc, w)
			w = nil
		}
	}
	if w != nil {
		wc = append(wc, w)
	}
	return wc, scanner.Err()
}

func (*T) Write(w io.Writer, wc waypoint.Collection) error {
	return nil
}
