package seeyou

import (
	"bitbucket.org/twpayne/waypoint"
	"bitbucket.org/twpayne/waypoint/internal/dmsh"

	//"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestReadWrite(t *testing.T) {
	for _, c := range []struct {
		s  string
		wc waypoint.Collection
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
			wc: waypoint.Collection{
				&waypoint.T{
					Id:          "Aconcagua",
					Latitude:    -dmsh.D(32, 39.2, 0),
					Longitude:   -dmsh.D(70, 0.7, 0),
					Altitude:    6962,
					Description: "Highest mountain in south-america",
				},
				&waypoint.T{
					Id:          "",
					Latitude:    dmsh.D(51, 3.117, 0),
					Longitude:   dmsh.D(7, 42.367, 0),
					Altitude:    488,
					Description: "Rabbit holes, 20\" ditch south end of rwy",
				},
				&waypoint.T{
					Id:          "GGB",
					Latitude:    dmsh.D(37, 49.05, 0),
					Longitude:   -dmsh.D(122, 28.7, 0),
					Altitude:    227,
					Description: "",
				},
				&waypoint.T{
					Id:          "RedSqr",
					Latitude:    dmsh.D(55, 45.25, 0),
					Longitude:   dmsh.D(37, 37.2, 0),
					Altitude:    123,
					Description: "",
				},
				&waypoint.T{
					Id:          "Opera",
					Latitude:    -dmsh.D(33, 51.417, 0),
					Longitude:   dmsh.D(151, 12.916, 0),
					Altitude:    5,
					Description: "",
				},
			},
		},
	} {
		if got, err := New().Read(strings.NewReader(c.s)); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if i < len(got) && !waypoint.Equal(w, got[i]) {
					t.Errorf("i=%d w=%#v got[%d]=%#v", i, w, *got[i])
				}
			}
			t.Errorf("Read(strings.NewReader(%v)) == %v, %v, want %v, nil", c.s, got, err, c.wc)
		}
		/*
			w := bytes.NewBuffer(nil)
			if err := New().Write(w, c.wc); err != nil {
				t.Errorf("Write(%v) == %v. want nil", c.wc, err)
			}
			if w.String() != c.s {
				checkStrings(t, w.String(), c.s)
				t.Errorf("w.String() == %v. want %v", w.String(), c.s)
			}
		*/
	}
}

func checkStrings(t *testing.T, s1, s2 string) {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}
	line := 1
	col := 0
	for i := 0; i < n; i++ {
		col++
		if s1[i] != s2[i] {
			t.Errorf("strings differ a line %d column %d (%v != %v)", line, col, s1[i], s2[i])
			break
		}
		if s1[i] == '\n' {
			line++
			col = 0
		}
	}
}
