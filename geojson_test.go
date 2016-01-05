package waypoint

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestGeoJSONReadWrite(t *testing.T) {
	for _, c := range []struct {
		s  string
		wc Collection
	}{
		{
			s: `{"features":[` +
				`{"geometry":{"coordinates":[-32.653333,-70.011667,6962],"type":"Point"},"id":"Aconcagu","properties":{"description":"Aconcagua"},"type":"Feature"},` +
				`{"geometry":{"coordinates":[51.05195,7.706117,488],"type":"Point"},"id":"Bergneus","properties":{"description":"Bergneustadt [A]"},"type":"Feature"},` +
				`{"geometry":{"coordinates":[37.8175,-122.478333,227],"type":"Point"},"id":"Golden G","properties":{"description":"Golden Gate Bridge"},"type":"Feature"},` +
				`{"geometry":{"coordinates":[55.754167,37.62,123],"type":"Point"},"id":"Red Squa","properties":{"description":"Red Square"},"type":"Feature"},` +
				`{"geometry":{"coordinates":[-33.85695,151.215267,5],"type":"Point"},"id":"Sydney O","properties":{"description":"Sydney Opera"},"type":"Feature"}` +
				`],"type":"FeatureCollection"}` + "\n",
			wc: Collection{
				&T{
					Id:          "Aconcagu",
					Latitude:    -32.653333,
					Longitude:   -70.011667,
					Altitude:    6962,
					Description: "Aconcagua",
				},
				&T{
					Id:          "Bergneus",
					Latitude:    51.05195,
					Longitude:   7.706117,
					Altitude:    488,
					Description: "Bergneustadt [A]",
				},
				&T{
					Id:          "Golden G",
					Latitude:    37.8175,
					Longitude:   -122.478333,
					Altitude:    227,
					Description: "Golden Gate Bridge",
				},
				&T{
					Id:          "Red Squa",
					Latitude:    55.754167,
					Longitude:   37.62,
					Altitude:    123,
					Description: "Red Square",
				},
				&T{
					Id:          "Sydney O",
					Latitude:    -33.85695,
					Longitude:   151.215267,
					Altitude:    5,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		if got, err := NewGeoJSONFormat().Read(strings.NewReader(c.s)); err != nil || !reflect.DeepEqual(got, c.wc) {
			for i, w := range c.wc {
				if err := Equal(w, got[i]); err != nil {
					t.Errorf("want %#v got=%#v, %v", w, got[i], err)
				}
			}
			t.Errorf("Read(strings.NewReader(%v)) == %v, %v, want %v, nil", c.s, got, err, c.wc)
		}
		w := bytes.NewBuffer(nil)
		if err := NewGeoJSONFormat().Write(w, c.wc); err != nil {
			t.Errorf("Write(%v) == %v. want nil", c.wc, err)
		}
		if w.String() != c.s {
			t.Errorf("w.String(),\n got %v\nwant %v", w.String(), c.s)
		}
		_, f, err := Read(strings.NewReader(c.s))
		if err != nil {
			t.Errorf("Read(...) return %v, expected nil", err)
		}
		if _, ok := f.(*GeoJSONFormat); !ok {
			t.Errorf("Read(...) returned a %T, expected a GeoJSONFormat", f)
		}
	}
}
