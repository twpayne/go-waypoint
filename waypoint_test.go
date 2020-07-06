package waypoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestData(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join("testdata", "*"))
	require.NoError(t, err)
	for _, match := range matches {
		t.Run(match, func(t *testing.T) {
			f, err := os.Open(match)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, f.Close())
			}()
			_, _, err = Read(f)
			require.NoError(t, err)
		})
	}
}
