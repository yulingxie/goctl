package sdk

import (
	"testing"
)

func TestGenerator_genSngSdkGo(t *testing.T) {
	if err := NewGenerator(".").genSngSdkGo("user-service"); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestGenerator_GenAllSngServiceSdk(t *testing.T) {
	NewGenerator("/Users/zwf/gitlab/k7game/server/supports/sng-sdk-go").GenAllSngServiceSdk()
}
