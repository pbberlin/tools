package weedout

import "golang.org/x/net/html"

// !DOCTYPE html head
// !DOCTYPE html body
//        0    1    2
const cScaffoldLvls = 2

const numTotal = 3 // comparable html docs, including base doc
const stageMax = 3 // weedstages

const cTestHostDev = "localhost:8085"

func attrX(attributes []html.Attribute, key string) (s string) {
	for _, a := range attributes {
		if key == a.Key {
			s = a.Val
			break
		}
	}
	return
}
