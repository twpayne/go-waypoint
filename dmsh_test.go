package waypoint

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func Test_dmsh(t *testing.T) {
	for i, tc := range []struct {
		x  float64
		hs hemisphere
		d  int
		m  int
		s  float64
		h  uint8
	}{
		{x: 0, hs: ns, d: 0, m: 0, s: 0, h: 'N'},
		{x: 1, hs: ns, d: 1, m: 0, s: 0, h: 'N'},
		{x: -1, hs: ns, d: 1, m: 0, s: 0, h: 'S'},
		{x: 1.5, hs: ns, d: 1, m: 30, s: 0, h: 'N'},
		{x: 1.75, hs: ns, d: 1, m: 45, s: 0, h: 'N'},
		{x: float64(1) / 3600, hs: ns, d: 0, m: 0, s: 1, h: 'N'},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			d, m, s, h := dmsh(tc.x, tc.hs)
			assert.Equal(t, tc.d, d)
			assert.Equal(t, tc.m, m)
			assert.Equal(t, tc.s, s)
			assert.Equal(t, tc.h, h)
		})
	}
}
