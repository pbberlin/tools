package big_query

// https://godoc.org/code.google.com/p/google-api-go-client/bigquery/v2
// https://developers.google.com/bigquery/bigquery-api-quickstart
import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"fmt"
	"net/http"

	bq "code.google.com/p/google-api-go-client/bigquery/v2"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"appengine"

	newappengine "google.golang.org/appengine" // https://github.com/golang/oauth2
)

// print it to http writer
func printPlaintextTable(w http.ResponseWriter, r *http.Request, vVDest [][]byte) {

	//c := appengine.NewContext(r)
	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(b1.Bytes())
	}()

	for i0 := 0; i0 < len(vVDest); i0++ {
		b1.WriteString("--")
		b1.Write(vVDest[i0])
		b1.WriteString("--")
		b1.WriteString("\n")
	}

}

func queryIntoDatastore(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	limitUpper := util.MonthsBack(1)
	limitLower := util.MonthsBack(25)

	var q bq.QueryRequest = bq.QueryRequest{}
	q.Query = `
		SELECT
		  repository_language
		, LEFT(repository_pushed_at,7) monthx
		, CEIL( count(*)/1000) Tausend
		FROM githubarchive:github.timeline
		where 1=1
			AND  LEFT(repository_pushed_at,7) >= '` + limitLower + `'
			AND  LEFT(repository_pushed_at,7) <= '` + limitUpper + `'
			AND  repository_language in ('Go','go','Golang','golang','C','Java','PHP','JavaScript','C++','Python','Ruby')
			AND  type="PushEvent"
		group by monthx, repository_language
		order by repository_language   , monthx
		;
	`

	c := appengine.NewContext(r)

	// The following client will be authorized by the App Engine
	// app's service account for the provided scopes.
	// "https://www.googleapis.com/auth/bigquery"
	// "https://www.googleapis.com/auth/devstorage.full_control"

	// 2015-06: instead of oauth2.NoContext we get a new type of context
	var ctx context.Context = newappengine.NewContext(r)
	oauthHttpClient, err := google.DefaultClient(
		ctx, "https://www.googleapis.com/auth/bigquery")

	if err != nil {
		log.Fatal(err)
	}

	bigqueryService, err := bq.New(oauthHttpClient)

	util_err.Err_http(w, r, err, false)

	fmt.Fprint(w, "s1<br>\n")

	// Create a query statement and query request object
	//  query_data = {'query':'SELECT TOP(title, 10) as title, COUNT(*) as revision_count FROM [publicdata:samples.wikipedia] WHERE wp_namespace = 0;'}
	//  query_request = bigquery_service.jobs()
	// Make a call to the BigQuery API
	//  query_response = query_request.query(projectId=PROJECT_NUMBER, body=query_data).execute()

	js := bq.NewJobsService(bigqueryService)
	jqc := js.Query("347979071940", &q)

	fmt.Fprint(w, "s2 "+util.TimeMarker()+" <br>\n")
	resp, err := jqc.Do()
	util_err.Err_http(w, r, err, false)

	rows := resp.Rows
	var vVDest [][]byte = make([][]byte, len(rows))

	c.Errorf("%#v", rows)

	for i0, v0 := range rows {

		cells := v0.F

		b_row := new(bytes.Buffer)
		b_row.WriteString(fmt.Sprintf("r%0.2d -- ", i0))
		for i1, v1 := range cells {
			val1 := v1.V
			b_row.WriteString(fmt.Sprintf("c%0.2d: %v  ", i1, val1))
		}
		vVDest[i0] = []byte(b_row.Bytes())
	}

	key_combi, _ := dsu.BufPut(w, r, dsu.WrapBlob{Name: "bq_res1", VVByte: vVDest}, "bq_res1")
	dsObj, _ := dsu.BufGet(w, r, key_combi)

	printPlaintextTable(w, r, dsObj.VVByte)

	fmt.Fprint(w, "s3 "+util.TimeMarker()+" <br>\n")

}

