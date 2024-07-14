package history

import (
	"time"

	aw "github.com/deanishe/awgo"
)

type Entry struct {
	ID            int64
	URL           string
	Title         string
	VisitCount    int64
	LastVisitTime time.Time
	hidden        int64
	Icon          *aw.Icon
}

func toEntry(dao *EntryDao) *Entry {
	return &Entry{
		ID:            dao.ID,
		URL:           dao.URL,
		Title:         dao.Title,
		VisitCount:    dao.VisitCount,
		LastVisitTime: ConvertChromeTime(dao.LastVisitTime),
		hidden:        dao.hidden,
	}
}
