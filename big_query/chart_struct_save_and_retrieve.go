package big_query

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/net/http/loghttp"
	"google.golang.org/appengine"
)

// Chart data
//
// Normally, datastore types are restricted.
// For intstance a
//    map[string]map[string]float64
// can not be a datastore field.
//
// Therefore, this package takes a complex struct
// and *globs* it into a byte array,
// quasi normalizing it.
//
// It then saves the byte array within a dsu.WrapBlob
//
// This way, any struct can be saved into into datastore
// using dsu.WrapBlob.
type CData struct {
	M          map[string]map[string]float64
	VPeriods   []string
	VLangs     []string
	F_max      float64
	unexported string
}

// http://stackoverflow.com/questions/12854125/go-how-do-i-dump-the-struct-into-the-byte-array-without-reflection

func SaveChartDataToDatastore(w http.ResponseWriter, r *http.Request, cd CData, key string) string {

	c := appengine.NewContext(r)

	internalType := fmt.Sprintf("%T", cd)
	//buffBytes, _	 := StringToVByte(s)  // instead of []byte(s)

	// CData to []byte
	serializedStruct := new(bytes.Buffer)
	enc := gob.NewEncoder(serializedStruct)
	err := enc.Encode(cd)
	loghttp.E(w, r, err, false)

	key_combi, err := dsu.BufPut(c,
		dsu.WrapBlob{Name: key, VByte: serializedStruct.Bytes(), S: internalType}, key)
	loghttp.E(w, r, err, false)

	return key_combi
}

func GetChartDataFromDatastore(w http.ResponseWriter, r *http.Request, key string) *CData {

	c := appengine.NewContext(r)

	key_combi := "dsu.WrapBlob__" + key

	dsObj, err := dsu.BufGet(c, key_combi)
	loghttp.E(w, r, err, false)

	serializedStruct := bytes.NewBuffer(dsObj.VByte)
	dec := gob.NewDecoder(serializedStruct)
	newCData := new(CData) // hell, it was set to nil above - causing this "unassignable value" error
	err = dec.Decode(newCData)
	loghttp.E(w, r, err, false, "VByte was ", dsObj.VByte[:10])

	return newCData
}

func testGobDecodeEncode(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	nx := 24

	// without custom implementation
	// everything is encoded/decoded except field unexported
	orig := CData{
		M:          map[string]map[string]float64{"lang1": map[string]float64{"2012-09": 0.2}},
		VPeriods:   []string{"2011-11", "2014-11"},
		VLangs:     []string{"C", "++"},
		F_max:      44.2,
		unexported: "val of unexported",
	}
	fmt.Fprintf(w, "orig\n%#v\n", &orig)

	// writing to []byte
	serializedStruct := new(bytes.Buffer)
	enc := gob.NewEncoder(serializedStruct)
	err := enc.Encode(orig)
	loghttp.E(w, r, err, false)

	sx := serializedStruct.String()
	lx := len(sx)
	fmt.Fprintf(w, "byte data: \n%#v...%#v\n", sx[0:nx], sx[lx-nx:])

	// saving to ds
	key_combi, err := dsu.BufPut(c,
		dsu.WrapBlob{Name: "chart_data_test_1", VByte: serializedStruct.Bytes(), S: "chart data"}, "chart_data_test_1")
	loghttp.E(w, r, err, false)
	// restoring from ds
	dsObj, err := dsu.BufGet(c, key_combi)
	loghttp.E(w, r, err, false)

	p := r.FormValue("p")

	// reading
	rest1 := new(CData)
	if p == "" {
		sx := string(dsObj.VByte)
		lx := len(sx)
		fmt.Fprintf(w, "byte data: \n%#v...%#v\n", sx[0:nx], sx[lx-nx:])

		readr := bytes.NewBuffer(dsObj.VByte)
		dec := gob.NewDecoder(readr)
		err = dec.Decode(rest1)
		loghttp.E(w, r, err, false)
	} else {
		readr := bytes.NewBuffer(serializedStruct.Bytes())
		dec := gob.NewDecoder(readr)
		err = dec.Decode(rest1)
		loghttp.E(w, r, err, false)
	}

	fmt.Fprintf(w, "resl\n%#v\n", rest1)

	fmt.Fprintf(w, "\n\n")
	SaveChartDataToDatastore(w, r, orig, "chart_data_test_2")

	dsObj2, err := dsu.BufGet(c, "dsu.WrapBlob__chart_data_test_2")
	loghttp.E(w, r, err, false)
	{

		rest2 := new(CData)

		sx := string(dsObj2.VByte)
		lx := len(sx)
		fmt.Fprintf(w, "byte data: \n%#v...%#v\n", sx[0:nx], sx[lx-nx:])

		readr := bytes.NewBuffer(dsObj2.VByte)
		dec := gob.NewDecoder(readr)
		err = dec.Decode(rest2)
		loghttp.E(w, r, err, false)

		fmt.Fprintf(w, "res2\n%#v\n", rest2)

	}

	f1 := GetChartDataFromDatastore(w, r, "chart_data_test_2")
	fmt.Fprintf(w, "resl\n%#v\n", f1)

}

//  if we wanted to gob.Encode/Decode unexported fields,
//  like CData.unexported, then we have to implement
//  every field ourselves
//  => uncomment following ...
/*
func (d *CData)GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(d.unexported)
	if err!=nil { return nil, err }
	return w.Bytes(), nil
}

func (d *CData)GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	return  decoder.Decode(&d.unexported)
}
*/
