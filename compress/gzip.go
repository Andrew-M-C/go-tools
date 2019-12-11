package compress

import (
	"bytes"
	"compress/gzip"
	"strings"
	"time"
)

func gzipCompress(in []byte, opt Options) (out []byte, err error) {
	if nil == in {
		return nil, ErrNilInput
	}
	outBuff := bytes.Buffer{}
	w := gzip.NewWriter(&outBuff)
	defer w.Close()

	w.Name = opt.Name
	w.Comment = opt.Comment
	w.ModTime = time.Now().Local()

	w.Write(in)
	w.Flush()
	return outBuff.Bytes(), nil
}

func gzipUncompress(in []byte, opt Options) (out []byte, err error) {
	if nil == in {
		return nil, ErrNilInput
	}

	var n int
	inBuff := bytes.NewBuffer(in)
	r, err := gzip.NewReader(inBuff)
	if err != nil {
		return
	}
	defer r.Close()

	if opt.MaxRead > 0 {
		outBuff := make([]byte, opt.MaxRead)
		n, err = r.Read(outBuff)
		if err != nil {
			return
		}
		return outBuff[:n], nil
	}

	out = []byte{}
	outBuff := make([]byte, 2*len(in))
	for {
		n, err = r.Read(outBuff)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				err = nil
			}
			return
		}
		if 0 == n {
			return
		}
		out = append(out, outBuff[:n]...)
	}
	return
}
