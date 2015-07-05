package backend

import (
	sc "github.com/pbberlin/tools/dsu/distributed_unancestored"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	ss "github.com/pbberlin/tools/shared_structs"

	"appengine"

	// "github.com/davecgh/go-spew/spew"
	"net/http"
	"sort"
	"strings"

	"github.com/pbberlin/tools/util"
)

/*
	I would love to use slices of *interfaces*,
	instead of fixed types

	VB1        []interface{}   <---- preferable
	VB1        []b1            <---- restrictive

	But it is complicated to reach
	the pointers to the actual structs,
	if the structs are wrapped behind an interface.

		myB1 := &(myB0.VB1a[i1])
		myB1.Fieldname += "_BBBBB" // works

		myB1, ok := myB0.VB1[i1].(b1)
		if !ok ...
		myB1.Fieldname += "_BBBBB" // not working

		myB1 := &(myB0.VB1[i1].(b1))   // impossible - "cannot take the address ..."
		myB1 := &(myB0.VB1[i1]).(*b1)  // impossible

		(myB0.VB1[i1].(b1)).Heading += "_AAAAA" // impossible - "can not assign ..."


	This solution seems not applicable without writing accessor methods for each field:
	https://groups.google.com/forum/#!topic/golang-nuts/bxbUdSYdTSI


	These are possible solutions, both using reflection
	http://stackoverflow.com/questions/6395076/in-golang-using-reflect-how-do-you-set-the-value-of-a-struct-field
	http://stackoverflow.com/questions/21379466/golang-get-a-pointer-to-a-field-of-a-struct-trough-an-interface

	I implemented the first solution in util.GetIntField(),
	but it's not working with generic arguments.

	The second solution is too much complication for me -
	 I fall back to explicit typing:

	for i1, _ := range myB0.VB1 {

		myB0.VB1a[i1].Heading += "_AAAAA" // works

		myB1 := &(myB0.VB1a[i1])
		myB1.Heading += "_BBBBB" // works

		myB2 := (&myB0).VB1a[i1]
		myB2.Heading += "_CCCCC" // works NOT

		myB3 := myB0.VB1a[i1]
		myB3.Heading += "_DDDDD" // works NOT

	}

*/

//var myB0 = b0{NumYSectors: 2}
var myB0 ss.B0 = ss.B0{NumYSectors: 2}

