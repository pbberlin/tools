// Package foscam allows setting and
// clearing alarm mode of such cam via http.
package foscam

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"

	"appengine"
	"appengine/urlfetch"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

var serveraddress string = `header of ["SERVER_ADDR"]`
var pos_01 int = strings.Index(serveraddress, "192.")

var dns_router string = "see func init"
var dns_cam string = "see func init"

const (
	debug               = false
	credentialsExcluded = "&usr=[a...n]&pwd=[pb...05]"
	path_get_alarm      = "/cgi-bin/CGIProxy.fcgi?cmd=getMotionDetectConfig"
	path_set_alarm      = "/cgi-bin/CGIProxy.fcgi?cmd=setMotionDetectConfig"

	path_snap_config     = "/cgi-bin/CGIProxy.fcgi?cmd=setSnapConfig&snapPicQuality=0&saveLocation=2"
	path_snap_retrieval  = "/cgi-bin/CGIProxy.fcgi?cmd=snapPicture2" + credentials
	path_video_retrieval = "/cgi-bin/CGIStream.cgi?cmd=GetMJStream" + credentials

	path_get_log = "/cgi-bin/CGIProxy.fcgi?cmd=getLog&count=20&offset=0"
)

// we have to upper case of fields to really get the values - to annoyed to think why
type CGI_Result struct {
	Result          string `xml:"result"`
	IsEnable        string `xml:"isEnable"`
	Linkage         string `xml:"linkage"`
	SnapInterval    string `xml:"snapInterval"`
	Sensitivity     string `xml:"sensitivity"`
	TriggerInterval string `xml:"triggerInterval"`
	Schedule0       string `xml:"schedule0"`
	Schedule1       string `xml:"schedule1"`
	Schedule2       string `xml:"schedule2"`
	Schedule3       string `xml:"schedule3"`
	Schedule4       string `xml:"schedule4"`
	Schedule5       string `xml:"schedule5"`
	Schedule6       string `xml:"schedule6"`
	Area0           string `xml:"area0"`
	Area1           string `xml:"area1"`
	Area2           string `xml:"area2"`
	Area3           string `xml:"area3"`
	Area4           string `xml:"area4"`
	Area5           string `xml:"area5"`
	Area6           string `xml:"area6"`
	Area7           string `xml:"area7"`
	Area8           string `xml:"area8"`
	Area9           string `xml:"area9"`
	Log0            string `xml:"log0"`
	Log1            string `xml:"log1"`
	Log2            string `xml:"log2"`
	Log3            string `xml:"log3"`
	Log4            string `xml:"log4"`
	Log5            string `xml:"log5"`
	Log6            string `xml:"log6"`
	Log7            string `xml:"log7"`
	Log8            string `xml:"log8"`
	Log9            string `xml:"log9"`
	// XMLName  xml.Name `xml:"account"`
}

func urlParamTS() string {
	ts := time.Now().UnixNano()
	return spf("%v", ts)
}

func makeRequest(w http.ResponseWriter, r *http.Request, path string) CGI_Result {

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	url_exe := spf(`http://%s%s%s&ts=%s`, dns_cam, path, credentials, urlParamTS())
	url_dis := spf(`http://%s%s&ts=%s`, dns_cam, path, urlParamTS())
	wpf(w, "<div style='font-size:10px; line-height:11px;'>requesting %v<br></div>\n", url_dis)
	resp1, err := client.Get(url_exe)
	loghttp.E(w, r, err, false)

	bcont, err := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()
	loghttp.E(w, r, err, false)

	cgiRes := CGI_Result{}
	xmlerr := xml.Unmarshal(bcont, &cgiRes)
	loghttp.E(w, r, xmlerr, false)

	if cgiRes.Result != "0" {
		wpf(w, "<b>RESPONSE shows bad mood:</b><br>\n")
		psXml := stringspb.IndentedDump(cgiRes)
		dis := strings.Trim(psXml, "{}")
		wpf(w, "<pre style='font-size:10px;line-height:11px;'>%v</pre>", dis)
	}

	if debug {
		scont := string(bcont)
		wpf(w, "<pre style='font-size:10px;line-height:11px;'>%v</pre>", scont)
	}

	return cgiRes

}

func imageRetrieve(w http.ResponseWriter, r *http.Request) {

	makeRequest(w, r, path_snap_config)
	wpf(w, "<img src='http://%s%s' width='60%' /><br>", dns_cam, path_snap_retrieval)

}

func logRetrieve(w http.ResponseWriter, r *http.Request) {

	cgiRes := makeRequest(w, r, path_get_log)

	sl := []string{cgiRes.Log0, cgiRes.Log1, cgiRes.Log2, cgiRes.Log3, cgiRes.Log4,
		cgiRes.Log5, cgiRes.Log6, cgiRes.Log7, cgiRes.Log8, cgiRes.Log9}

	for _, v := range sl {
		sl1 := strings.Split(v, "%2B")
		if len(sl1) < 4 {
			continue
		}

		// 		 time+user+ip+logID
		unixTS := sl1[0]
		usr := sl1[1]
		ip := sl1[2]
		eventId := sl1[3]
		eventDesc := ""
		switch eventId {
		case "0":
			eventDesc = "Power On"
		case "1":
			eventDesc = "Motion Alarm"
		case "3":
			eventDesc = "Login"
		case "4":
			eventDesc = "Logout"
		case "5":
			eventDesc = "Offline"
		default:
			eventDesc = "unkown event id: " + eventId
		}
		_, _, _ = eventDesc, usr, ip

		ts := util.TimeFromUnix(unixTS)
		tsf := ts.Format("2.1.2006 15:04:05")

		tn := time.Now()
		since := tn.Sub(ts)
		iHours := int(math.Floor(since.Hours()))
		iMinutes := util.Round(since.Minutes()) - iHours*60

		if eventId == "1" {
			wpf(w, "Last Alarm <b>%3vhrs %2vmin</b> ago (%v)<br>\n", iHours, iMinutes, tsf)
			break
		}
	}

}

