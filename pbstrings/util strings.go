package pbstrings

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var nonAscii = regexp.MustCompile(`[^a-zA-Z0-9\.\_]+`)
var mutatedVowels = strings.NewReplacer("ä", "ae", "ö", "oe", "ü", "ue", "Ä", "ae", "Ö", "oe", "Ü", "ue")

func StringNormalize(s string) string {
	return LowerCasedUnderscored(s)
}

// LowerCasedUnderscored gives us a condensed filename
// cleansed of all non Ascii characters
// where word boundaries are encoded by "_"
//
// whenever we want a transformation of user input
// into innoccuous lower case - sortable - searchable
// ascii - the we should look to this func

// in addition - extensions are respected and cleansed
func LowerCasedUnderscored(s string) string {

	//log.Printf("%v\n", s)

	s = mutatedVowels.Replace(s)

	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"' `)

	replaced := nonAscii.ReplaceAllString(s, "_")

	replaced = strings.Trim(replaced, `_`)
	replaced = strings.ToLower(replaced)

	// clean the  file extension
	replaced = strings.Replace(replaced, "_.", ".", -1)
	replaced = strings.Replace(replaced, "._", ".", -1)

	//log.Printf("%v\n", replaced)

	return replaced
}

// ToLen chops or extends string to the exactly desired length
// format strings like %4v do not restrict.
func ToLen(s string, nx int) string {
	if len(s) < nx {
		return fmt.Sprintf("%v%v", s, strings.Repeat(" ", nx-len(s)))
	} else if len(s) > nx {
		return s[:nx]
	} else {
		return s
	}
}

//  followed by ... and n trailing characters
func Ellipsoider(s string, nx int) string {

	if len(s) == 0 {
		return "[empty]"
	}

	if len(s) <= 2*nx {
		return s
	}

	// len(s) > 2*nx
	const ellip = "..."
	return fmt.Sprintf("%v%v%v", s[:nx-1], ellip, s[len(s)-nx+1:])

}

// SplitByWhitespace splits by *any* combination of \t \n or space
func SplitByWhitespace(s1 string) (s2 []string) {
	s1 = strings.TrimSpace(s1)
	s2 = regexp.MustCompile(`[\s]+`).Split(s1, -1)

	return
}

// Breaker breaks a string into n equal sized substrings
func Breaker(s string, nx int) []string {

	if len(s) == 0 {
		return make([]string, 0)
	}

	rows := len(s) / nx
	if (len(s) % nx) != 0 {
		rows++
	}
	var ret []string = make([]string, rows)
	for i := 0; i < rows; i++ {
		if i < rows-1 {
			ret[i] = s[i*nx : (i+1)*nx]

		} else {
			ret[i] = s[i*nx:]

		}
	}
	return ret

}

/*
	IncrementString takes the last Character or Symbol
	and "increments" it.

	This is for all datastore indexes where we want to
	filter by
		field >= someString
		field <  nextBiggerString


	Note: We assume that s is already converted to lower case,

	If we wanted maintain case sensitive filtering,
	then we would do something like
		uTFCodePointUpperCase :=  uTFCodePointLastChar - 'A' + 'a'

	And then we would construct four filters
		.Filter("title >=", "cowgirls")
		.Filter("title < ", "cowgirlt")
		.Filter("title >=", "Cowgirls")
		.Filter("title < ", "Cowgirlt")

*/
func IncrementString(s string) string {

	if s == "" {
		panic("Increment String is undefined for an empty string")
	}

	uTFCodePointLastChar, itsSize := utf8.DecodeLastRuneInString(s)
	if uTFCodePointLastChar == utf8.RuneError {
		panic(fmt.Sprint("Following string is invalid utf8: %q", s))
	}
	sReduced := s[:len(s)-itsSize]

	uTFCodePointLastChar++
	oneHigherChar := fmt.Sprintf("%c", uTFCodePointLastChar)

	return sReduced + oneHigherChar

}

// IndentedDump is the long awaited spew alternative, that is *safe*.
// It takes any structure and converts it to a hierarchical string.
// It has no external dependencies.
//
// Big disadvantage: no unexported fields.
// For unexported fields fall back to
//		fmt.Println(spew.Sdump(nd))
func IndentedDump(v interface{}) *string {

	var reverseJSONTagEscaping = strings.NewReplacer(`\u003c`, "<", `\u003e`, ">", `\n`, "\n")

	firstColLeftMostPrefix := " "
	bytes, err := json.MarshalIndent(v, firstColLeftMostPrefix, "\t")
	if err != nil {
		s := fmt.Sprintf("error indent: %v\n", err)
		return &s
	}
	s := reverseJSONTagEscaping.Replace(string(bytes))

	return &s
}

func SliceDumpI(sl [][]int) {
	for i := 0; i < len(sl); i++ {
		fmt.Printf("%2v: ", i)
		for j := 0; j < len(sl[i]); j++ {
			fmt.Printf("%2v %2v; ", j, sl[i][j])
		}
		fmt.Printf("\n")
	}
}

func init() {
	// log.Println(LowerCasedUnderscored(`" geh du alter Äsel äh? - "" `))
	// log.Println(LowerCasedUnderscored(` 'Theo - wir fahrn nach Łódź .PnG'`))
}
