package transposablematrix

type SearchResult int

const (
	nothing SearchResult = iota
	success
	failure
)

var AmorphsExhaustedError = epf("Amorphs exhausted")

func (c SearchResult) String() string {
	switch c {
	case nothing:
		return "nothing"
	case success:
		return "success"
	case failure:
		return "failure"
	}
	return ""
}
