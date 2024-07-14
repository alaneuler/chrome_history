package history

import "time"

type Entry struct {
	ID            int64
	URL           string
	Title         string
	VisitCount    int64
	LastVisitTime time.Time
	hidden        int64
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

func ToEntries(daoList []*EntryDao) []*Entry {
	var entries []*Entry
	for _, dao := range daoList {
		entries = append(entries, toEntry(dao))
	}
	return entries
}
