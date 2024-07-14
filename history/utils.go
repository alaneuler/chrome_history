package history

import "time"

func ConvertChromeTime(msec int64) time.Time {
	sec := msec / 1000000
	nanoSec := (msec % 1000000) * 1000

	// Chrome's epoch starts from year 1601
	// https://stackoverflow.com/questions/20458406/what-is-the-format-of-chromes-timestamps
	t := time.Unix(sec, nanoSec)
	t = t.AddDate(1970-1601, 0, 0)
	return t
}