func backend3(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	var nColsBlock = 4
	if r.FormValue("nColsBlock") != "" {
		nColsBlock = util.Stoi(r.FormValue("nColsBlock"))
	}
	var nColsViewport = 6
	if r.FormValue("nColsViewport") != "" {
		nColsViewport = util.Stoi(r.FormValue("nColsViewport"))
	}

	myB0.VB1 = ss.X
	myB0.NumB1 = len(myB0.VB1)
	myB0.NumB2 = 0
	myB0.NColsViewport = nColsViewport

	// compute basic meta data
	for i1, _ := range myB0.VB1 {
		myB0.NumB2 += len(myB0.VB1[i1].VB2)
		for i2, _ := range myB0.VB1[i1].VB2 {
			// number of chars
			ro := myB0.VB1[i1].VB2[i2] // read only
			myB0.VB1[i1].VB2[i2].Size = len(ro.Linktext) + len(ro.Description)
			myB0.VB1[i1].VB2[i2].EditorialIndex = i2
		}
	}

	// compute NCols - NRows for the block
	for i1, _ := range myB0.VB1 {
		myB0.VB1[i1].NCols = nColsBlock
		if myB0.VB1[i1].NColsEditorial > 0 {
			myB0.VB1[i1].NCols = myB0.VB1[i1].NColsEditorial
		}
		if len(myB0.VB1[i1].VB2) < nColsBlock && len(myB0.VB1[i1].VB2) > 0 {
			myB0.VB1[i1].NCols = len(myB0.VB1[i1].VB2)
		}
		myB0.VB1[i1].NRows = complementRowsOrCols(len(myB0.VB1[i1].VB2), myB0.VB1[i1].NCols)
		myB0.VB1[i1].Discrepancy = myB0.VB1[i1].NCols*myB0.VB1[i1].NRows - len(myB0.VB1[i1].VB2)

		myB0.MaxNCols = util.Max(myB0.MaxNCols, myB0.VB1[i1].NCols)
		myB0.MaxNRows = util.Max(myB0.MaxNRows, myB0.VB1[i1].NRows)
	}

	// compute NCols - NRows - sizeup to MaxNRows
	for i1, _ := range myB0.VB1 {
		if myB0.VB1[i1].NRows < myB0.MaxNRows {
			myB0.VB1[i1].NRows = myB0.MaxNRows
			myB0.VB1[i1].NCols = complementRowsOrCols(len(myB0.VB1[i1].VB2), myB0.VB1[i1].NRows)
			myB0.VB1[i1].Discrepancy = myB0.VB1[i1].NCols*myB0.VB1[i1].NRows - len(myB0.VB1[i1].VB2)
		}
	}

	// is first or last
	for i1, _ := range myB0.VB1 {
		for i2, _ := range myB0.VB1[i1].VB2 {
			myB0.VB1[i1].VB2[i2].IsFirst = false
			myB0.VB1[i1].VB2[i2].IsLast = false
			if i2%myB0.VB1[i1].NCols == 0 {
				myB0.VB1[i1].VB2[i2].IsFirst = true
			}
			if i2%myB0.VB1[i1].NCols == (myB0.VB1[i1].NCols - 1) {
				myB0.VB1[i1].VB2[i2].IsLast = true
			}
			//c.Infof("first-last %v %v \n", i2, i2%myB0.VB1[i1].NCols)
		}
	}

	// create slices with the data to be sorted
	for i1, _ := range myB0.VB1 {
		sh1 := make([]ss.Order, len(myB0.VB1[i1].VB2))
		myB0.VB1[i1].BySize = ss.ByInt(sh1)
		sh2 := make([]ss.Order, len(myB0.VB1[i1].VB2))
		myB0.VB1[i1].ByHeading = ss.ByStr(sh2)
		// fill in the data - to be sorted later
		for i2, _ := range myB0.VB1[i1].VB2 {
			ro := myB0.VB1[i1].VB2[i2] // read only
			myB0.VB1[i1].BySize[i2].IdxSrc = i2
			myB0.VB1[i1].BySize[i2].ByI = len(ro.Linktext) + len(ro.Description)
			myB0.VB1[i1].ByHeading[i2].IdxSrc = i2
			myB0.VB1[i1].ByHeading[i2].ByS = strings.ToLower(ro.Linktext)
		}
	}

	// actual rearranging of the sorting date
	for i1, _ := range myB0.VB1 {
		sort.Sort(myB0.VB1[i1].BySize)
		sort.Sort(myB0.VB1[i1].ByHeading)
		c.Infof("-- Sorting %v", myB0.VB1[i1].Heading)
		// for i, v := range myB0.VB1[i1].BySize {
		// 	c.Infof("---- %v %v %v", i, v.IdxSrc, v.ByI)
		// }
		// for i, v := range myB0.VB1[i1].ByHeading {
		// 	c.Infof("---- %v %v %v", i, v.IdxSrc, v.ByS)
		// }
	}

	path := m["dir"].(string) + m["base"].(string)

	cntr, _ := sc.Count(c, path)

	add, tplExec := tplx.FuncTplBuilder(w, r)
	add("n_html_title", "Backend", nil)

	add("n_cont_0", "<style>"+htmlfrag.CSSColumnsWidth(nColsViewport)+"</style>", "")
	add("n_cont_1", tplx.PrefixLff+"backend3_body", myB0)
	add("tpl_legend", tplx.PrefixLff+"backend3_body_embed01", "")
	add("n_cont_2", "<p>{{.}} views</p>", cntr)

	sDumped := ""
	//sDumped = spew.Sdump(myB0)
	add("n_cont_3", "<pre>{{.}} </pre>", sDumped)

	tplExec(w, r)

}

func prepareLayout(l ss.B0) {

}

func complementRowsOrCols(nBlocks int, nRowOrCol int) int {
	if nRowOrCol < 1 {
		panic("Count of Cols or Rows must be given")
	}
	nCompl := nBlocks / nRowOrCol
	if nBlocks%nRowOrCol != 0 {
		nCompl++
	}
	return nCompl

}

func init() {
	prepareLayout(myB0)
	http.HandleFunc("/backend3", loghttp.Adapter(backend3))
}
