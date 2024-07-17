package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelSpeed(t *testing.T) {
	cacheDir = filepath.Join(os.TempDir(), "chrome_history_cache/TestParallelSpeed")

	cleanDir(cacheDir)
	start := time.Now()
	Query("", 30, false)
	sequentialElapsed := time.Since(start)

	cleanDir(cacheDir)
	start = time.Now()
	Query("", 30, true)
	parallelElapsed := time.Since(start)

	assert.Greater(t, sequentialElapsed, parallelElapsed)
}

func cleanDir(dir string) {
	_ = os.RemoveAll(dir)
	_ = os.MkdirAll(dir, os.ModePerm)
}

func TestSorting(t *testing.T) {
	entries := Query("", 10, true)
	assert.Greater(t, len(entries), 2)
	for i := 1; i < len(entries); i++ {
		assert.True(t,
			entries[i-1].LastVisitTime.After(entries[i].LastVisitTime) ||
				entries[i-1].LastVisitTime.Equal(entries[i].LastVisitTime))
	}
}
