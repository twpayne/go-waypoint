package compegps

import (
	"bitbucket.org/twpayne/waypoint"

	"bytes"
	"image/color"
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
			s: "G  WGS 84\r\n" +
				"U  1\r\n" +
				"W  ACONCA A 32.6533333333\x5cS 70.0116666667\x5cW 27-MAR-62 00:00:00 6962.000000 Highest mountain in south-america\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  BERGNE A 51.0519500000\x5cN 7.70611700000\x5cE 27-MAR-62 00:00:00 488.000000 Rabbit holes, 20\" ditch south end of rwy\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  GOLDEN A 37.8175000000\x5cN 122.478333333\x5cW 27-MAR-62 00:00:00 227.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  REDSQU A 55.7541670000\x5cN 37.6200000000\x5cE 27-MAR-62 00:00:00 123.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  SYDNEY A 33.8569500000\x5cS 151.215267000\x5cE 27-MAR-62 00:00:00 5.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n",
			wc: waypoint.Collection{
				&waypoint.T{
					Id:          "ACONCA",
					Latitude:    -32.6533333333,
					Longitude:   -70.0116666667,
					Altitude:    6962,
					Description: "Highest mountain in south-america",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&waypoint.T{
					Id:          "BERGNE",
					Latitude:    51.05195,
					Longitude:   7.706117,
					Altitude:    488,
					Description: "Rabbit holes, 20\" ditch south end of rwy",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&waypoint.T{
					Id:          "GOLDEN",
					Latitude:    37.8175,
					Longitude:   -122.478333333,
					Altitude:    227,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&waypoint.T{
					Id:          "REDSQU",
					Latitude:    55.754167,
					Longitude:   37.62,
					Altitude:    123,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&waypoint.T{
					Id:          "SYDNEY",
					Latitude:    -33.85695,
					Longitude:   151.215267,
					Altitude:    5,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
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
