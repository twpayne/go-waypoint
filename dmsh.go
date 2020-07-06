package waypoint

import (
	"math"
)

type hemisphere []uint8

var (
	ew hemisphere = []uint8{'E', 'W'}
	ns hemisphere = []uint8{'N', 'S'}
)

// D converts degrees, minutes, and seconds into decimal degrees.
func D(d, m, s float64) float64 {
	return d + m/60 + s/3600
}

// dmh converts degrees into degrees, decimal minutes, and a hemisphere. hs
// should be "NS" for latitude and "EW" for latitude.
func dmh(x float64, hs hemisphere) (d int, m float64, h uint8) {
	if x < 0 {
		h = hs[1]
		x = -x
	} else {
		h = hs[0]
	}
	d = int(x)
	m = math.Mod(60*x, 60)
	return
}

// dmsh converts degrees to degrees, minutes, decimal seconds, and a
// hemisphere. hs should be "NS" for latitude and "EW" for longitude.
func dmsh(x float64, hs hemisphere) (d, m int, s float64, h uint8) {
	if x < 0 {
		h = hs[1]
		x = -x
	} else {
		h = hs[0]
	}
	d = int(x)
	m = int(60*x) % 60
	s = math.Mod(3600*x, 60)
	return
}
