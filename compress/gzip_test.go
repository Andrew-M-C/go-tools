package compress

import (
	"encoding/hex"
	"testing"
)

func TestGzip(t *testing.T) {
	raw := []byte(`Hello, GZip! `)

	comp, err := Compress(Gzip, raw)
	if err != nil {
		t.Errorf("compress gzip error: %v", err)
		return
	}

	t.Logf("compressed:\n%s", hex.Dump(comp))

	ucmp, err := Uncompress(Gzip, comp)
	if err != nil {
		t.Errorf("uncompress gzip error: %v", err)
		return
	}

	t.Logf("uncompressed:\n%s", hex.Dump(ucmp))
	if string(ucmp) != string(raw) {
		t.Errorf("uncompress error")
		return
	}
	return
}
