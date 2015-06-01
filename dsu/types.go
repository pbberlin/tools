package dsu

var memoryInstanceStore = make(map[string]*WrapBlob)

// WrapString for scalar values of type string
type WrapString struct {
	S string
}

// WrapInt for scalar values of type int
type WrapInt struct {
	I int
}

// WrapBlob - an universal struct for buffered put and get
//     for datastore and memcache
type WrapBlob struct {

	Name string // the key - dsu.WrapBlob__ 'some name'
	Desc string
	Category string

	// scalars
	S string  // Also contains the mime type
	I int
	F float64

	// for storage >> 500 byte
	VByte  []byte   // simple vector, for images + globbed structs
	VVByte [][]byte // vector of vector, i.e. for table data

}
