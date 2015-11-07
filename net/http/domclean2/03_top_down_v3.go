package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

// Now this third implementation finally condenses *selectively*.
// Not all boats from each pond are lifted equally.
// We achieve tremendous structural simplification.
// It also starts from top, pulling lower levels up.
// Unlike implementation #1, that started from the middle.
func topDownV3(l1 *html.Node, l2Types map[string]bool, l3Types map[string]bool) {

	if l1.Type != html.ElementNode &&
		l1.Type != html.DocumentNode {
		return // cannot assign to - do not unable to have children
	}
	if l1.Data == "span" || l1.Data == "a" {
		return // want not condense into
	}

	// dig two levels deep

	// isolate l2,l3
	l2s := []*html.Node{}
	l3s := map[*html.Node][]*html.Node{}

	for l2 := l1.FirstChild; l2 != nil; l2 = l2.NextSibling {

		l2s = append(l2s, l2)
		// l2s = append([]*html.Node{l2}, l2s...) // order inversion

		for l3 := l2.FirstChild; l3 != nil; l3 = l3.NextSibling {
			l3s[l2] = append(l3s[l2], l3)
			// l3s[l2] = append(map[*html.Node][]*html.Node{l2: []*html.Node{l3}}, l3s[l2]...) // order inversion
		}
	}

	postponedRemoval := map[*html.Node]bool{}

	//
	//
	// check types for each l2 subtree distinctively
	for _, l2 := range l2s {

		l2Match := l2.Type == html.ElementNode && l2Types[l2.Data] // l2 is a div

		l3Match := true
		for _, l3 := range l3s[l2] {
			l3Match = l3Match && (l3.Type == html.ElementNode && l3Types[l3.Data])
		}

		// act
		if l2Match && l3Match {

			// detach l3 from l2
			for _, l3 := range l3s[l2] {
				// if ml3[l3] > 0 {
				// 	fmt.Printf("rmd_%v_%v ", ml3[l3], l3.Data)
				// }
				l2.RemoveChild(l3)
				// ml3[l3]++
			}

			// Since we still need l2 below
			// We have to postpone detaching l2 from l1
			// to the bottom
			// NOT HERE: l1.RemoveChild(l2)
			postponedRemoval[l2] = true

			for _, l3 := range l3s[l2] {
				// attach l3 to l1

				if l3.Data != "a" && l3.Data != "span" {
					l1.InsertBefore(l3, l2)
				} else {
					wrap := dom.Nd("p")
					wrap.Attr = []html.Attribute{html.Attribute{Key: "cfrm", Val: "noth"}}
					wrap.AppendChild(l3)
					// NOT  wrap.FirstChild = l3
					l1.InsertBefore(wrap, l2)
				}
			}

		}

	}

	for k, _ := range postponedRemoval {
		l1.RemoveChild(k) // detach l2 from l1
	}

}
