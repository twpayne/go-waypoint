package waypoint

import (
	"testing"
)

func Test_dmsh(t *testing.T) {
	for _, c := range []struct {
		x  float64
		hs Hemisphere
		d  int
		m  int
		s  float64
		h  uint8
	}{
		{x: 0, hs: NS, d: 0, m: 0, s: 0, h: 'N'},
		{x: 1, hs: NS, d: 1, m: 0, s: 0, h: 'N'},
		{x: -1, hs: NS, d: 1, m: 0, s: 0, h: 'S'},
		{x: 1.5, hs: NS, d: 1, m: 30, s: 0, h: 'N'},
		{x: 1.75, hs: NS, d: 1, m: 45, s: 0, h: 'N'},
		{x: float64(1) / 3600, hs: NS, d: 0, m: 0, s: 1, h: 'N'},
	} {
		if d, m, s, h := dmsh(c.x, c.hs); d != c.d || m != c.m || s != c.s || h != c.h {
			t.Errorf("dmsh(%v, %v) == %v, %v, %v, %v, want %v, %v, %v, %v", c.x, c.hs, d, m, s, h, c.d, c.m, c.s, c.h)
		}
	}
}
