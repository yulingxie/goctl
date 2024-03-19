package git

import (
	"testing"
)

func TestGit_CloneProject(t *testing.T) {
	NewGit("").CloneProject([]string{""}, []int{1485})
}

func TestGit_UpdateAllSngService(t *testing.T) {
	if err := NewGit("").UpdateAllSngService(); err != nil {
		t.Error(err.Error())
	}
}

func TestGit_UpdateAllSngGateway(t *testing.T) {
	if err := NewGit("").UpdateAllSngGateway(); err != nil {
		t.Error(err.Error())
	}
}
