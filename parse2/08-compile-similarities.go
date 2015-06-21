package parse2

import "fmt"

func compileSimarities() {

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
	bytes2File("outp_frags.txt", bfrags)

}

func weedOut() {

	for _, v := range frags {

		if len(v.Similars) >= numTotal-1 {

			lvlHighest := v.Lvl
			for _, v1 := range v.Similars {
				if v1.Lvl < lvlHighest {
					lvlHighest = v1.Lvl
				}
			}

		}
	}
}
