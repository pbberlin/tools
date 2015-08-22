// +build split
// go test -tags=split

package domclean2

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"
	"golang.org/x/net/html"
)

var testDocs = make([]string, 2)

func init() {

	testDocs[0] = `<!DOCTYPE html><html><head>
		</head><body>

			<p>
				<a href="/some/first/page.html">Links1:
					no img
				</a>
				<a href="/some/first/page.html">Links2:
					<span>text bef</span>
					<img src="/img1src" />
					<span>text aft</span>
				</a>
				
			</p>
			<div>
				<div>
					<a href="/some/first/page.html">Links3:
						<span>text2 bef</span>
						<span>
							<span>text3 bef</span>
							<img src="/img1src" />
							<span>text3 aft</span>
						</span>
						<span>text2 aft</span>
					</a>
				</div>
			</div>
		</body></html>`

	testDocs[1] = `<!DOCTYPE html><html><head></head><body>
	<p>
		<a href="/some/first/page.html">Links1:
					no img 
		</a>
		<a href="/some/first/page.html">Links2: 
			<span>text bef 
			</span>
		</a>
		<img src="/img1src"/>
		<a href="/some/first/page.html">
			<span>text aft 
			</span>
		</a>
	</p>
	<div>
		<div>
			<a href="/some/first/page.html">Links3: 
				<span>text2 bef 
				</span>
				<span>
					<span>text3 bef 
					</span>
				</span>
			</a>
			<img src="/img1src"/>
			<a href="/some/first/page.html">
				<span>
					<span>text3 aft 
					</span>
				</span>
				<span>text2 aft 
				</span>
			</a>
		</div>
	</div></body></html>`

}

func Test2(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	doc, err := html.Parse(strings.NewReader(testDocs[0]))
	if err != nil {
		lge(err)
		return
	}
	removeCommentsAndIntertagWhitespace(NdX{doc, 0})

	splitAnchSubtreesByImage(doc)

	removeCommentsAndIntertagWhitespace(NdX{doc, 0})
	reIndent(doc, 0)
	var b bytes.Buffer
	err = html.Render(&b, doc)
	lge(err)
	if b.String() != testDocs[1] {
		t.Errorf("output unexpted")
	}

	// osutilpb.Dom2File("xx.html", doc)
	// osutilpb.Bytes2File("xw.html", []byte(testDocs[0]))

	lg("end")

}
