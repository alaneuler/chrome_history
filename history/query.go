package history

import (
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

func Query(query string, limit int) []*Entry {
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
	return toEntries(entries)
}

func toEntries(daoList []*EntryDao) []*Entry {
	db, err := openFaviconDb()
	if err != nil {
		slog.Error("Open Favicon database", "error", err)
	}

	var entries []*Entry
	m := make(map[string]int)
	for _, dao := range daoList {
		if strings.HasPrefix(dao.URL, "http://") {
			continue
		}

		if _, ok := m[dao.URL]; !ok {
			m[dao.URL] = 1

			entry := toEntry(dao)
			entry.Icon = ObtainIcon(db, dao)
			entries = append(entries, entry)
		}
	}
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
