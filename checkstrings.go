package waypoint

import (
	"testing"
)

func checkStrings(t *testing.T, s1, s2 string) {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}
	line := 1
	col := 0
	for i := 0; i < n; i++ {
		col++
		if s1[i] != s2[i] {
			t.Errorf("strings differ a line %d column %d (%v != %v)", line, col, s1[i], s2[i])
			break
		}
		if s1[i] == '\n' {
			line++
			col = 0
		}
	}
}
