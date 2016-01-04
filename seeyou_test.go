package waypoint

import (
	//"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestSeeYouReadWrite(t *testing.T) {
	for _, c := range []struct {
		s  string
		wc Collection
	}{
		{
			s: "name,code,country,lat,lon,elev,style,rwdir,rwlen,freq,desc\r\n" +
				"\"Aconcagua\",\"Aconcagua\",,3239.200S,07000.700W,6962.0m,7,0,0.0m,\"\",\"Highest mountain in south-america\"\r\n" +
				//"\"Bergneustadt\",\"\",,5103.117N  ,00742.367E,  488.0m,  5  ,040,590m,\"123.650\" , \"Rabbit holes, 20\" ditch south end of rwy\"\r\n" +
				"\"Bergneustadt\",\"\",,5103.117N  ,00742.367E,  488.0m,  5  ,040,590m,\"123.650\", \"Rabbit holes, 20\" ditch south end of rwy\"\r\n" +
				"\"Golden Gate Bridge\",\"GGB\",,3749.050N,12228.700W,227.0m,14,0,0.005NM,\"\",\"\"\r\n" +
				"\"Red Square\",\"RedSqr\",,5545.250N,03737.200E,123.0m,3,90,0.01ml,\"\",\"\"\r\n" +
				"\"Sydney Opera\",\"Opera\",,3351.417S,15112.916E,5.0m,1,0,0.0m,\"\",\"\"\r\n" +
				"-----Related Tasks-----\r\n",
			wc: Collection{
				&T{
					Id:          "Aconcagua",
					Latitude:    -D(32, 39.2, 0),
					Longitude:   -D(70, 0.7, 0),
					Altitude:    6962,
					Description: "Highest mountain in south-america",
				},
				&T{
					Id:          "",
					Latitude:    D(51, 3.117, 0),
					Longitude:   D(7, 42.367, 0),
					Altitude:    488,
					Description: "Rabbit holes, 20\" ditch south end of rwy",
				},
				&T{
					Id:          "GGB",
					Latitude:    D(37, 49.05, 0),
					Longitude:   -D(122, 28.7, 0),
					Altitude:    227,
					Description: "",
				},
				&T{
					Id:          "RedSqr",
					Latitude:    D(55, 45.25, 0),
					Longitude:   D(37, 37.2, 0),
					Altitude:    123,
					Description: "",
				},
				&T{
					Id:          "Opera",
					Latitude:    -D(33, 51.417, 0),
					Longitude:   D(151, 12.916, 0),
					Altitude:    5,
					Description: "",
				},
			},
		},
	} {
		if got, err := NewSeeYouFormat().Read(strings.NewReader(c.s)); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if err := Equal(w, got[i]); err != nil {
					t.Errorf("want %#v got=%#v, %v", w, got[i], err)
				}
			}
			t.Errorf("Read(strings.NewReader(%v)) == %v, %v, want %v, nil", c.s, got, err, c.wc)
		}
		_, f, err := Read(strings.NewReader(c.s))
		if err != nil {
			t.Errorf("Read(...) return %v, expected nil", err)
		}
		if _, ok := f.(*SeeYouFormat); !ok {
			t.Errorf("Read(...) returned a %T, expected a SeeYouFormat", f)
		}
		/*
			w := bytes.NewBuffer(nil)
			if err := NewSeeYouFormat().Write(w, c.wc); err != nil {
				t.Errorf("Write(%v) == %v. want nil", c.wc, err)
			}
			if w.String() != c.s {
				checkStrings(t, w.String(), c.s)
				t.Errorf("w.String() == %v. want %v", w.String(), c.s)
			}
		*/
	}
}
