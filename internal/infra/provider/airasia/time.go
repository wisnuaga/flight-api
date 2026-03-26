package airasia

import (
	"time"

	"github.com/wisnuaga/flight-api/internal/util"
)

func parseTime(s string) (time.Time, error) {
	return util.ParseTime(s)
}
