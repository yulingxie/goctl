package mongo

import "testing"

func Test_generateModel(t *testing.T) {
	gen, err := NewGenerator("./test", "common", "test", "nitro_test,goctl_test")
	if err != nil {
		t.Error(err)
	}

	if err := gen.GenModel(); err != nil {
		t.Error(err)
	}
}
