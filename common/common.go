package common

import "time"

// TimeUnitMilis miliseconds in a second
const TimeUnitMilis = 1e6 // 1 milisecond

// TimeInMilis converts time to miiseconds
func TimeInMilis(t time.Time) int64 {
	return t.UTC().UnixNano() / TimeUnitMilis
}

// ElapsedTimeInMilis time elapsed from star time in milis
func ElapsedTimeInMilis(startTime int64) int64 {
	return TimeInMilis(time.Now()) - startTime
}
