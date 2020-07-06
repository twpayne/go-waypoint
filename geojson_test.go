package waypoint

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeoJSONReadWrite(t *testing.T) {
	for i, tc := range []struct {
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
					ID:          "Aconcagu",
					Latitude:    -32.653333,
					Longitude:   -70.011667,
					Altitude:    6962,
					Description: "Aconcagua",
				},
				&T{
					ID:          "Bergneus",
					Latitude:    51.05195,
					Longitude:   7.706117,
					Altitude:    488,
					Description: "Bergneustadt [A]",
				},
				&T{
					ID:          "Golden G",
					Latitude:    37.8175,
					Longitude:   -122.478333,
					Altitude:    227,
					Description: "Golden Gate Bridge",
				},
				&T{
					ID:          "Red Squa",
					Latitude:    55.754167,
					Longitude:   37.62,
					Altitude:    123,
					Description: "Red Square",
				},
				&T{
					ID:          "Sydney O",
					Latitude:    -33.85695,
					Longitude:   151.215267,
					Altitude:    5,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			wc, err := NewGeoJSONFormat().Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			assert.Equal(t, tc.wc, wc)

			w := &bytes.Buffer{}
			require.NoError(t, NewGeoJSONFormat().Write(w, tc.wc))
			assert.Equal(t, tc.s, w.String())

			_, f, err := Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			require.IsType(t, &GeoJSONFormat{}, f)
		})
	}
}
