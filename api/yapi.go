package api

import (
	"encoding/json"
	"errors"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/client/httpx"
)

type Yapi struct {
	client *httpx.Client
}

func getCppType(goType string) string {
	switch goType {
	case "uint16":
		return "uint16_t"
	case "uint32":
		return "uint32_t"
	case "uint64":
		return "uint64_t"
	case "int16":
		return "int16_t"
	case "int32":
		return "int32_t"
	case "int64":
		return "int64_t"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "const std::string&"
	case "bool":
		return "bool"
	}
	panic(errors.New("unknow cpp type: " + goType))
}

func getCppJsonFunc(goType string) string {
	switch goType {
	case "uint16":
		return "Uint"
	case "uint32":
		return "Uint"
	case "uint64":
		return "Uint"
	case "int16":
		return "Int"
	case "int32":
		return "Int"
	case "int64":
		return "Int"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "String"
	case "bool":
		return "Bool"
	}
	panic(errors.New("unknow cpp type: " + goType))
}

func NewYapi() *Yapi {
	conf := &httpx.HttpConfig{
		Host: "http://yapi.kaiqitech.com",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return &Yapi{
		client: httpx.NewClient(
			httpx.Config(conf),
			httpx.Codec(httpx.NewJsonCodec(conf)),
		),
	}
}

type GetProjectInfoReq struct {
	Token     string `http:"token,query"`
	ProjectId int    `http:"project_id,query"`
	Page      int    `http:"page,query"`
	Limit     int    `http:"limit,query"`
}

type GetProjectInfoRsp struct {
	Errcode int `json:"errcode,omitempty"`
	Data    struct {
		Count int `json:"count,omitempty"`
		Total int `json:"total,omitempty"`
		List  []*struct {
			Id int `json:"_id,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

func (self *Yapi) GetProjectInfo(token string, projectId int) ([]int, *httpx.Serror) {
	getProjectInfoRsp := &GetProjectInfoRsp{}
	if err := self.client.Get("/api/interface/list", &GetProjectInfoReq{
		Token:     token,
		ProjectId: projectId,
		Page:      1,
		Limit:     100,
	}, getProjectInfoRsp); err != nil {
		return nil, err
	}
	ids := []int{}
	for _, listItem := range getProjectInfoRsp.Data.List {
		ids = append(ids, listItem.Id)
	}
	return ids, nil
}

type GetApiInfoReq struct {
	Token string `http:"token,query"`
	Id    int    `http:"id,query"`
}

type GetApiInfoRsp struct {
	Errcode int `json:"errcode,omitempty"`
	Data    struct {
		Title    string   `json:"title,omitempty"`
		Path     string   `json:"path,omitempty"`
		Tag      []string `json:"tag,omitempty"`
		Method   string   `json:"method,omitempty"`
		ReqQuery []struct {
			Name string `json:"name,omitempty"`
			Type string `json:"example,omitempty"`
			Desc string `json:"desc,omitempty"`
		} `json:"req_query,omitempty"`
		ReqHeaders []struct {
			Name string `json:"name,omitempty"`
			Type string `json:"example,omitempty"`
			Desc string `json:"desc,omitempty"`
		} `json:"req_headers,omitempty"`
		ReqBodyOther string `json:"req_body_other,omitempty"`
		ResBody      string `json:"res_body,omitempty"`
		Markdown     string `json:"markdown,omitempty"`
	} `json:"data,omitempty"`
}

func (self *GetApiInfoRsp) ToApiInfo() *ApiInfo {
	if self.Errcode != 0 {
		return nil
	}

	names := strings.Split(self.Data.Title, " ")
	if len(names) < 2 {
		return nil
	}
	names = strings.Split(names[1], ".")
	if len(names) < 2 {
		return nil
	}

	apiInfo := &ApiInfo{
		StructName:   names[0],
		FuncName:     names[1],
		FileName:     stringx.From(names[1]).ToSnake(),
		Path:         self.Data.Path,
		Method:       self.Data.Method,
		MethodFunc:   stringx.From(self.Data.Method).ToCamel(),
		HandlerType:  self.Data.Tag[0],
		QueryFileds:  []*QueryOrHeaderFiled{},
		HeaderFileds: []*QueryOrHeaderFiled{},
		ReqBody:      &ReqBody{},
		RspBody:      &RspBody{},
	}

	for _, query := range self.Data.ReqQuery {
		apiInfo.QueryFileds = append(apiInfo.QueryFileds, &QueryOrHeaderFiled{
			Name:      query.Name,
			FieldName: stringx.From(strings.Replace(query.Name, "-", "_", -1)).ToCamel(),
			Type:      query.Type,
			// CppType:      getCppType(query.Type),
			// CppJsonFunc:  getCppJsonFunc(query.Type),
			CppFiledName: stringx.From(strings.Replace(query.Name, "-", "_", -1)).ToCamel(),
			Comment:      query.Desc,
			Dot:          "`",
		})
	}

	for _, header := range self.Data.ReqHeaders {
		if header.Name == "Content-Type" || header.Name == "x-k7-stage" {
			continue
		}
		apiInfo.HeaderFileds = append(apiInfo.HeaderFileds, &QueryOrHeaderFiled{
			Name:      header.Name,
			FieldName: stringx.From(strings.Replace(header.Name, "-", "_", -1)).ToCamel(),
			Type:      header.Type,
			// CppType:     getCppType(header.Type),
			// CppJsonFunc: getCppJsonFunc(header.Type),
			Comment: header.Desc,
			Dot:     "`",
		})
	}

	if len(self.Data.ReqBodyOther) > 0 {
		if err := json.Unmarshal([]byte(self.Data.ReqBodyOther), apiInfo.ReqBody); err != nil {
			return nil
		}
	}

	// if len(self.Data.ResBody) > 0 {
	// 	if err := json.Unmarshal([]byte(self.Data.ResBody), apiInfo.RspBody); err != nil {
	// 		return nil
	// 	}
	// }

	return apiInfo
}

