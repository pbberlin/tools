package util

/* This is a closure,
  an anonymous func
  with surrounding variables

  myIntSeq01(), myIntSeq02()
	 yield independent values for i

They also demonstrate "static instance memory",
as the "global" variables are kept as long as the app lives

*/
func intSeq() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

var MyIntSeq01 func() int = intSeq()
var MyIntSeq02 func() int = intSeq()
