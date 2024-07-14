package utils

import (
	"regexp"
	"strconv"
	"time"
)

// ParseDuration parses an ISO 8601 string representing a duration,
// and returns the resultant golang time.Duration instance.
func ParseDuration(isoDuration string) time.Duration {
	re := regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:.\d+)?)S)?$`)
	matches := re.FindStringSubmatch(isoDuration)
	if matches == nil {
		return 0
	}

	seconds := 0.0

	// Day.
	if matches[3] != "" {
		f, err := strconv.ParseFloat(matches[3], 32)
		if err != nil {
			return 0
		}
		seconds += (f * 24 * 60 * 60)
	}

	// Hour.
	if matches[4] != "" {
		f, err := strconv.ParseFloat(matches[4], 32)
		if err != nil {
			return 0
		}
		seconds += (f * 60 * 60)
	}

	// Minute.
	if matches[5] != "" {
		f, err := strconv.ParseFloat(matches[5], 32)
		if err != nil {
			return 0
		}
		seconds += (f * 60)
	}

	// Second & millisecond.
	if matches[6] != "" {
		f, err := strconv.ParseFloat(matches[6], 32)
		if err != nil {
			return 0
		}
		seconds += f
	}

	d, _ := time.ParseDuration(strconv.FormatFloat(seconds, 'f', -1, 32) + "s")
	return d

}

// Thousands to format int to thousands string format.
func Thousands(num int) string {
	str := strconv.Itoa(num)
	lStr := len(str)
	digits := lStr
	if num < 0 {
		digits--
	}
	commas := (digits+2)/3 - 1
	lBuf := lStr + commas
	var sbuf [32]byte // pre allocate buffer at stack rather than make([]byte,n)
	buf := sbuf[0:lBuf]
	// copy str from the end
	for si, bi, c3 := lStr-1, lBuf-1, 0; ; {
		buf[bi] = str[si]
		if si == 0 {
			return string(buf)
		}
		si--
		bi--
		// insert comma every 3 chars
		c3++
		if c3 == 3 && (si > 0 || num > 0) {
			buf[bi] = ','
			bi--
			c3 = 0
		}
	}
}