func mockDateIntoDatastore(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	rand.Seed(time.Now().UnixNano())

	row_max := 100
	col_max := 3

	var languages []string = []string{"C", "C++", "Rambucto"}

	var vVDest [][]byte = make([][]byte, row_max)
	for i0 := 0; i0 < row_max; i0++ {

		vVDest[i0] = make([]byte, col_max)

		b_row := new(bytes.Buffer)
		b_row.WriteString(fmt.Sprintf("r%0.2d -- ", i0))

		for i1 := 0; i1 < col_max; i1++ {
			if i1 == 0 {
				val := languages[i0/10%3]
				b_row.WriteString(fmt.Sprintf(" c%0.2d: %-10.8v  ", i1, val))
			} else if i1 == 2 {
				val := rand.Intn(300)
				b_row.WriteString(fmt.Sprintf(" c%0.2d: %10v  ", i1, val))
			} else {

				f2 := "2006-01-02 15:04:05"
				f2 = "2006-01"
				tn := time.Now()
				//tn  = tn.Add( - time.Hour * 85 *24 )
				tn = tn.Add(-time.Hour * time.Duration(i0) * 24)
				val := tn.Format(f2)
				b_row.WriteString(fmt.Sprintf(" c%0.2d: %v  ", i1, val))
			}
		}
		vVDest[i0] = []byte(b_row.Bytes())

	}

	key_combi, _ := dsu.BufPut(w, r, dsu.WrapBlob{Name: "bq_res_test", VVByte: vVDest}, "bq_res_test")
	dsObj, _ := dsu.BufGet(w, r, key_combi)

	printPlaintextTable(w, r, dsObj.VVByte)

}

func regroupFromDatastore01(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-Type", "text/html")
		w.Write(b1.Bytes())
	}()

	var vVSrc [][]byte

	if util_appengine.IsLocalEnviron() {
		vVSrc = bq_statified_res1
	} else {
		dsObj1, _ := dsu.BufGet(w, r, "dsu.WrapBlob__bq_res1")
		vVSrc = dsObj1.VVByte
	}

	if r.FormValue("mock") == "1" {
		dsObj1, _ := dsu.BufGet(w, r, "dsu.WrapBlob__bq_res_test")
		vVSrc = dsObj1.VVByte
	}

	var vVDest [][]byte = make([][]byte, len(vVSrc))

	for i0 := 0; i0 < len(vVSrc); i0++ {

		s_row := string(vVSrc[i0])
		v_row := pbstrings.SplitByWhitespace(s_row)
		b_row := new(bytes.Buffer)

		b_row.WriteString(fmt.Sprintf("%16.12s   ", v_row[3])) // leading spaces
		b_row.WriteString(fmt.Sprintf("%16.12s   ", v_row[5]))
		b_row.WriteString(fmt.Sprintf("%16.8s", v_row[7]))

		vVDest[i0] = []byte(b_row.Bytes())

	}

	key_combi, _ := dsu.BufPut(w, r, dsu.WrapBlob{Name: "res_processed_01", S: "[][]byte", VVByte: vVDest}, "res_processed_01")
	dsObj2, _ := dsu.BufGet(w, r, key_combi)

	printPlaintextTable(w, r, dsObj2.VVByte)

}

func regroupFromDatastore02(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-Type", "text/html")
		w.Write(b1.Bytes())
	}()

	var vVSrc [][]byte
	dsObj1, err := dsu.BufGet(w, r, "dsu.WrapBlob__res_processed_01")
	util_err.Err_http(w, r, err, false)
	vVSrc = dsObj1.VVByte

	d := make(map[string]map[string]float64)

	distinctLangs := make(map[string]interface{})
	distinctPeriods := make(map[string]interface{})
	f_max := 0.0
	for i0 := 0; i0 < len(vVSrc); i0++ {
		//vVDest[i0] = []byte( b_row.Bytes() )
		s_row := string(vVSrc[i0])
		v_row := pbstrings.SplitByWhitespace(s_row)

		lang := v_row[0]
		period := v_row[1]
		count := v_row[2]
		fCount := util.Stof(count)
		if fCount > f_max {
			f_max = fCount
		}

		distinctLangs[lang] = 1
		distinctPeriods[period] = 1

		if _, ok := d[period]; !ok {
			d[period] = map[string]float64{}
		}
		d[period][lang] = fCount

	}
	//fmt.Fprintf(w,"%#v\n",d2)
	//fmt.Fprintf(w,"%#v\n",f_max)

	sortedPeriods := util.StringKeysToSortedArray(distinctPeriods)
	sortedLangs := util.StringKeysToSortedArray(distinctLangs)

	cd := CData{}
	_ = cd

	cd.M = d
	cd.VPeriods = sortedPeriods
	cd.VLangs = sortedLangs
	cd.F_max = f_max

	SaveChartDataToDatastore(w, r, cd, "chart_data_01")

	/*
		if r.FormValue("f") == "table" {
			showAsTable(w,r,cd)
		} else {
			showAsChart(w,r,cd)
		}
	*/

}

func init() {
	http.HandleFunc("/big-query/query-into-datastore", util_appengine.Adapter(queryIntoDatastore))
	http.HandleFunc("/big-query/mock-data-into-datastore", util_appengine.Adapter(mockDateIntoDatastore))
	http.HandleFunc("/big-query/regroup-data-01", util_appengine.Adapter(regroupFromDatastore01))
	http.HandleFunc("/big-query/regroup-data-02", util_appengine.Adapter(regroupFromDatastore02))
}
