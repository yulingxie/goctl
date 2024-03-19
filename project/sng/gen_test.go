package sng

import (
	"testing"
)

func TestGenerator_GenService(t *testing.T) {
	g, _ := NewGenerator("./cfg.yaml")
	if err := g.GenSngServiceAndTest(); err != nil {
		t.Error(err)
	}
}
