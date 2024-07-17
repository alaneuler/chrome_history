package history

import (
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	profileKey     = "CHROME_PROFILE"
	defaultProfile = "Default"

	chromePath = "Library/Application Support/Google/Chrome"

	historyTableName = "urls"
)

var (
	historyDatabasePath string
	faviconDatabasePath string
)

func init() {
	profile := os.Getenv(profileKey)
	if profile == "" {
		profile = defaultProfile
	}

	home, _ := os.UserHomeDir()
	historyDatabasePath = filepath.Join(home, chromePath, profile, "History")
	faviconDatabasePath = filepath.Join(home, chromePath, profile, "Favicons")
}

func Query(query string, limit int, parallel bool) []*Entry {
	db, err := openHistoryDb()
	if err != nil {
		slog.Error("Open Chrome history database", "error", err)
		return nil
	}

	db = db.Where("visit_count > 0").Where("hidden = 0")
	db = db.Order("last_visit_time desc")
	if query != "" {
		titleOrUrl := "%" + query + "%"
		slog.Info("Starting query:", "titleOrUrl", titleOrUrl)
		db = db.Where("title like ? or url like ?", titleOrUrl, titleOrUrl)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}

	var entries []*EntryDao
	db.Find(&entries)
	entries = filter(entries)
	if parallel {
		return toEntriesParallel(entries)
	}
	return toEntries(entries)
}

func filter(daoList []*EntryDao) []*EntryDao {
	var rtn []*EntryDao
	for _, dao := range daoList {
		if strings.HasPrefix(dao.URL, "http://") {
			continue
		}

		rtn = append(rtn, dao)
	}
	return rtn
}

func toEntries(daoList []*EntryDao) []*Entry {
	slog.Info("toEntries sequentially", "len", len(daoList))
	start := time.Now()
	db, err := openFaviconDb()
	if err != nil {
		slog.Error("Open Favicon database", "error", err)
	}

	var entries []*Entry
	for _, dao := range daoList {
		entry := toEntry(dao)
		entry.Icon = ObtainIcon(db, dao)
		entries = append(entries, entry)
	}
	slog.Info("toEntries", "elapsed", time.Since(start))
	return entries
}

func toEntriesParallel(daoList []*EntryDao) []*Entry {
	slog.Info("toEntriesParallel", "len", len(daoList))
	start := time.Now()
	db, err := openFaviconDb()
	if err != nil {
		slog.Error("Open Favicon database", "error", err)
	}

	var entries []*Entry
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	for _, dao := range daoList {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			entry := toEntry(dao)
			entry.Icon = ObtainIcon(db, dao)
			entries = append(entries, entry)
		}()
	}
	wg.Wait()
	slog.Info("toEntriesParallel", "elapsed", time.Since(start))
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].LastVisitTime.After(entries[j].LastVisitTime)
	})
	return entries
}

func encodePath(path string) string {
	dsn, _ := url.Parse(path)
	dsn.Scheme = "file"
	q := dsn.Query()
	q.Set("mode", "ro")
	q.Set("immutable", "1")
	q.Set("_query_only", "1")
	dsn.RawQuery = q.Encode()
	return dsn.String()
}

func openHistoryDb() (*gorm.DB, error) {
	path := encodePath(historyDatabasePath)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.Table(historyTableName), nil
}

func openFaviconDb() (*gorm.DB, error) {
	path := encodePath(faviconDatabasePath)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
