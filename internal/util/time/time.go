package time

import (
	"fmt"
	"math"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int64(math.Floor(d.Hours()))
	tmp := d - time.Duration(hours)*time.Hour
	minutes := int64(math.Floor(tmp.Minutes()))

	return fmt.Sprintf("%02dh%02dm", hours, minutes)
}
