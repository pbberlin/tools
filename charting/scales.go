package charting

import (
	"bytes"
	"fmt"
	"math"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlpb"
	"github.com/pbberlin/tools/util"
)

// scales are characterized by the number of ticks
// being forced to set the number in advance
// we would resort to 10 ticks, or 12 ticks or 8 ticks, or 5 or 4.
// => yielding the maximum number of integers upon division

// once the maximum of ticks is specified
// we need to "inscribe" the tick *values*
// based on the maximum function value
// i.e. we could label a for-tick sale
// either 2-4-6-8 or 1-2-3-4,
// depending on max(f) being 7 or 4

// Thus we distill following variables
// => num_ticks
// => vector max_tick_values
// => f_max
// => best max_tick_value for f_max

// we will harcode several num_ticks - max_tick_value structures
// we will then compute all single tick values
// we will finally provide an exported function returning the
//  most appropriate max_tick_value for a given f_max and num_ticks

// in the case of time ranges:
//    integer index the times-points
//     0:2012-01, 1:2012-02 ... 23:2014-12
// then get a scale for the index values [0:23]+1
//   then provide a conversion function from the tick_values
//   0, 2.0 , 4.0 ... 12.0
// to the tick-values in time
//   i.e.   func(dateStart date,delta float)string{ return dateStart.AddDate( delta * date.Month).format("YY-mm") }

// vm stands for *v*ector of *m*aps
var Scale_x_vm, Scale_y_vm map[float64][]float64

func init() {

	// our standard chart has 5 y-axis ticks
	// and 12 x-axis ticks

	// we set possible max tick values
	// easily divisible by the number of ticks:
	Scale_y_vm = map[float64][]float64{
		7.5: nil,
		5.0: nil,
		2.5: nil,
		2.0: nil,
		1.0: nil,
	}
	prepareTickValues(Scale_y_vm, 5) // filling in 5 tick labels

	// preparing scales for the x-axis with 12 ticks ...
	Scale_x_vm = map[float64][]float64{
		9.0: nil,
		6.0: nil,
		2.4: nil,
		1.2: nil,
	}
	prepareTickValues(Scale_x_vm, 12)

	http.HandleFunc("/blob/scale-test", test)

}

// returns a format string with as few as needed
//   post-decimal digits ;  1000 => 1000 , but 0.0400 => 0.04
func practicalFormat(mv float64) (floatFormat string, exponent int) {

	//Log x    =    Ln x / Ln 10
	fExponent := math.Log(mv) / math.Log(10)
	//exponent  = int(fExponent)
	exponent = util.Round(fExponent)

	sExponent := fmt.Sprint(util.Abs(exponent) + 2)
	floatFormat = "%12.0f"
	if mv < 10 {
		floatFormat = "%12.1f"
	}
	if mv < 1 {
		floatFormat = "%12." + sExponent + "f"
	}

	return
}

// helper for outputting
func printScale(s map[float64][]float64) string {

	b1 := new(bytes.Buffer)
	b1.WriteString("<hr>\n")
	for max_val, vs := range s {
		b1.WriteString(fmt.Sprint("<b>", max_val, "</b><br>\n"))
		for i, val := range vs {
			quot := fmt.Sprintf("%-4.2f", val)
			if len(quot) > 0 && quot[len(quot)-1:] == "0" {
				quot = quot[:len(quot)-1]
			}
			b1.WriteString(fmt.Sprintf("<pre style='margin:0'> %-6d   %s</pre>\n", i, quot))
		}
		b1.WriteString("<br>\n")
	}
	return b1.String()
}

// takes number of ticks and a max value -
// and fills tick-values {0;num_ticks} with
// values for
func prepareTickValues(s map[float64][]float64, num_ticks int) {

	for max_val, vs := range s {
		vs = make([]float64, num_ticks+1)
		for i := 0; i <= num_ticks; i++ {
			ftick := max_val / float64(num_ticks)
			ftick_val := float64(i) * ftick
			vs[i] = ftick_val
		}
		s[max_val] = vs // unclear why this is necessary
	}
}

