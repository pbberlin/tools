package util

import (
	"fmt"
	"strconv"
	"time"
)

// TimeMarker gives us a simple string time stamp
func TimeMarkerPretty() string {
	tn := time.Now()
	//tn  = tn.Add( - time.Hour * 85 *24 )
	f := "2006-01-02 15:04:05"
	s := tn.Format(f)
	return s
}

// TimeMarker gives us a simple string time stamp
func TimeMarker() string {
	tn := time.Now()
	f := "2006-01-02-15-04-05"
	s := tn.Format(f)
	return s
}

func monthsBackComplicated() string {

	f := "2006-01-02 15:04:05 (MST)"
	s := ""

	tn := time.Now()
	y, m, d := tn.Date()

	loc, err := time.LoadLocation("Local")
	if err != nil {
		panic(fmt.Sprintln("could not load local time zone", err))
		loc = time.UTC
	}
	monthBegin := time.Date(y, m, d, 0, 0, 0, 0, loc)

	monthBegin = monthBegin.AddDate(0, -1, 0)

	s = monthBegin.Format(f)

	return s

}

// MonthsBack takes the current month
//  and returns a string "Year-Month" previous
func MonthsBack(monthsEarlier int) string {

	tn := time.Now()
	tn = tn.AddDate(0, -monthsEarlier, 0)
	s := tn.Format("2006-01")

	return s

}

func TimeFromUnix(uts string) time.Time {
	iuts, err := strconv.ParseInt(uts, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(iuts, 0)
	return tm
}
