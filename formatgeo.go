package waypoint

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	formatGeoIDRegexp   = regexp.MustCompile(`\A\$FormatGEO\s*\z`)
	formatGeoLineRegexp = regexp.MustCompile(`\A(.*)\s+([NS])\s+(\d+)\s+(\d+)\s+(\d+\.\d+)\s+([EW])\s+(\d+)\s+(\d+)\s+(\d+\.\d+)\s+(-?\d+)(?:\s+(.*))?\z`)
)

// A FormatGeoFormat is a GeoFormat format.
type FormatGeoFormat struct{}

// NewFormatGeoFormat returns a new FormatGeoFormat.
func NewFormatGeoFormat() *FormatGeoFormat {
	return &FormatGeoFormat{}
}

// Extension returns f's extension.
func (f *FormatGeoFormat) Extension() string {
	return "wpt"
}

// Name returns f's name.
func (f *FormatGeoFormat) Name() string {
	return "formatgeo"
}

// Read reads a Collection from r.
func (f *FormatGeoFormat) Read(r io.Reader) (Collection, error) {
	var wc Collection
	scanner := bufio.NewScanner(r)
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Text()
		switch {
		case lineno == 1:
			if formatGeoIDRegexp.FindString(line) == "" {
				return nil, errSyntax{LineNo: lineno, Line: line}
			}
		default:
			ss := formatGeoLineRegexp.FindStringSubmatch(line)
			if ss == nil {
				continue
			}
			latDeg, _ := strconv.ParseInt(ss[3], 10, 64)
			latMin, _ := strconv.ParseInt(ss[4], 10, 64)
			latSec, _ := strconv.ParseFloat(ss[5], 64)
			lat := D(float64(latDeg), float64(latMin), latSec)
			if ss[2] == "S" {
				lat = -lat
			}
			lngDeg, _ := strconv.ParseInt(ss[7], 10, 64)
			lngMin, _ := strconv.ParseInt(ss[8], 10, 64)
			lngSec, _ := strconv.ParseFloat(ss[9], 64)
			lng := D(float64(lngDeg), float64(lngMin), lngSec)
			if ss[6] == "W" {
				lng = -lng
			}
			alt, _ := strconv.ParseFloat(ss[10], 64)
			id := strings.TrimSpace(ss[1])
			description := strings.TrimSpace(ss[11])
			w := &T{
				ID:          id,
				Latitude:    lat,
				Longitude:   lng,
				Altitude:    alt,
				Description: description,
			}
			wc = append(wc, w)
		}
	}
	return wc, scanner.Err()
}

// Write writes wc to w.
func (f *FormatGeoFormat) Write(w io.Writer, wc Collection) error {
	if _, err := fmt.Fprintf(w, "$FormatGEO\r\n"); err != nil {
		return err
	}
	for _, wp := range wc {
		latDeg, latMin, latSec, latHemi := dmsh(wp.Latitude, ns)
		lngDeg, lngMin, lngSec, lngHemi := dmsh(wp.Longitude, ew)
		if _, err := fmt.Fprintf(w, "%-8s  %c %02d %02d %05.2f    %c %03d %02d %05.2f  %4d  %s\r\n", wp.ID, latHemi, latDeg, latMin, latSec, lngHemi, lngDeg, lngMin, lngSec, int(wp.Altitude), wp.Description); err != nil {
			return err
		}
	}
	return nil
}
