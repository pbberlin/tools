package gitkit

import (
	"net/http"
	"time"

	aelog "google.golang.org/appengine/log"

	"google.golang.org/appengine"

	aeOrig "appengine"

	"appengine/datastore"
)

var (
	weekdays = []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}
)

type FavWeekday struct {
	// User ID. Serves as primary key in datastore.
	ID string
	// 0 is Sunday.
	Weekday time.Weekday
}

// weekdayForUser fetches the favorite weekday for the user from the datastore.
// Sunday is returned if no such data is found.
func weekdayForUser(r *http.Request, u *User) time.Weekday {
	c := aeOrig.NewContext(r)
	c2 := appengine.NewContext(r)
	k := datastore.NewKey(c, "FavWeekday", u.ID, 0, nil)
	d := FavWeekday{}
	err := datastore.Get(c, k, &d)
	if err != nil {
		if err != datastore.ErrNoSuchEntity {
			aelog.Errorf(c2, "Failed to fetch the favorite weekday for user %+v: %s", *u, err)
		}
		return time.Sunday
	}
	return d.Weekday
}

// updateWeekdayForUser updates the favorite weekday for the user.
func updateWeekdayForUser(r *http.Request, u *User, d time.Weekday) {
	c := aeOrig.NewContext(r)
	c2 := appengine.NewContext(r)
	k := datastore.NewKey(c, "FavWeekday", u.ID, 0, nil)
	_, err := datastore.Put(c, k, &FavWeekday{u.ID, d})
	if err != nil {
		aelog.Errorf(c2, "Failed to update the favorite weekday for user %+v: %s", *u, err)
	}
}
