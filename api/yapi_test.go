package api

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/pretty"
)

func TestYapi_GetProjectInfo(t *testing.T) {
	yapi := NewYapi()
	ids, err := yapi.GetProjectInfo("ea18b184ba39801cf29721dfb445d14e241d8a49fa6fe505daf0fa7af00d05ef", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}

func TestYapi_GetApiInfo(t *testing.T) {
	yapi := NewYapi()
	getApiInfo, err := yapi.GetApiInfo("201fa3a26d5b79b65057e90eed4d9d94618bd561c3201da137274fcdff1dd92d", 13704)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(getApiInfo)
}

func TestYapi_UpdateApiInfo(t *testing.T) {
	yapi := NewYapi()
	getApiInfo, err := yapi.GetApiInfo("201fa3a26d5b79b65057e90eed4d9d94618bd561c3201da137274fcdff1dd92d", 13704)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(getApiInfo)
	t.Log(string(pretty.Pretty(data)))
}
