package domclean2

import (
	"bytes"
	"net/url"
	"path/filepath"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/os/osutilpb"
	"golang.org/x/net/html"
)

type CleaningOptions struct {
	FNamer func() string

	Proxify    bool
	ProxyHost  string
	RemoteHost string

	AddOutline bool
	AddID      bool

	Beautify bool // make pretty at the end, removes <a> linktext trailing space
}

func FileNamer(logdir string, fileNumber int) func() string {
	cntr := -2
	return func() string {
		cntr++
		if cntr == -1 {
			return spf("outp_%03v", fileNumber) // prefix/filekey
		} else {
			fn := spf("outp_%03v_%v", fileNumber, cntr) // filename with stage
			fn = filepath.Join(logdir, fn)
			return fn
		}
	}
}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}

func fileDump(doc *html.Node, fNamer func() string) {
	if fNamer != nil {
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
	}
}

func DomClean(b []byte, opt CleaningOptions) (*html.Node, error) {

	lg, lge := loghttp.Logger(nil, nil)
	_ = lg

	b = globFixes(b)
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		lge(err)
		return nil, err
	}

	if opt.FNamer != nil {
		osutilpb.Dom2File(opt.FNamer()+".html", doc)
	}

	//
	//
	cleanseDom(doc, 0)
	removeCommentsAndIntertagWhitespace(NdX{doc, 0})
	fileDump(doc, opt.FNamer)

	//
	//
	condenseTopDown(doc, 0, 0)
	removeEmptyNodes(doc, 0)
	fileDump(doc, opt.FNamer)

	//
	//
	removeCommentsAndIntertagWhitespace(NdX{doc, 0}) // prevent spacey textnodes around singl child images
	breakoutImagesFromAnchorTrees(doc)
	fileDump(doc, opt.FNamer)

	//
	//
	condenseBottomUpV3(doc, 0, 7, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 6, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 5, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 4, map[string]bool{"div": true})
	condenseTopDown(doc, 0, 0)

	removeEmptyNodes(doc, 0)
	removeEmptyNodes(doc, 0)

	fileDump(doc, opt.FNamer)

	//
	//
	if opt.Proxify {
		if opt.ProxyHost == "" {
			opt.ProxyHost = routes.AppHost()
		}

		proxify(doc, opt.ProxyHost, &url.URL{Scheme: "http", Host: opt.RemoteHost})
		fileDump(doc, opt.FNamer)
	}

	if opt.Beautify {
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
		reIndent(doc, 0)

	}

	//
	//
	if opt.AddOutline {
		addOutlineAttr(doc, 0, []int{0})
	}
	if opt.AddID {
		addIdAttr(doc, 0, 1)
	}
	if opt.AddOutline || opt.AddID {
		fileDump(doc, opt.FNamer)
	}

	//
	computeXPathStack(doc, 0)
	if opt.FNamer != nil {
		osutilpb.Bytes2File(opt.FNamer()+".txt", xPathDump)
	}

	return doc, nil

}

func DomFormat(doc *html.Node) {
	removeEmptyNodes(doc, 0)
	removeCommentsAndIntertagWhitespace(NdX{doc, 0})
	reIndent(doc, 0)
}
