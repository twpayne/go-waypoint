package seeyou

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"bitbucket.org/twpayne/waypoint"
	"bitbucket.org/twpayne/waypoint/internal/dmsh"
)

var (
	latRegexp    = regexp.MustCompile(`\A\s*(\d{2})(\d{2}\.\d+)([NS])\s*\z`)
	lngRegexp    = regexp.MustCompile(`\A\s*(\d{3})(\d{2}\.\d+)([EW])\s*\z`)
	altRegexp    = regexp.MustCompile(`\A\s*(\d+(?:\.\d*)?)(m)\s*?\z`)
	headerFields = strings.Split("name,code,country,lat,lon,elev,style,rwdir,rwlen,freq,desc", ",")
)

type T struct{}

func New() *T {
	return &T{}
}

func (*T) Read(r io.Reader) (waypoint.Collection, error) {
	var wc waypoint.Collection
	csvr := csv.NewReader(r)
	csvr.FieldsPerRecord = -1
	csvr.LazyQuotes = true
	csvr.TrimLeadingSpace = true
	lineno := 0
	for {
		lineno++
		record, err := csvr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		switch lineno {
		case 1:
			if len(record) != len(headerFields) {
				return nil, waypoint.ErrSyntax{LineNo: lineno}
			}
			for i, f := range record {
				if f != headerFields[i] {
					return nil, waypoint.ErrSyntax{LineNo: lineno}
				}
			}
		default:
			if len(record) == 1 && record[0] == "-----Related Tasks-----" {
				break
			}
			if len(record) != len(headerFields) {
				continue
			}
			id := record[1]
			ss := latRegexp.FindStringSubmatch(record[3])
			if ss == nil {
				continue
			}
			latDeg, _ := strconv.ParseInt(ss[1], 10, 64)
			latMin, _ := strconv.ParseFloat(ss[2], 64)
			lat := float64(latDeg) + latMin/60
			if ss[3] == "S" {
				lat = -lat
			}
			ss = lngRegexp.FindStringSubmatch(record[4])
			if ss == nil {
				continue
			}
			lngDeg, _ := strconv.ParseInt(ss[1], 10, 64)
			lngMin, _ := strconv.ParseFloat(ss[2], 64)
			lng := float64(lngDeg) + lngMin/60
			if ss[3] == "W" {
				lng = -lng
			}
			ss = altRegexp.FindStringSubmatch(record[5])
			if ss == nil {
				continue
			}
			alt, _ := strconv.ParseFloat(ss[1], 64)
			if ss[2] != "m" {
				alt *= 0.3048
			}
			description := record[10]
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
	return wc, nil
}

func (*T) Write(w io.Writer, wc waypoint.Collection) error {
	csvw := csv.NewWriter(w)
	csvw.UseCRLF = true
	if err := csvw.Write(headerFields); err != nil {
		return err
	}
	record := make([]string, 10)
	for _, wp := range wc {
		// FIXME record[0] = wp.Name
		record[1] = wp.Id
		latDeg, latMin, latHemi := dmsh.DMH(wp.Latitude, "NS")
		record[3] = fmt.Sprintf("%02d%06.3f%c", latDeg, latMin, latHemi)
		lngDeg, lngMin, lngHemi := dmsh.DMH(wp.Longitude, "EW")
		record[4] = fmt.Sprintf("%02d%06.3f%c", lngDeg, lngMin, lngHemi)
		record[5] = fmt.Sprintf("%.1fm", wp.Altitude)
		record[10] = wp.Description
		if err := csvw.Write(record); err != nil {
			return err
		}
	}
	return nil
}
