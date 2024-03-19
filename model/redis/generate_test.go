package redis

import "testing"

func Test_generateModel(t *testing.T) {
	gen, err := NewGenerator("./test", "common", "Common,NitroTest")
	if err != nil {
		t.Error(err)
	}

	if err := gen.GenModel(); err != nil {
		t.Error(err)
	}
}
