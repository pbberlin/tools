package weedout

import "golang.org/x/net/html"

// !DOCTYPE html head
// !DOCTYPE html body
//        0    1    2
const cScaffoldLvls = 2

const numTotal = 3 // comparable html docs, including base doc
const stageMax = 3 // weedstages

var URLs = []string{
	"www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html",
	"www.economist.com/news/britain/21663648-hard-times-hard-hats-making-britain-make-things-again-proving-difficult",
	"www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power",
}

func attrX(attributes []html.Attribute, key string) (s string) {
	for _, a := range attributes {
		if key == a.Key {
			s = a.Val
			break
		}
	}
	return
}
