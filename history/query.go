package history

import (
	"log/slog"
	"net/url"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbLoc     = "/Users/alaneuler/Library/Application Support/Google/Chrome/Default/History"
	tableName = "urls"
)

func Query() []*Entry {
	db, err := open()
	if err != nil {
		slog.Error("Open Chrome history database error", err)
		return nil
	}

	var entries []*EntryDao
	db.Where("visit_count > 0").Where("hidden = 0").Order("last_visit_time desc").Find(&entries)

	return ToEntries(entries)
}

func open() (*gorm.DB, error) {
	dsn, err := url.Parse(dbLoc)
	if err != nil {
		return nil, err
	}

	dsn.Scheme = "file"
	db, err := gorm.Open(sqlite.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.Table(tableName), nil
}
