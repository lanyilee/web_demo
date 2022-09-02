package compress

import "testing"

func TestCompressGZ(t *testing.T) {
	err := CompressGZ("../../../hack", "test.tgz")
	if err != nil {
		t.Error(err)
	}
}

func TestDeCompressGZ(t *testing.T) {
	err := DeCompressGZ("../../../bin/files/3.tgz", "./data")
	if err != nil {
		t.Error(err)
	}
}
