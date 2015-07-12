package backend

// block level root
type B0 struct {
	NColsViewport int // either NCols or NRows determines the other value
	NumXSectors   int
	NumYSectors   int
	VB1           []B1
	NumB1         int
	NumB2         int

	MaxNCols int
	MaxNRows int
}

// block level 1
type B1 struct {
	Heading        string
	VB2            []B2  // editorial sorting - as ordered by the editor
	BySize         ByInt // sorted by Number of Characters
	ByHeading      ByStr // alphabetically sorted
	NCols          int   // either NCols or NRows determines the other value
	NColsEditorial int   // fixed by editor
	NRows          int   // NRows := len(VB2) / NCols ; if len(VB2) % NCols != 0 {NRows++}
	Discrepancy    int   // NRows*NCols - len(VB2)
	StartNewRow    bool

	//OrderedBySize []B2 // obsolete
}

type B2 struct {
	Linktext       string
	Url            string
	Target         string
	Description    string
	IsFirst        bool
	IsLast         bool
	Size           int
	EditorialIndex int // for referral from a sorted helper array
}

var X = []B1{
	B1{
		Heading:        "Diverse",
		NColsEditorial: 4,
		VB2: []B2{
			B2{Linktext: "Login", Url: "/login"},
			B2{Linktext: "Schreib-Methoden", Url: "/write-methods", Description: "FPrintf, bytes.Buffer, io.WriteString and others"},
			B2{Linktext: "Letzte Email", Url: "/email-view"},
			B2{Linktext: "Blob List", Url: "/blob/list"},
			B2{Linktext: "Template Demo 1", Url: "/tpl/demo1"},
			B2{Linktext: "Template Demo 2", Url: "/tpl/demo2"},
			B2{Linktext: "Http fetch", Url: "/fetch-url"},
			B2{Linktext: "Instance Info", Url: "/instance-info/view"},
			B2{Linktext: "Gob encode decode", Url: "/big-query/test-gob-codec"},

			B2{Linktext: "JSON encode", Url: "/json-encode"},
			B2{Linktext: "JSON decode", Url: "/json-decode"},

			B2{Linktext: "Fulltext put", Url: "/fulltext-search/put"},
			B2{Linktext: "Fulltext get", Url: "/fulltext-search/get"},
		},
	},

	B1{
		Heading:     "Guestbook",
		StartNewRow: false,
		VB2: []B2{
			B2{Linktext: "Neuer Eintrag", Url: "/guest-entry"},
			B2{Linktext: "Einträge auflisten", Url: "/guest-view"},
			B2{Linktext: "Einträge auflisten - paged - serialized cursor", Url: "/guest-view-cursor"},
		},
	},
	B1{
		Heading: "Drawing",
		VB2: []B2{
			B2{Linktext: "Drawing a static chart", Url: "/image/draw-lines-example"},
		},
	},
	B1{
		Heading:     "Big Query",
		StartNewRow: true,
		VB2: []B2{
			B2{Linktext: "Get real data", Url: "/big-query/query-into-datastore"},
			B2{Linktext: "Get mocked data", Url: "/big-query/mock-data-into-datastore"},
		},
	},
	B1{
		Heading:     "... with Chart",
		StartNewRow: false,
		VB2: []B2{
			B2{Linktext: "Process Data 1 (mock=1},", Url: "/big-query/regroup-data-01?mock=0"},
			B2{Linktext: "Process Data 2", Url: "/big-query/regroup-data-02?f=table"},
			B2{Linktext: "Show as Table", Url: "/big-query/show-table"},
			B2{Linktext: "Show as Chart", Url: "/big-query/show-chart"},
			B2{Linktext: "As HTML", Url: "/big-query/html"},
		},
	},
	B1{
		Heading: "Request Images",
		VB2: []B2{
			B2{Linktext: "WrapBlob from Datastore", Url: "/image/img-from-datastore?p=chart1"},
			B2{Linktext: "base64 from Datastore", Url: "/image/base64-from-datastore?p=chart1"},
			B2{Linktext: "base64 from Variable", Url: "/image/base64-from-var?p=1"},
			B2{Linktext: "base64 from File", Url: "/image/base64-from-file?p=static/pberg1.png"},
		},
	},
	B1{
		Heading: "Namespaces and Task Queues",
		VB2: []B2{
			B2{Linktext: "Increment", Url: "/namespaced-counters/increment"},
			B2{Linktext: "Read", Url: "/namespaced-counters/read"},
			B2{Linktext: "Push to task-queue", Url: "/namespaced-counters/queue-push"},
		},
	},
	B1{
		Heading: "URLs with/without Ancestors",
		VB2: []B2{
			B2{Linktext: "Add without ancestor", Url: "/save-url/save-no-anc"},
			B2{Linktext: "Add <b>with</b> ancestor", Url: "/save-url/save-wi-anc"},
			B2{Linktext: "Query without ancestor", Url: "/save-url/view-no-anc"},
			B2{Linktext: "Query <b>with</b> ancestor", Url: "/save-url/view-wi-anc"},
			B2{Linktext: "Query URL Backend", Url: "/save-url/backend"},
		},
	},
	B1{
		Heading: "x",
		VB2:     []B2{},
	},
}

func init() {

}
