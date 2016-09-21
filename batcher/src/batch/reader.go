package batch

import (
	"io"
	"os"
	"strings"
)

type ReadHandler struct {
	dbpath   string
	hnd      *os.File
	FileSize int64
}

const BATCH_SIZE = 256000

/* Removes characters after the last new line ("\n") */
func removeAfterLastNL(dataAsString string) (string, int64) {
	lastIndex := strings.LastIndex(dataAsString, "\n")
	origLen := len(dataAsString)
	toCut := origLen - lastIndex
	dataAsString = dataAsString[0 : origLen-toCut]
	return dataAsString, int64(toCut)
}

/* Rewind the file to the last new line ("\n") */
func rewindUntilEOL(p *ReadHandler, toCut int64) {
	_, _ = p.hnd.Seek(-toCut+1, os.SEEK_CUR)
}

func NewHandler(dbpath string) *ReadHandler {
	p := new(ReadHandler)
	p.dbpath = dbpath

	f, err := os.Open(p.dbpath)
	if err != nil {
		panic(err)
	}

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	p.FileSize = fi.Size()
	p.hnd = f

	return p
}

func ReadNext(p *ReadHandler) string {

	data := make([]byte, BATCH_SIZE)
	n, err := p.hnd.Read(data)
	if err != nil {
		if err == io.EOF {
			panic("EOF will be handled")
		}
		panic(err)
	}
	dataAsString := string(data[:n])
	dataAsString, toCut := removeAfterLastNL(dataAsString)
	rewindUntilEOL(p, toCut)
	return dataAsString
}

func Reset(p *ReadHandler) {
	_, err := p.hnd.Seek(0, 0)

	if err != nil {
		panic(err)
	}
}

func Terminate(p *ReadHandler) {
	p.hnd.Close()
}
