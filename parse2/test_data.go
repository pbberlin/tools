package parse2

import "fmt"

var tests = make([]string, 2)

func init() {

	tests[0] = `<!DOCTYPE html><html><head>
		<script type="text/javascript" src="./article01_files/empty.js"></script>
		<link href="./article01_files/vendor.css" rel="stylesheet" type="text/css"/>
		</head><body><p>Links:
				<span>span01</span>
				<span>span02-line1<br>span02-line2</span>
				<span>span03</span>
			</p>
			<style> p {font-size:17px}</style>
			<ul>
				<li id='332' ><a   href="/some/first/page.html">Linktext1 <span>inside</span></a>
				<li><a   href="/snd/page" title="wi-title">LinkT2</a>
			</ul>
			<div>
				<div>div-1-content</div>
				<div>div-2-content</div>
				<p>pararaph in between</p>
				<div>div-3-content with iimmage<img alt="alt-cnt" title='title-cnt' 
				href='some-long-href-some-long-href-some-long-href-some-long-href'>after img</div>
			</div>
			</body></html>`

	tests[1] = `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>Gepriesen sei der Mann, der wenn er nichts</li>
				<li>zu sagen hat</li>
			</ul>`

	const offSetFilename = 4

	// write out to http doc root
	for i := 0; i < len(tests); i++ {
		fn := fmt.Sprintf(docRoot+"/handelsblatt.com/art%02v.html", i+offSetFilename)
		bytes2File(fn, []byte(tests[i]))
	}

}
