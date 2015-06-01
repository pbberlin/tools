package extract

import (
	"fmt"
	"net/http"
)

type Extractable interface {
	ExtractFromHtmlPost(r *http.Request) error
}

func Extract(r *http.Request, e Extractable) error {
	return e.ExtractFromHtmlPost(r)
}

// ===============================================

type User struct {
	Name, Surname string
}

func (u *User) ExtractFromHtmlPost(r *http.Request) error {

	if r == nil {
		return fmt.Errorf("request is empty")
	}

	u.Name = "extracted from"
	u.Surname = "http request"
	return nil

}

func (u User) String() string {
	return fmt.Sprintf("SName: %q LName: %q", u.Surname, u.Name)
}