// for a given data set with f(x) = {0;f_max}
//    we need the the best scale of possibleScales
//   	for convenience we multiply tick label values for f_max
func BestScale(f_max float64, possibleScales map[float64][]float64) (bestScale []string, exponent int, msg string) {

	keyToBeestScale := 0.0
	b1 := new(bytes.Buffer)

	b1.WriteString(fmt.Sprintf("searching maxval for %#v<br>\n", f_max))

	// finding mantissa * 10^exponent for f_max
	//   we also could use     exponent := math.Log(f_max) / math.Log(10)
	lp := f_max
	if lp < 0 {
		lp = lp * -1
	}
	if lp >= 1 {
		for {
			lp = lp / 10
			if lp < 1 {
				break
			}
			exponent++
		}
	} else {
		exponent = -1
		for {
			lp = lp * 10
			if lp > 1 {
				break
			}
			exponent--
		}
	}
	mantissa := f_max / (math.Pow10(exponent))
	b1.WriteString(fmt.Sprintf("mantissa <b>%6.2f</b> - exponent  %#v<br>\n", mantissa, exponent))

	// now we find the scale best fitting our mantissa
	smallest_dist := 10.0
	for max_scale_val, _ := range possibleScales {
		lp_dist := max_scale_val - mantissa
		if lp_dist >= 0 && lp_dist <= smallest_dist {
			keyToBeestScale = max_scale_val
			smallest_dist = lp_dist
		}
		b1.WriteString(fmt.Sprintf("  &nbsp;  &nbsp; cealinged by %4.2f  --- dist %4.2f - new min %4.2f<br>\n", max_scale_val, lp_dist, smallest_dist))
	}
	b1.WriteString(fmt.Sprintf("found scale <b>%#v</b> in between<br>\n", keyToBeestScale))

	// overflow - special case
	//		find the "dead" range
	smallest_scale := 10.0
	largest_scale := 0.0
	for max_scale_val, _ := range possibleScales {
		if max_scale_val < smallest_scale {
			smallest_scale = max_scale_val
		}
		if max_scale_val > largest_scale {
			largest_scale = max_scale_val
		}
	}
	// if mantissa is in dead range - then scale up smallest tick
	if 10*smallest_scale > mantissa && mantissa > largest_scale {
		keyToBeestScale = smallest_scale
		exponent++
		b1.WriteString(fmt.Sprintf("<br>found scale <b>%#v</b> - in loop over\n", keyToBeestScale))
		b1.WriteString(fmt.Sprintf(" &nbsp;  &nbsp;  &nbsp; -- smallest scale %#v - largest scale  %#v -- <br>\n", smallest_scale, largest_scale))
	}

	b1.WriteString("<hr>\n")
	msg = b1.String()

	// now that we scale the tick label values according to the exponent

	floatFormat, _ := practicalFormat(f_max)

	factor := math.Pow10(exponent)
	lenSrc := len(possibleScales[keyToBeestScale])
	bestScale = make([]string, lenSrc)
	for i, v := range possibleScales[keyToBeestScale] {
		bestScale[i] = fmt.Sprintf(floatFormat, factor*v)
	}

	return

}

func test(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	//fmt.Fprint(w,  util.PrintMap(util.Map_example_right))

	// vector max values
	vMaxVal := []float64{5555.0094, 9400, 0.0094, 1.0, 110.94, 120, 0.00001233}

	scale_test := Scale_x_vm
	for _, mv := range vMaxVal {
		optScale, _, _ := BestScale(mv, scale_test)

		floatFormat, exp := practicalFormat(mv)

		funcSpanner := htmlpb.GetSpanner()

		fmt.Fprintf(w, "optimal scale for "+floatFormat+" (exp %d) is <br>", mv, exp)
		for tick, val := range optScale {
			dis1 := fmt.Sprint(tick)
			dis2 := fmt.Sprintf("%s<br>", val)

			next_idx := tick + 1
			if next_idx > len(optScale)-1 {
				next_idx = len(optScale) - 1
			}
			if util.Stof(val) <= mv && mv <= util.Stof(optScale[next_idx]) {
				dis2 = fmt.Sprintf("<b>%s</b> <br>", val)
			}
			//dis4a :=  fmt.Sprintf(" %f <=  "+floatFormat+" &&   "+floatFormat+" < %f<br>",	  util.Stof(val),mv,mv, util.Stof(optScale[next_idx]) )
			//fmt.Fprintf(w,funcSpanner(dis4a,433)	)

			fmt.Fprintf(w, funcSpanner(" ", 64))
			fmt.Fprintf(w, funcSpanner(dis1, 120))
			fmt.Fprintf(w, funcSpanner(dis2, 233))
			fmt.Fprintf(w, "<br>")

		}
		//fmt.Fprintf(w,"key %v -  pot %v <br> %v",key,pot,msg)

		fmt.Fprintf(w, "<br>")
	}

	fmt.Fprintf(w, printScale(Scale_y_vm))
	fmt.Fprintf(w, printScale(Scale_x_vm))

}
