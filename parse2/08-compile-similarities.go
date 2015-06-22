package parse2

import (
	"fmt"
	"os"
)

func similaritiesToFile(frags []fragment, stage int) {

	// bfrags := pbstrings.IndentedDumpBytes(frags)
	bfrags := []byte{}
	for _, v := range frags {
		bfrags = append(bfrags, v.ArticleUrl...)
		bfrags = append(bfrags, ' ')
		bfrags = append(bfrags, fmt.Sprintf("%v", v.Lvl)...)
		bfrags = append(bfrags, ' ')
		bfrags = append(bfrags, fmt.Sprintf("%-8v", string(v.Outline))...)
		bfrags = append(bfrags, "             "...)
		bfrags = append(bfrags, string(v.Text)...)
		bfrags = append(bfrags, '\n')
		for _, v1 := range v.Similars {
			bfrags = append(bfrags, v1.ArticleUrl...)
			bfrags = append(bfrags, ' ')
			bfrags = append(bfrags, fmt.Sprintf("%v", v1.Lvl)...)
			bfrags = append(bfrags, ' ')
			bfrags = append(bfrags, fmt.Sprintf("%-8v", string(v1.Outline))...)
			bfrags = append(bfrags, "    "...)
			bfrags = append(bfrags, spf("%2v ", v1.AbsLevenshtein)...)
			bfrags = append(bfrags, spf("%-5.2v ", v1.RelLevenshtein)...)
			bfrags = append(bfrags, string(v1.Text)...)
			bfrags = append(bfrags, '\n')
		}
		bfrags = append(bfrags, '\n')
	}
	bytes2File(spf("outp_frags_st%v.txt", stage), bfrags)

}

func assembleWeedout(frags []fragment, ret map[string]map[string]bool) map[string]map[string]bool {

	for _, v := range frags {

		// if len(v.Similars) >= numTotal-1 {
		if len(v.Similars) >= 2 {

			lvlHighest := v.Lvl
			for _, v1 := range v.Similars {
				if v1.Lvl < lvlHighest {
					lvlHighest = v1.Lvl
				}
			}

			for _, v1 := range v.Similars {
				if v1.Lvl == lvlHighest {
					if ret[v1.ArticleUrl] == nil {
						pf("WANT %v\n", v1.ArticleUrl)
						for k, _ := range ret {
							pf("     %v\n", k)
						}
						os.Exit(1)
					}
					ret[v1.ArticleUrl][v1.Outline] = true
				}
			}
			if v.Lvl == lvlHighest {
				ret[v.ArticleUrl][v.Outline] = true
			}

		}

	}

	// pf("%v\n", ret)
	return ret
}
