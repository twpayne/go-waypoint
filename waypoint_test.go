package waypoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestTestData(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join("testdata", "*"))
	assert.NoError(t, err)
	for _, match := range matches {
		t.Run(match, func(t *testing.T) {
			f, err := os.Open(match)
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, f.Close())
			}()
			collection, format, err := Read(f)
			assert.NoError(t, err)
			assert.NotZero(t, len(collection))
			assert.NotZero(t, format)
		})
	}
}

func assertIsType[T any](t testing.TB, expected T, actual any) {
	_, ok := actual.(T)
	if ok {
		return
	}
	t.Helper()
	t.Fatalf("Expected %v to be of type %T", actual, expected)
}