func foscamStatus(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	htmlfrag.SetNocacheHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	logRetrieve(w, r)

	cgiRes := makeRequest(w, r, path_get_alarm)

	psXml := stringspb.IndentedDump(cgiRes)
	dis := strings.Trim(psXml, "{}")
	dis = strings.Replace(dis, "\t", "", -1)
	dis = strings.Replace(dis, " ", "", -1)
	dis = strings.Replace(dis, "\"", "", -1)
	dis = strings.Replace(dis, "\n", " ", -1)
	dis = strings.Replace(dis, "Area0", "\nArea0", -1)
	dis = strings.Replace(dis, "Schedule0", "\nSchedule0", -1)
	dis = strings.Replace(dis, "Log0", "\nLog0", -1)
	wpf(w, "<pre style='font-size:10px;line-height:11px;'>%v</pre>", dis)

	if cgiRes.IsEnable == "0" {
		wpf(w, "Status <b>DISabled</b><br><br>\n")
	} else {
		wpf(w, "Status <b>ENabled</b><br><br>\n")
	}

	imageRetrieve(w, r)

}

func foscamToggle(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	htmlfrag.SetNocacheHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ssecs := r.FormValue("sleep")
	if ssecs != "" {
		secs := util.Stoi(ssecs)
		wpf(w, "sleeping %v secs ... <br><br>\n", secs)
		time.Sleep(time.Duration(secs) * time.Second)
	}

	prevStat := makeRequest(w, r, path_get_alarm)

	wpf(w, "||%s||<br>\n", prevStat.IsEnable)
	if strings.TrimSpace(prevStat.IsEnable) == "0" {
		prevStat.IsEnable = "1"
	} else {
		prevStat.IsEnable = "0"
	}
	prevStat.Area0 = "255"
	prevStat.Area1 = "255"
	prevStat.Area2 = "255"
	prevStat.Area3 = "255"
	prevStat.Area4 = "255"
	prevStat.Area5 = "255"
	prevStat.Area6 = "255"
	prevStat.Area7 = "255"
	prevStat.Area8 = "255"
	prevStat.Area9 = "255"

	// ugly: XML dump to query string
	s2 := spf("%+v", prevStat)
	s2 = strings.Trim(s2, "{}")
	s2 = strings.Replace(s2, ":", "=", -1)
	s2 = strings.Replace(s2, " ", "&", -1)

	// even worse: we have to lower the case again
	pairs := strings.Split(s2, "&")
	recombined := ""
	for i, v := range pairs {
		fchar := v[:1]
		fchar = strings.ToLower(fchar)
		recombined += fchar + v[1:]
		if i < len(pairs)-1 {
			recombined += "&"
		}
	}

	wpf(w, "<pre>")
	// disS2 := stringspb.Breaker(s2, 50)
	// for _, v := range disS2 {
	// 	wpf(w, "%v\n", v)
	// }
	disRecombined := stringspb.Breaker(recombined, 50)
	for _, v := range disRecombined {
		wpf(w, "%v\n", v)
	}
	wpf(w, "</pre>")
	// wpf(w, "<pre>%v</pre>\n", recombined)

	toggleRes := makeRequest(w, r, path_set_alarm+"&"+recombined)
	if toggleRes.Result == "0" {
		wpf(w, "<br>end foscam toggle - success<br>\n")
		if prevStat.IsEnable == "0" {
			wpf(w, "<b>DISabled</b><br>\n")
		} else {
			wpf(w, "<b>ENabled</b><br>\n")
		}
	}

}

func foscamWatch(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	htmlfrag.SetNocacheHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Foscam live watch"}))

	/*

		There is no way to access a real video stream.
		Thus we use this suggestion: http://foscam.us/forum/post43654.html#p43654

	*/
	str := `<img 
	width='640' 
	src="http://` + dns_cam + `/CGIProxy.fcgi?cmd=snapPicture2&usr=visitor&pwd=visitor&t=" 
	onload='setTimeout(function() {src = src.substring(0, (src.lastIndexOf("t=")+2))+(new Date()).getTime()}, 1000)' 
	onerror='setTimeout(function() {src = src.substring(0, (src.lastIndexOf("t=")+2))+(new Date()).getTime()}, 5000)' 
	alt='' />`
	w.Write([]byte(str))

	w.Write([]byte(tplx.Foot))

}

func init() {

	if util_appengine.IsLocalEnviron() {
		dns_router = "192.168.1.1"
		dns_cam = "192.168.1.4:8081"
	} else {
		dns_router = "ds7934.myfoscam.org"
		dns_cam = "ds7934.myfoscam.org:8081"
	}

	http.HandleFunc("/foscam-status", loghttp.Adapter(foscamStatus))
	http.HandleFunc("/foscam-toggle", loghttp.Adapter(foscamToggle))
	http.HandleFunc("/foscam-watch", loghttp.Adapter(foscamWatch))

}
