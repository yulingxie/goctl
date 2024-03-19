package sng

import (
	"testing"
)

func TestOpenProject(t *testing.T) {
	openProject("money-test")
}

func TestCreateSngServiceProject(t *testing.T) {
	gen, _ := NewGenerator("./cfg.yaml")
	gen.GenSngServiceAndTest()
}
