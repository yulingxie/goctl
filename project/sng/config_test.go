package sng

import (
	"testing"
)

func TestLoadSngServiceProjectConfigFromFile(t *testing.T) {
	cfg, err := LoadSngServiceProjectConfigFromFile("./cfg.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", cfg)
}
