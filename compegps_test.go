package waypoint

import (
	"bytes"
	"image/color"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompeGPSReadWrite(t *testing.T) {
	for i, tc := range []struct {
		s  string
		wc Collection
	}{
		{
			s: "G  WGS 84\r\n" +
				"U  1\r\n" +
				"W  ACONCA A 32.6533333333\xbdS 70.0116666667\xbdW 27-MAR-62 00:00:00 6962.000000 Highest mountain in south-america\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  BERGNE A 51.0519500000\xbdN 7.70611700000\xbdE 27-MAR-62 00:00:00 488.000000 Rabbit holes, 20\" ditch south end of rwy\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  GOLDEN A 37.8175000000\xbdN 122.478333333\xbdW 27-MAR-62 00:00:00 227.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  REDSQU A 55.7541670000\xbdN 37.6200000000\xbdE 27-MAR-62 00:00:00 123.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n" +
				"W  SYDNEY A 33.8569500000\xbdS 151.215267000\xbdE 27-MAR-62 00:00:00 5.000000\r\n" +
				"w  box,0,0.0,16777215,255,1,7,,0.0\r\n",
			wc: Collection{
				&T{
					ID:          "ACONCA",
					Latitude:    -32.6533333333,
					Longitude:   -70.0116666667,
					Altitude:    6962,
					Description: "Highest mountain in south-america",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&T{
					ID:          "BERGNE",
					Latitude:    51.05195,
					Longitude:   7.706117,
					Altitude:    488,
					Description: "Rabbit holes, 20\" ditch south end of rwy",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&T{
					ID:          "GOLDEN",
					Latitude:    37.8175,
					Longitude:   -122.478333333,
					Altitude:    227,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&T{
					ID:          "REDSQU",
					Latitude:    55.754167,
					Longitude:   37.62,
					Altitude:    123,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
				&T{
					ID:          "SYDNEY",
					Latitude:    -33.85695,
					Longitude:   151.215267,
					Altitude:    5,
					Description: "",
					Color:       color.RGBA{R: 255, G: 255, B: 255},
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			wc, err := NewCompeGPSFormat().Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			assert.Equal(t, tc.wc, wc)

			w := &bytes.Buffer{}
			require.NoError(t, NewCompeGPSFormat().Write(w, tc.wc))

			_, f, err := Read(strings.NewReader(tc.s))
			require.NoError(t, err)
			require.IsType(t, &CompeGPSFormat{}, f)

			// FIXME
			// if w.String() != c.s {
			//	checkStrings(t, w.String(), c.s)
			//	t.Errorf("w.String() == %v. want %v", w.String(), c.s)
			// }
		})
	}
}
