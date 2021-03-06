/*
// http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go

// do something like this:
type Entity interface {
    UnmarshalHTTP(*http.Request) error
}

func GetEntity(r *http.Request, v Entity) error {
    return v.UnmarshalHTTP(r)
}

// Where the GetEntity function takes an interface value
// that is guaranteed to have an UnmarshalHTTP method.
// To make use of this, we would define on our User object
// some method that allows the User to describe
// how it would get itself out of an HTTP request
func (u *User) UnmarshalHTTP(r *http.Request) error {
   // ...
}

// in your application code,
// you would declare a var of User type,
// and then pass a pointer to this function into GetEntity
var u User
if err := GetEntity(req, &u); err != nil {
    // ...
}
*/
package extract
