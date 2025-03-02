package utils

import "time"

type Time interface {
	Now() time.Time
	Pattern() string
}

type RealTime struct{}

func NewRealTime() *RealTime {
	return &RealTime{}
}

func (rt *RealTime) Now() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func (rt *RealTime) Pattern() string {
	return "2006-01-02"
}

func Pattern() string {
	return "2006-01-02"
}

func YMDKey(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}
