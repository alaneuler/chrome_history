package history

type EntryDao struct {
	ID            int64
	URL           string
	Title         string
	VisitCount    int64
	LastVisitTime int64
	hidden        int64
}

type IconMappingDao struct {
	ID      int64
	PageURL string
	IconId  int64
}

type IconBitmapDao struct {
	ID        int64
	IconId    int64
	ImageData []byte
}
