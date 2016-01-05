package waypoint

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestFormatGeoReadWrite(t *testing.T) {
	for _, c := range []struct {
		s  string
		wc Collection
	}{
		{
			s: "$FormatGEO\r\n" +
				"Aconcagu  S 32 39 12.00    W 070 00 42.00  6962  Aconcagua\r\n" +
				"Bergneus  N 51 03 07.02    E 007 42 22.02   488  Bergneustadt [A]\r\n" +
				"Golden G  N 37 49 03.00    W 122 28 42.00   227  Golden Gate Bridge\r\n" +
				"Red Squa  N 55 45 15.00    E 037 37 12.00   123  Red Square\r\n" +
				"Sydney O  S 33 51 25.02    E 151 12 54.96     5  Sydney Opera\r\n",
			wc: Collection{
				&T{
					Id:          "Aconcagu",
					Latitude:    -D(32, 39, 12),
					Longitude:   -D(70, 0, 42),
					Altitude:    6962,
					Description: "Aconcagua",
				},
				&T{
					Id:          "Bergneus",
					Latitude:    D(51, 3, 7.02),
					Longitude:   D(7, 42, 22.02),
					Altitude:    488,
					Description: "Bergneustadt [A]",
				},
				&T{
					Id:          "Golden G",
					Latitude:    D(37, 49, 3),
					Longitude:   -D(122, 28, 42),
					Altitude:    227,
					Description: "Golden Gate Bridge",
				},
				&T{
					Id:          "Red Squa",
					Latitude:    D(55, 45, 15),
					Longitude:   D(37, 37, 12),
					Altitude:    123,
					Description: "Red Square",
				},
				&T{
					Id:          "Sydney O",
					Latitude:    -D(33, 51, 25.02),
					Longitude:   D(151, 12, 54.96),
					Altitude:    5,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		if got, err := NewFormatGeoFormat().Read(strings.NewReader(c.s)); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if err := equal(w, got[i]); err != nil {
					t.Errorf("want %#v got=%#v, %v", w, got[i], err)
				}
			}
			t.Errorf("Read(strings.NewReader(%v)) == %v, %v, want %v, nil", c.s, got, err, c.wc)
		}
		w := bytes.NewBuffer(nil)
		if err := NewFormatGeoFormat().Write(w, c.wc); err != nil {
			t.Errorf("Write(%v) == %v. want nil", c.wc, err)
		}
		if w.String() != c.s {
			t.Errorf("w.String() == %v. want %v", w.String(), c.s)
		}
		_, f, err := Read(strings.NewReader(c.s))
		if err != nil {
			t.Errorf("Read(...) return %v, expected nil", err)
		}
		if _, ok := f.(*FormatGeoFormat); !ok {
			t.Errorf("Read(...) returned a %T, expected a FormatGeoFormat", f)
		}
	}
}
