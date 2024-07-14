package history

import (
	"log/slog"
	"os"
	"path/filepath"

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
	historyDatabasePath = "file://" + filepath.Join(home, chromePath, profile, "History")
	faviconDatabasePath = "file://" + filepath.Join(home, chromePath, profile, "Favicons")
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
	for _, dao := range daoList {
		entry := toEntry(dao)
		entry.Icon = ObtainIcon(db, dao)
		entries = append(entries, entry)
	}
	return entries
}

func openHistoryDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(historyDatabasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.Table(historyTableName), nil
}

func openFaviconDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(faviconDatabasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
