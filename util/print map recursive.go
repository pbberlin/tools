package util

import "fmt"
import "bytes"

// this does not work generically:
var Map_example_expl_1 map[float64]map[int]float64

// this is generic - but unnecessary:
var Map_example_expl_2 map[interface{}]map[interface{}]interface{}

var Map_example_right map[interface{}]interface{}

func init() {

	Map_example_expl_1 = map[float64]map[int]float64{
		10:  map[int]float64{0: 0, 1: 2},
		7.5: map[int]float64{0: 0, 1: 1.5},
		5.0: map[int]float64{0: 0, 1: 1},
	}

	Map_example_expl_2 = map[interface{}]map[interface{}]interface{}{
		10:  map[interface{}]interface{}{0: 0, 1: 2},
		7.5: map[interface{}]interface{}{0: 0, 1: 1.5},
	}

	Map_example_right = map[interface{}]interface{}{
		"scalar_k1":                        "scalar_val",
		"scalar_k2":                        0.11111,
		"key to invalidly typed submap =>": map[interface{}]float64{"xx": 0.111, 1: 0.5},
		"lvl 1 submap":                     map[interface{}]interface{}{0: 0, 1: 2},
		32168:                              map[interface{}]interface{}{"any type": 0, 1: "anywhere"},
		"lvl 2 submaps": map[interface{}]interface{}{
			"l1a": map[interface{}]interface{}{"l2a": 0, 1: 0.5},
			"l1b": map[interface{}]interface{}{0: 0, 1: 0.5},
			13:    "mixing in a scalar",
		},
	}

}

//  PrintMap does print a map without recursion
func PrintMap(m map[interface{}]interface{}) string {
	return PrintMapRecursive(m, 0)
}

//  PrintMap prints a map - recursing through submaps
// concrete types, like map[float64]map[int]float64
//	  would have to be converted element by element
//   to map[interface{}]interface{}
// the json.indent function might be similar
// but I did not bring it to work
func PrintMapRecursive(m1 map[interface{}]interface{}, depth int) string {
	b1 := new(bytes.Buffer)
	fiveSpaces := "&nbsp;  &nbsp;  &nbsp;" // html spaces
	cxIndent := fmt.Sprint(depth * 40)
	depthStyle := ""
	if depth == 0 {
		depthStyle = "margin-top:8px; border-top: 1px solid #aaa;"
	} else {
		darkening := 16 - 2*depth
		bgc := fmt.Sprintf("%x%x%x", darkening, darkening, darkening)
		depthStyle = "margin:0;padding:0;background-color:#" + bgc + ";"
	}
	for k1, m2 := range m1 {
		b1.WriteString("<div style='" + depthStyle + "margin-left:" + cxIndent + "px;'>\n")
		b1.WriteString(fmt.Sprint(k1))
		//b1.WriteString( "<br>\n" )
		switch typedvar := m2.(type) {
		default:
			type_unknown := fmt.Sprintf("<br>\n &nbsp; &nbsp;  --this type switch does not support type %#v", typedvar)
			//panic(type_unknown)
			b1.WriteString(fiveSpaces + fmt.Sprintf("%#v", typedvar) + type_unknown)
		case string:
			b1.WriteString(fiveSpaces + typedvar)
		case int, float64:
			b1.WriteString(fiveSpaces + fmt.Sprint(typedvar))
		case nil:
			b1.WriteString("nil interface")
		case map[interface{}]interface{}:
			b1.WriteString(PrintMapRecursive(typedvar, depth+1))
		}
		b1.WriteString("<br>\n")
		b1.WriteString("</div>\n")
	}
	return b1.String()
}

/*
func Test_print_map(t *testing.T){

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Errorf("could not get a context")
	}

	s := PrintMap(Map_example_right)
	c.Infof("testing print map recursive ...")
	if  Test_want != s {
		c.Errorf("want %s - got %s", Test_want, s )
	}
}
*/

const Test_want = `<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
scalar_k1&nbsp;  &nbsp;  &nbsp;scalar_val<br>
</div>
<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
scalar_k2&nbsp;  &nbsp;  &nbsp;0.11111<br>
</div>
<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
key to invalidly typed submap =>&nbsp;  &nbsp;  &nbsp;map[interface {}]float64{"xx":0.111, 1:0.5}<br>
 &nbsp; &nbsp;  --this type switch does not support type map[interface {}]float64{"xx":0.111, 1:0.5}<br>
</div>
<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
lvl 1 submap<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
0&nbsp;  &nbsp;  &nbsp;0<br>
</div>
<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
1&nbsp;  &nbsp;  &nbsp;2<br>
</div>
<br>
</div>
<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
32168<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
any type&nbsp;  &nbsp;  &nbsp;0<br>
</div>
<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
1&nbsp;  &nbsp;  &nbsp;anywhere<br>
</div>
<br>
</div>
<div style='margin-top:8px; border-top: 1px solid #aaa;margin-left:0px;'>
lvl 2 submaps<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
l1a<div style='margin:0;padding:0;background-color:#ccc;margin-left:80px;'>
l2a&nbsp;  &nbsp;  &nbsp;0<br>
</div>
<div style='margin:0;padding:0;background-color:#ccc;margin-left:80px;'>
1&nbsp;  &nbsp;  &nbsp;0.5<br>
</div>
<br>
</div>
<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
l1b<div style='margin:0;padding:0;background-color:#ccc;margin-left:80px;'>
0&nbsp;  &nbsp;  &nbsp;0<br>
</div>
<div style='margin:0;padding:0;background-color:#ccc;margin-left:80px;'>
1&nbsp;  &nbsp;  &nbsp;0.5<br>
</div>
<br>
</div>
<div style='margin:0;padding:0;background-color:#eee;margin-left:40px;'>
13&nbsp;  &nbsp;  &nbsp;mixing in a scalar<br>
</div>
<br>
</div>
`
