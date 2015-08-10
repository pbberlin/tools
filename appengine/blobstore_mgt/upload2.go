package blobstore_mgt

import (
	"html/template"
	"net/http"
	"net/url"

	"appengine"
	"appengine/blobstore"

	"fmt"
	"path"

	"appengine/image"

	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"

	"bytes"
	"strings"
	"time"

	"appengine/datastore"
	"appengine/user"
)

const upload2HTML = `
		<form action="{{.}}" method="POST" enctype="multipart/form-data">
			Upload File: <input type="file" name="post_field_file"><br>
			<input type="text"   name="title"  value="enter Title"><br>
			<input type="text"   name="descr"  value="enter Description" size=60><br>
			<input type="submit" name="submit" value="Submit">
		</form>
`

var upload2 = template.Must(template.New("topLevelTemplateName").Parse(upload2HTML))

const restrictForm = `
		<form 
			method="post"
			aaaction="/blob2"
		    enctype="application/x-www-form-urlencoded" 
		>
			<input type="text"   name="nameFilter"  value="nameFilter"><br>
			<input type="submit" name="submit"    value="Submit">
		</form>
`

const openHTML = `
<html>
	<body style='margin:8px;'>
	<style> * { margin: 4 0 } </style>
	<style>
		.ib {
			vertical-align:middle;
			display:inline-block;
			width:95px;
		}
	</style>
`

const closeHTML = `
	</body>
</html>`

func submitUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/blob2/processing-new-upload", nil)
	loghttp.E(w, r, err, false)

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	err = upload2.Execute(w, uploadURL)
	loghttp.E(w, r, err, false)
}

func processUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.Write(b1.Bytes())
	}()

	blobs, otherFormFields, err := MyParseUpload(r)
	loghttp.E(w, r, err, false, fmt.Errorf("Fehler beim Parsen: %v", err))

	// s1 := stringspb.IndentedDump(blobs)
	// if len(s1) > 2 {
	// 	b1.WriteString("<br>blob: " + s1)
	// }

	// s2 := stringspb.IndentedDump(otherFormFields)
	// if len(s2) > 2 {
	// 	b1.WriteString("<br>otherF: " + s2+"<br>")
	// }

	numFiles := len(blobs["post_field_file"]) // this always yields (int)1
	numOther := len(otherFormFields["post_field_file"])
	if numFiles == 0 && numOther == 0 {
		//http.Redirect(w, r, "/blob2/upload", http.StatusFound)
		b1.WriteString("<a href='/blob2/upload' >No file uploaded? Try again.</a><br>")
		b1.WriteString("<a href='/blob2' >List</a><br>")
		return
	}

	if numFiles > 0 {

		blob0 := blobs["post_field_file"][0]

		dataStoreClone(w, r, blob0, otherFormFields)

		urlSuccessX := "/blob2/serve-full?blobkey=" + string(blob0.BlobKey)
		urlThumb := "/blob2/thumb?blobkey=" + string(blob0.BlobKey)
		b1.WriteString("<a href='/blob2' >List</a><br>")
		b1.WriteString("<a href='" + urlSuccessX + "' >View Full: " + fmt.Sprintf("%v", blob0) + " - view it</a><br>\n")
		b1.WriteString("<a href='" + urlThumb + "' >View Thumbnail</a><br>\n")

	}

	//http.Redirect(w, r, urlSuccess, http.StatusFound)
}

func serveFull(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobkey")))
}

