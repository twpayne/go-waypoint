package formatgeo

import (
	"bitbucket.org/twpayne/waypoint"
	"bitbucket.org/twpayne/waypoint/internal/dmsh"

	"reflect"
	"testing"
)

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		ls   []string
		want bool
	}{
		{ls: []string{}, want: false},
		{ls: []string{"$FormatGEO"}, want: true},
		{ls: []string{"$FormatGEO "}, want: true},
		{ls: []string{"$FormatUTM"}, want: false},
	} {
		if got := New().Detect(c.ls); got != c.want {
			t.Errorf("Detect(%v) == %v, want %v", c.ls, got, c.want)
		}
	}
}

func TestReadWrite(t *testing.T) {
	for _, c := range []struct {
		ls []string
		wc waypoint.Collection
	}{
		{
			ls: []string{
				"$FormatGEO",
				"Aconcagu  S 32 39 12.00    W 070 00 42.00  6962  Aconcagua",
				"Bergneus  N 51 03 07.02    E 007 42 22.02   488  Bergneustadt [A]",
				"Golden G  N 37 49 03.00    W 122 28 42.00   227  Golden Gate Bridge",
				"Red Squa  N 55 45 15.00    E 037 37 12.00   123  Red Square",
				"Sydney O  S 33 51 25.02    E 151 12 54.96     5  Sydney Opera",
			},
			wc: waypoint.Collection{
				&waypoint.T{
					Id:          "Aconcagu",
					Latitude:    -dmsh.D(32, 39, 12),
					Longitude:   -dmsh.D(70, 0, 42),
					Altitude:    6962,
					Description: "Aconcagua",
				},
				&waypoint.T{
					Id:          "Bergneus",
					Latitude:    dmsh.D(51, 3, 7.02),
					Longitude:   dmsh.D(7, 42, 22.02),
					Altitude:    488,
					Description: "Bergneustadt [A]",
				},
				&waypoint.T{
					Id:          "Golden G",
					Latitude:    dmsh.D(37, 49, 3),
					Longitude:   -dmsh.D(122, 28, 42),
					Altitude:    227,
					Description: "Golden Gate Bridge",
				},
				&waypoint.T{
					Id:          "Red Squa",
					Latitude:    dmsh.D(55, 45, 15),
					Longitude:   dmsh.D(37, 37, 12),
					Altitude:    123,
					Description: "Red Square",
				},
				&waypoint.T{
					Id:          "Sydney O",
					Latitude:    -dmsh.D(33, 51, 25.02),
					Longitude:   dmsh.D(151, 12, 54.96),
					Altitude:    5,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		if got, err := New().Read(c.ls); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if !reflect.DeepEqual(w, got[i]) {
					t.Errorf("i=%d w=%v got[%d]=%v", i, w, got[i])
				}
			}
			t.Errorf("Read(%v) == %v, %v, want %v, nil", c.ls, got, err, c.wc)
		}
	}
}
