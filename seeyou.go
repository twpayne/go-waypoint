package waypoint

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	seeYouLatRegexp    = regexp.MustCompile(`\A\s*(\d{2})(\d{2}\.\d+)([NS])\s*\z`)
	seeYouLngRegexp    = regexp.MustCompile(`\A\s*(\d{3})(\d{2}\.\d+)([EW])\s*\z`)
	seeYouAltRegexp    = regexp.MustCompile(`\A\s*(\d+(?:\.\d*)?)(m)\s*?\z`)
	seeYouHeaderFields = strings.Split("name,code,country,lat,lon,elev,style,rwdir,rwlen,freq,desc", ",")
)

type SeeYouFormat struct{}

func NewSeeYouFormat() *SeeYouFormat {
	return &SeeYouFormat{}
}

func (*SeeYouFormat) Extension() string {
	return "cup"
}

func (*SeeYouFormat) Name() string {
	return "seeyou"
}

func (*SeeYouFormat) Read(r io.Reader) (Collection, error) {
	var wc Collection
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
			if len(record) != len(seeYouHeaderFields) {
				return nil, errSyntax{LineNo: lineno}
			}
			for i, f := range record {
				if f != seeYouHeaderFields[i] {
					return nil, errSyntax{LineNo: lineno}
				}
			}
		default:
			if len(record) == 1 && record[0] == "-----Related Tasks-----" {
				break
			}
			if len(record) != len(seeYouHeaderFields) {
				continue
			}
			id := record[1]
			ss := seeYouLatRegexp.FindStringSubmatch(record[3])
			if ss == nil {
				continue
			}
			latDeg, _ := strconv.ParseInt(ss[1], 10, 64)
			latMin, _ := strconv.ParseFloat(ss[2], 64)
			lat := float64(latDeg) + latMin/60
			if ss[3] == "S" {
				lat = -lat
			}
			ss = seeYouLngRegexp.FindStringSubmatch(record[4])
			if ss == nil {
				continue
			}
			lngDeg, _ := strconv.ParseInt(ss[1], 10, 64)
			lngMin, _ := strconv.ParseFloat(ss[2], 64)
			lng := float64(lngDeg) + lngMin/60
			if ss[3] == "W" {
				lng = -lng
			}
			ss = seeYouAltRegexp.FindStringSubmatch(record[5])
			if ss == nil {
				continue
			}
			alt, _ := strconv.ParseFloat(ss[1], 64)
			if ss[2] != "m" {
				alt *= 0.3048
			}
			description := record[10]
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
	return wc, nil
}

func (*SeeYouFormat) Write(w io.Writer, wc Collection) error {
	csvw := csv.NewWriter(w)
	csvw.UseCRLF = true
	if err := csvw.Write(seeYouHeaderFields); err != nil {
		return err
	}
	record := make([]string, 11)
	for _, wp := range wc {
		record[0] = wp.ID
		record[1] = wp.ID
		latDeg, latMin, latHemi := DMH(wp.Latitude, NS)
		record[3] = fmt.Sprintf("%02d%06.3f%c", latDeg, latMin, latHemi)
		lngDeg, lngMin, lngHemi := DMH(wp.Longitude, EW)
		record[4] = fmt.Sprintf("%02d%06.3f%c", lngDeg, lngMin, lngHemi)
		record[5] = fmt.Sprintf("%.1fm", wp.Altitude)
		record[10] = wp.Description
		if err := csvw.Write(record); err != nil {
			return err
		}
	}
	csvw.Flush()
	return csvw.Error()
}
