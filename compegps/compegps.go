package compegps

import (
	"bitbucket.org/twpayne/waypoint"

	"bufio"
	"fmt"
	"image/color"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	headerRegexps = []*regexp.Regexp{
		regexp.MustCompile(`\AG\s+WGS\s+84\s*\z`),
		regexp.MustCompile(`\AU\s+1\s*\z`),
	}
	wRegexp  = regexp.MustCompile(`\AW\s+(.{6})\s+A\s+(\d+(?:\.\d*)?).?([NS])\s+(\d+(?:\.\d*)?).?([EW])\s+\S+\s+\S+\s+(\d+(?:\.\d*)?)(.*)\z`)
	w2Regexp = regexp.MustCompile(`\Aw\s+[^,]*,[^,]*,[^,]*,(\d*),[^,]*,[^,]*,[^,]*,[^,]*,[^,]*\s*\z`)
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
		case lineno <= 1:
			if headerRegexps[lineno-1].FindString(line) == "" {
				return nil, waypoint.ErrSyntax{LineNo: lineno, Line: line}
			}
		case w == nil:
			ss := wRegexp.FindStringSubmatch(line)
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
			ss := w2Regexp.FindStringSubmatch(line)
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
	for _, s := range []string{
		"G  WGS 84\r\n",
		"U  1\r\n",
	} {
		if _, err := fmt.Fprintf(w, s); err != nil {
			return err
		}
	}
	for _, wp := range wc {
		latHemi := 'N'
		if wp.Latitude < 0 {
			latHemi = 'S'
		}
		lngHemi := 'E'
		if wp.Longitude < 0 {
			lngHemi = 'W'
		}
		// FIXME find correct format specifiers for lat and lng
		if _, err := fmt.Fprintf(w, "W  %6s A %.10f\x5c%c %.11f\x5c%c 27-MAR-62 00:00:00 %.6f %s\r\n",
			wp.Id, math.Abs(wp.Latitude), latHemi, math.Abs(wp.Longitude), lngHemi, wp.Altitude, wp.Description); err != nil {
			return err
		}
		rgb := 0xffffff
		if wp.Color != nil {
			r, g, b, _ := wp.Color.RGBA()
			rgb = int(r/0x101)<<16 + int(g/0x101)<<8 + int(b/0x101)
		}
		if _, err := fmt.Fprintf(w, "w  box,0,0.0,%d,255,1,7,,0.0\r\n", rgb); err != nil {
			return err
		}
	}
	return nil
}
