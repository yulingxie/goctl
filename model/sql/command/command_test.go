package command

import "testing"

func Test_fromDB(t *testing.T) {
	if err := fromDB("", "common", "yygsubcasualgame", "game_archive", "sqlm", false); err != nil {
		t.Error(err)
	}
}
