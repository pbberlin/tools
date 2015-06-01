package conv


import (
	"bytes"
	"github.com/pbberlin/tools/util_err"
	"github.com/pbberlin/tools/util"
	
	"io"
	"strings"

)


func String_to_VVByte( base64_img string )( [][]byte, *bytes.Buffer){

	bMsg := new(bytes.Buffer)

	const chunksize = 400		// 
	
	var size_o int
	if len(base64_img) % chunksize == 0 {
		size_o = len(base64_img) / chunksize
	} else {
		size_o = len(base64_img) / chunksize + 1
	}
	

	VVByte := make([][]byte,size_o)


	cntr := -1
	b := make([]byte,chunksize)
	rdr := strings.NewReader( base64_img )
	for {
		cntr++
		n,err := rdr.Read(b)
		if err == io.EOF {
			break
		}
		util_err.Err_log(err)
		if n < 1 {
			break
		}

		indep_copy := make([]byte,n)			
		copy( indep_copy, b )
		VVByte[cntr] = indep_copy  
		
		bMsg.WriteString( "reading " +  util.Itos(n) + " bytes:\n"  )
		//bMsg.Write(  VVByte[util.Itos(cntr)]  )
	}
	
	return VVByte, bMsg

}

// based on bytes.Buffer and Writing into it
func VVByte_to_string( m [][]byte)( *bytes.Buffer, *bytes.Buffer){

	bRet := new(bytes.Buffer)
	bMsg := new(bytes.Buffer)

	//for i,v := range m {
	for i := 0; i < len(m); i++ {
		n,err := bRet.Write( m[i] )
		util_err.Err_log(err)
		bMsg.WriteString( " lp"+util.Itos(i)+": writing "+util.Itos(n)+" bytes: \n" )
	}
	return bRet,bMsg
}

