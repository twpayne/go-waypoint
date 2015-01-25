package oziexplorer

import (
	"bitbucket.org/twpayne/waypoint"

	"bytes"
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
			s: "OziExplorer Waypoint File Version 1.0\r\n" +
				"WGS 84\r\n" +
				"Reserved 2\r\n" +
				"Reserved 3\r\n" +
				"   1,Aconcagua, -32.653333, -70.011667,40652.2883218,0, 1, 3, 0, 65535,Aconcagua                               , 0, 0, 0, 22841\r\n" +
				"   2,Bergneustadt,  51.051950,   7.706117,40652.2883218,0, 1, 3, 0, 65535,Bergneustadt                            , 0, 0, 0, 1601\r\n" +
				"   3,Golden Gate Bridge,  37.817500,-122.478333,40652.2883218,0, 1, 3, 0, 65535,Golden Gate Bridge                      , 0, 0, 0, 745\r\n" +
				"   4,Red Square,  55.754167,  37.620000,40652.2883218,0, 1, 3, 0, 65535,Red Square                              , 0, 0, 0, 404\r\n" +
				"   5,Sydney Opera, -33.856950, 151.215267,40652.2883218,0, 1, 3, 0, 65535,Sydney Opera                            , 0, 0, 0, 16\r\n",
			wc: waypoint.Collection{
				&waypoint.T{
					Id:          "Aconcagua",
					Latitude:    -32.653333,
					Longitude:   -70.011667,
					Altitude:    0.3048 * 22841,
					Description: "Aconcagua",
				},
				&waypoint.T{
					Id:          "Bergneustadt",
					Latitude:    51.05195,
					Longitude:   7.706117,
					Altitude:    0.3048 * 1601,
					Description: "Bergneustadt",
				},
				&waypoint.T{
					Id:          "Golden Gate Bridge",
					Latitude:    37.8175,
					Longitude:   -122.478333,
					Altitude:    0.3048 * float64(745),
					Description: "Golden Gate Bridge",
				},
				&waypoint.T{
					Id:          "Red Square",
					Latitude:    55.754167,
					Longitude:   37.62,
					Altitude:    0.3048 * 404,
					Description: "Red Square",
				},
				&waypoint.T{
					Id:          "Sydney Opera",
					Latitude:    -33.85695,
					Longitude:   151.215267,
					Altitude:    0.3048 * 16,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		if got, err := New().Read(strings.NewReader(c.s)); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if i < len(got) && !waypoint.Equal(w, got[i]) {
					t.Errorf("i=%d w=%#v got[%d]=%#v", i, w, got[i])
				}
			}
			t.Errorf("Read(strings.NewReader(%v)) == %v, %v, want %v, nil", c.s, got, err, c.wc)
		}
		w := bytes.NewBuffer(nil)
		if err := New().Write(w, c.wc); err != nil {
			t.Errorf("Write(%v) == %v. want nil", c.wc, err)
		}
		if w.String() != c.s {
			checkStrings(t, w.String(), c.s)
			t.Errorf("w.String() == %v. want %v", w.String(), c.s)
		}
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
