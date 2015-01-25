package formatgeo

import (
	"bitbucket.org/twpayne/waypoint"
	"bitbucket.org/twpayne/waypoint/internal/dmsh"

	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	idRegExp = regexp.MustCompile(`\A\$FormatGEO\s*\z`)
	regExp   = regexp.MustCompile(`\A(.*)\s+([NS])\s+(\d+)\s+(\d+)\s+(\d+\.\d+)\s+([EW])\s+(\d+)\s+(\d+)\s+(\d+\.\d+)\s+(-?\d+)(?:\s+(.*))?\z`)
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
		case lineno == 1:
			if idRegExp.FindString(line) == "" {
				return nil, waypoint.ErrSyntax{LineNo: lineno, Line: line}
			}
		default:
			ss := regExp.FindStringSubmatch(line)
			if ss == nil {
				continue
			}
			latDeg, _ := strconv.ParseInt(ss[3], 10, 64)
			latMin, _ := strconv.ParseInt(ss[4], 10, 64)
			latSec, _ := strconv.ParseFloat(ss[5], 64)
			lat := dmsh.D(int(latDeg), int(latMin), latSec)
			if ss[2] == "S" {
				lat = -lat
			}
			lngDeg, _ := strconv.ParseInt(ss[7], 10, 64)
			lngMin, _ := strconv.ParseInt(ss[8], 10, 64)
			lngSec, _ := strconv.ParseFloat(ss[9], 64)
			lng := dmsh.D(int(lngDeg), int(lngMin), lngSec)
			if ss[6] == "W" {
				lng = -lng
			}
			alt, _ := strconv.ParseFloat(ss[10], 64)
			id := strings.TrimSpace(ss[1])
			description := strings.TrimSpace(ss[11])
			w := &waypoint.T{
				Id:          id,
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

func (*T) Write(w io.Writer, wc waypoint.Collection) error {
	if _, err := fmt.Fprintf(w, "$FormatGEO\r\n"); err != nil {
		return err
	}
	for _, wp := range wc {
		latDeg, latMin, latSec, latHemi := dmsh.DMSH(wp.Latitude, "NS")
		lngDeg, lngMin, lngSec, lngHemi := dmsh.DMSH(wp.Longitude, "EW")
		if _, err := fmt.Fprintf(w, "%-8s  %s %02d %02d %05.2f    %s %03d %02d %05.2f  %4d  %s\r\n", wp.Id, latHemi, latDeg, latMin, latSec, lngHemi, lngDeg, lngMin, lngSec, int(wp.Altitude), wp.Description); err != nil {
			return err
		}
	}
	return nil
}