// working differently as in Python
//		//blob2s := blobstore.BlobInfo.gql("ORDER BY creation DESC")
func blobList(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		b1.WriteString(closeHTML)
		w.Write(b1.Bytes())
	}()

	errFormParse := r.ParseForm()
	if errFormParse != nil {
		b1.WriteString(fmt.Sprintf("Form parse Error %v", errFormParse))
		return
	}

	s1 := ""
	b1.WriteString(openHTML)
	b1.WriteString(restrictForm)

	nameFilter := r.PostFormValue("nameFilter")
	nameFilter = strings.TrimSpace(nameFilter)

	if nameFilter == "" {
		return
	} else {
		tn := time.Now()
		f := "2006-01-02 15:04:05"
		f = "January 02"
		prefix := tn.Format(f)
		// nameFilter = fmt.Sprintf("%v %s", prefix,nameFilter )
		if !strings.HasPrefix(nameFilter, prefix) || len(nameFilter) == len(prefix) {
			b1.WriteString(fmt.Sprintf("cannot filter by %v", nameFilter))
			return
		} else {
			nameFilter = nameFilter[len(prefix):]
		}
	}

	u := user.Current(c)
	if u != nil {
		b1.WriteString(fmt.Sprintf("%v %v %v<br>\n", u.ID, u.Email, u.FederatedIdentity))
	} else {
		b1.WriteString("nobody calling on the phone<br>")
	}

	// c := appengine.NewContext(r)
	q := datastore.NewQuery("__BlobInfo__")
	if nameFilter != "" {

		// nameFilter = strings.ToLower(nameFilter)
		b1.WriteString(fmt.Sprintf("Filtering by %v-%v<br>", nameFilter, stringspb.IncrementString(nameFilter)))

		q = datastore.NewQuery("__BlobInfo__").Filter("filename >=", nameFilter)
		q = q.Filter("filename <", stringspb.IncrementString(nameFilter))

	}
	for t := q.Run(c); ; {
		var bi BlobInfo
		dsKey, err := t.Next(&bi)

		if err == datastore.Done {
			// c.Infof("   No Results (any more) blob-list %v", err)
			break
		}
		// other err
		if err != nil {
			loghttp.E(w, r, err, false)
			return
		}

		//s1 = fmt.Sprintf("key %v %v %v %v %v %v<br>\n", dsKey.AppID(),dsKey.Namespace() , dsKey.Parent(), dsKey.Kind(), dsKey.StringID(), dsKey.IntID())
		//b1.WriteString( s1 )

		//s1 = fmt.Sprintf("blobinfo: %v %v<br>\n", bi.Filename, bi.Size)
		//b1.WriteString( s1 )

		ext := path.Ext(bi.Filename)
		base := path.Base(bi.Filename)
		base = base[:len(base)-len(ext)]

		//b1.WriteString( fmt.Sprintf("-%v-  -%v-",base, ext) )

		base = strings.Replace(base, "_", " ", -1)
		base = strings.Title(base)
		ext = strings.ToLower(ext)

		titledFilename := base + ext
		if strings.HasPrefix(titledFilename, "Backup") {
			showBackup := r.FormValue("backup")
			if len(showBackup) < 1 {
				continue
			}
		}

		s1 = fmt.Sprintf("<a class='ib' style='width:280px;margin-right:20px' target='_view' href='/blob2/serve-full?blobkey=%v'>%v</a> &nbsp; &nbsp; \n", dsKey.StringID(), titledFilename)
		b1.WriteString(s1)

		if bi.ContentType == "image/png" || bi.ContentType == "image/jpeg" {
			s1 = fmt.Sprintf("<img class='ib' style='width:40px;' src='/_ah/img/%v%v' />\n",
				dsKey.StringID(), "=s200-c")
			b1.WriteString(s1)

			s1 = fmt.Sprintf("<a class='ib' target='_view' href='/_ah/img/%v%v'>Thumb</a>\n",
				dsKey.StringID(), "=s200-c")
			b1.WriteString(s1)

		} else {
			s1 = fmt.Sprintf("<span class='ib' style='width:145px;'> &nbsp; no thb</span>")
			b1.WriteString(s1)
		}

		s1 = fmt.Sprintf("<a class='ib' target='_rename_delete' href='/blob2/rename-delete?action=delete&blobkey=%v'>Delete</a>\n",
			dsKey.StringID())
		b1.WriteString(s1)

		s1 = fmt.Sprintf(`
			<span class='ib' style='width:450px; border: 1px solid #aaa'>
				<form target='_rename_delete' action='/blob2/rename-delete' >
					<input name='blobkey'  value='%v'     type='hidden'/>
					<input name='action'   value='rename' type='hidden'/>
					<input name='filename' value='%v' size='42'/>
					<input type='submit'   value='Rename' />
				</form>
			</span>
			`, dsKey.StringID(), bi.Filename)
		b1.WriteString(s1)

		b1.WriteString("<br><br>\n\n")

	}

	b1.WriteString("<a accesskey='u' href='/blob2/upload' ><b>U</b>pload</a><br>")
	b1.WriteString(`<a  href='https://appengine.google.com/blobstore/explorer?&app_id=s~libertarian-islands' 
		>Delete via Console</a><br>`)

}

