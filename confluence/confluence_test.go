package confluence

import (
	"os"
	"testing"
)

var (
	user = "ci-pink"
	pass = "2bmuR3dC8u43E"
)

func TestConfluence_updatePage(t *testing.T) {
	t.Log(NewConfluence(user, pass).UpdateSngErrors("./test/code1.go"))
}

func TestConfluence_UpdateSngUnitTestCoverage(t *testing.T) {
	t.Log(NewConfluence(user, pass).UpdateSngDevBuild("test1-service", "1.0.0-dev-1", "89.9%", "http://www.baidu.com"))
}

func TestConfluence_createCodeFile(t *testing.T) {
	NewConfluence(user, pass).createCodeFile()
}

func TestParser(t *testing.T) {
	os.Rename("../lcode.proto", "../lcode.proto")
	t.Log(NewConfluence(user, pass).UpdateCacErrors("../lcode.proto"))
}
