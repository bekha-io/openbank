package fineract

import "time"

type TimeSlice []int

func (s TimeSlice) Time() time.Time {
	return time.Date(s[0], time.Month(s[1]), s[2], 0, 0, 0, 0, time.UTC)
}