func renameOrDelete(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	b1 := new(bytes.Buffer)
	s1 := ""

	defer func() {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.Write(b1.Bytes())
	}()

	// c := appengine.NewContext(r)

	bk := r.FormValue("blobkey")
	if bk == "" {
		b1.WriteString("No blob key given<br>")
		return
	} else {
		s1 = fmt.Sprintf("Blob key given %q<br>", bk)
		b1.WriteString(s1)
	}

	dsKey := datastore.NewKey(c, "__BlobInfo__", bk, 0, nil)

	q := datastore.NewQuery("__BlobInfo__").Filter("__key__=", dsKey)

	var bi BlobInfo
	var found bool

	for t := q.Run(c); ; {
		_, err := t.Next(&bi)
		if err == datastore.Done {
			c.Infof("   No Results (any more), blob-rename-delete %v", err)
			break
		}
		// other err
		if err != nil {
			logif.E(err)
			return
		}
		found = true
		break
	}

	if found {
		ac := r.FormValue("action")
		if ac == "delete" {
			b1.WriteString("deletion  ")

			// first the binary data
			keyBlob, err := blobstore.BlobKeyForFile(c, bi.Filename)
			logif.E(err)

			if err != nil {
				b1.WriteString(fmt.Sprintf(" ... failed (1) %v", err))
			} else {
				err = blobstore.Delete(c, keyBlob)
				logif.E(err)
				if err != nil {
					b1.WriteString(fmt.Sprintf(" ... failed (2) %v<br>", err))
				} else {
					// now the datastore record
					err = datastore.Delete(c, dsKey)
					logif.E(err)
					if err != nil {
						b1.WriteString(fmt.Sprintf(" ... failed (3) %v<br>%#v<br>", err, dsKey))
					} else {
						b1.WriteString(" ... succeeded<br>")
					}

				}
			}
		}

		if ac == "rename" {
			b1.WriteString("renaming ")

			nfn := r.FormValue("filename")
			if nfn == "" || len(nfn) < 4 {
				b1.WriteString(" ... failed - at LEAST 4 chars required<br>")
				return
			}
			nfn = strings.ToLower(nfn)
			bi.Filename = nfn
			_, err := datastore.Put(c, dsKey, &bi)
			logif.E(err)
			if err != nil {
				b1.WriteString(fmt.Sprintf(" ... failed. %v", err))
			} else {
				b1.WriteString(" ... succeeded<br>")
			}
		}

	} else {
		b1.WriteString("no blob found for given blobkey<br>")
	}
	b1.WriteString("<a href='/blob2'>Back to list</a><br>")

}

func serveThumb(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	// c := appengine.NewContext(r)
	k := appengine.BlobKey(r.FormValue("blobkey"))

	var o image.ServingURLOptions = *new(image.ServingURLOptions)
	o.Size = 200
	o.Crop = true
	url, err := image.ServingURL(c, k, &o)

	loghttp.E(w, r, err, false)

	http.Redirect(w, r, url.String(), http.StatusFound)
}

// This was an attempt, to "catch" the uploaded blob
// and store it by myself into the datastore -
// where I would be able to delete it.
// But this failed - the actual blob data does not even reach the appengine.
// Only the blob-info data.
func dataStoreClone(w http.ResponseWriter, r *http.Request,
	blob0 *BlobInfo, otherFormFields url.Values) {

	return

	c := appengine.NewContext(r)

	wbl := dsu.WrapBlob{}
	wbl.Category = "print"
	wbl.Name = otherFormFields["title"][0] + " - " + otherFormFields["descr"][0]
	wbl.Name += " - " + stringspb.LowerCasedUnderscored(blob0.Filename)
	wbl.Desc = fmt.Sprintf("%v", blob0.BlobKey)
	wbl.S = blob0.ContentType
	if len(otherFormFields["post_field_file"]) > 0 {
		filecontent := otherFormFields["post_field_file"][0]
		wbl.VByte = []byte(filecontent)
	}
	keyX2 := "bl" + time.Now().Format("060102_1504-05")
	_, errDS := dsu.BufPut(c, wbl, keyX2)
	loghttp.E(w, r, errDS, false)

}

func init() {
	http.HandleFunc("/print", loghttp.Adapter(blobList))
	http.HandleFunc("/blob2", loghttp.Adapter(blobList))
	http.HandleFunc("/blob2/upload", loghttp.Adapter(submitUpload))
	http.HandleFunc("/blob2/processing-new-upload", loghttp.Adapter(processUpload))
	http.HandleFunc("/blob2/serve-full", loghttp.Adapter(serveFull))
	http.HandleFunc("/blob2/thumb", loghttp.Adapter(serveThumb))
	http.HandleFunc("/blob2/rename-delete", loghttp.Adapter(renameOrDelete))
}
