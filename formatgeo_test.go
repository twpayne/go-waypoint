package waypoint

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatGeoReadWrite(t *testing.T) {
	for i, tc := range []struct {
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
					ID:          "Aconcagu",
					Latitude:    -D(32, 39, 12),
					Longitude:   -D(70, 0, 42),
					Altitude:    6962,
					Description: "Aconcagua",
				},
				&T{
					ID:          "Bergneus",
					Latitude:    D(51, 3, 7.02),
					Longitude:   D(7, 42, 22.02),
					Altitude:    488,
					Description: "Bergneustadt [A]",
				},
				&T{
					ID:          "Golden G",
					Latitude:    D(37, 49, 3),
					Longitude:   -D(122, 28, 42),
					Altitude:    227,
					Description: "Golden Gate Bridge",
				},
				&T{
					ID:          "Red Squa",
					Latitude:    D(55, 45, 15),
					Longitude:   D(37, 37, 12),
					Altitude:    123,
					Description: "Red Square",
				},
				&T{
					ID:          "Sydney O",
					Latitude:    -D(33, 51, 25.02),
					Longitude:   D(151, 12, 54.96),
					Altitude:    5,
					Description: "Sydney Opera",
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			wc, err := NewFormatGeoFormat().Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			assert.Equal(t, tc.wc, wc)

			w := &bytes.Buffer{}
			require.NoError(t, NewFormatGeoFormat().Write(w, tc.wc))
			assert.Equal(t, tc.s, w.String())

			_, f, err := Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			require.IsType(t, &FormatGeoFormat{}, f)
		})
	}
}
