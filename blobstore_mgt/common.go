package blobstore_mgt

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pbberlin/tools/pbstrings"

	"appengine"
)

type BlobInfo struct {
	BlobKey      appengine.BlobKey
	ContentType  string    `datastore:"content_type"`
	CreationTime time.Time `datastore:"creation"`
	Filename     string    `datastore:"filename"`
	MD5          string    `datastore:"md5_hash"`
	Size         int64     `datastore:"size"`
	Upload_id    string    `datastore:"upload_id"`
	ObjectName   string    `datastore:"gs_object_name"`
}

func (bi BlobInfo) String() string {

	b1 := new(bytes.Buffer)
	b1.WriteString("FN: " + pbstrings.LowerCasedUnderscored(bi.Filename))
	b1.WriteString(" Type: " + bi.ContentType)
	b1.WriteString(" " + fmt.Sprintf("%v", bi.Size/1024) + "KB")
	b1.WriteString(" BlobKey:" + fmt.Sprintf("%v", bi.BlobKey))

	return b1.String()

}