type QueryOrHeaderFiled struct {
	Name         string
	FieldName    string
	Type         string
	CppType      string
	CppJsonFunc  string
	CppFiledName string
	Comment      string
	Dot          string
}

type ApiInfo struct {
	StructName   string
	FuncName     string
	FileName     string
	Path         string
	Method       string
	MethodFunc   string
	HandlerType  string
	QueryFileds  []*QueryOrHeaderFiled
	HeaderFileds []*QueryOrHeaderFiled
	ReqBody      *ReqBody
	RspBody      *RspBody
}

type ReqBody struct {
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type RspBody struct {
	Properties struct {
		Data struct {
			Properties map[string]interface{} `json:"properties,omitempty"`
		} `json:"data,omitempty"`
	} `json:"properties,omitempty"`
}

func (self *Yapi) GetApiInfo(token string, id int) (*GetApiInfoRsp, *httpx.Serror) {
	getApiInfoRsp := &GetApiInfoRsp{}
	if err := self.client.Get("/api/interface/get", &GetApiInfoReq{
		Token: token,
		Id:    id,
	}, getApiInfoRsp); err != nil {
		return nil, err
	}
	return getApiInfoRsp, nil
}

type UpdateApiInfoReq struct {
	Token string `json:"token"`
	Data  struct {
		Title    string   `json:"title,omitempty"`
		Path     string   `json:"path,omitempty"`
		Tag      []string `json:"tag,omitempty"`
		Method   string   `json:"method,omitempty"`
		ReqQuery []struct {
			Name string `json:"name,omitempty"`
			Type string `json:"example,omitempty"`
			Desc string `json:"desc,omitempty"`
		} `json:"req_query,omitempty"`
		ReqHeaders []struct {
			Name string `json:"name,omitempty"`
			Type string `json:"example,omitempty"`
			Desc string `json:"desc,omitempty"`
		} `json:"req_headers,omitempty"`
		ReqBodyOther string `json:"req_body_other,omitempty"`
		ResBody      string `json:"res_body,omitempty"`
		Markdown     string `json:"markdown,omitempty"`
	} `json:"data,omitempty"`
}

type UpdateApiInfoRsp struct {
	Errcode int `json:"errcode,omitempty"`
}

func (self *Yapi) UpdateApiInfo(token string, updateApiInfoReq *UpdateApiInfoReq) (*UpdateApiInfoRsp, *httpx.Serror) {
	updateApiInfoRsp := &UpdateApiInfoRsp{}
	if err := self.client.Post("/api/interface/up", updateApiInfoReq, updateApiInfoRsp); err != nil {
		return nil, err
	}
	return updateApiInfoRsp, nil
}
