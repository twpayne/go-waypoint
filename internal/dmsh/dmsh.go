package dmsh

import (
	"math"
)

func D(d, m, s float64) float64 {
	return d + m/60 + s/3600
}

func DMSH(x float64, hs string) (d, m int, s float64, h string) {
	if x < 0 {
		h = hs[1:2]
		x = -x
	} else {
		h = hs[:1]
	}
	d = int(x)
	m = int(60*x) % 60
	s = math.Mod(3600*x, 60)
	return
}
