package output

import (
	"time"
)

func convertFromMsToTime(raw int64) time.Time {
	return time.Unix(0, raw*int64(time.Millisecond))
}
