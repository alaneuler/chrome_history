package history

type EntryDao struct {
	ID            int64
	URL           string
	Title         string
	VisitCount    int64
	LastVisitTime int64
	hidden        int64
}
