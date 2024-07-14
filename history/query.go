package history

import (
	"log/slog"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	profileKey     = "CHROME_PROFILE"
	defaultProfile = "Default"
	chromeLoc      = "~/Library/Application Support/Google/Chrome"
	dbFile         = "History"
	tableName      = "urls"
)

func Query(query string, limit int) []*Entry {
	db, err := open()
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
	return ToEntries(entries)
}

func obtainHistoryDbFile() (string, error) {
	profile := os.Getenv(profileKey)
	if profile == "" {
		profile = defaultProfile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	loc := filepath.Join(home, chromeLoc[2:])
	uri, err := url.Parse(loc)
	if err != nil {
		return "", err
	}

	uri.Path = path.Join(uri.Path, profile, dbFile)
	uri.Scheme = "file"
	return uri.String(), nil
}

func open() (*gorm.DB, error) {
	loc, err := obtainHistoryDbFile()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(loc), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.Table(tableName), nil
}
