package api

const apiTpl = `
package logic

import (
	"gitlab.kaiqitech.com/nitro/nitro/v3/api"
)

var Endpoints = []api.Endpoint{
  {{- range $apiInfo := .apiInfos }}
	{
		Name:    "{{$apiInfo.StructName}}.{{$apiInfo.FuncName}}",
		Path:    []string{"{{$apiInfo.Path}}"},
		Method:  []string{"{{$apiInfo.Method}}"},
		Handler: "{{$apiInfo.HandlerType}}",
	},
  {{- end }}
}
`

const reqTpl = `
package apim

type {{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req struct {
  {{- range $header := .apiInfo.HeaderFileds }}
    {{$header.FieldName}} {{$header.Type}} {{ .Dot }}http:"{{$header.Name}},header"{{ .Dot }}
  {{- end }}
  {{- range $query := .apiInfo.QueryFileds }}
    {{$query.FieldName}} {{$query.Type}} {{ .Dot }}http:"{{$query.Name}},query"{{ .Dot }}
  {{- end }}
  {{.reqBody}}
}
`

const rspTpl = `
package apim

type {{.apiInfo.StructName}}{{.apiInfo.FuncName}}Rsp struct {
	{{.rspBody}}
}
`

const filedTpl = `{{.camelName}} {{.type}} {{.dot}}json:"{{.name}},omitempty"{{.dot}}`

const structTpl = `
package logic

import (
	"context"

	"gitlab.kaiqitech.com/nitro/nitro/v3/merrors"
	"gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git/models/apim"
)

type {{.structName}} struct {
}

func New{{.structName}}() *{{.structName}} {
	return &{{.structName}}{}
}

{{ $structName := .structName}}
{{- range $funcName := .funcNames }}

func (self *{{$structName}}) {{$funcName}}(ctx context.Context, req *apim.{{$structName}}{{$funcName}}Req) *merrors.ServiceError {
	return merrors.Normal(nil)
}
{{- end }}
`

const apiTestInitTpl = `
package api

import (
	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/config"
	"gitlab.kaiqitech.com/nitro/nitro/v3/httpx"
)

var (
	sapi *httpx.Client
	err  error
)

const (
	USER_ID       uint32 = 1031560
	ERROR_USER_ID uint32 = 1
)

func init() {
	config.NewDevConfig("{{.serviceName}}")
	sapi, err = httpx.GetClient("sapi")
	if err != nil {
		panic(err)
	}
}
`

const apiTestRpcTpl = `
package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git/models/apim"
)

func Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}}(t *testing.T) {
	type DBInit struct{
	}

	tests := []struct {
		name string
		dbInit         *DBInit
		req  *apim.{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req
		want           *apim.{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Rsp
		wantErr        bool
		wantErrCode    int32
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 数据库初始化
			if tt.dbInit != nil {
			}
			// 发起请求
			got := &apim.{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Rsp{}
			err := sapi.{{.apiInfo.MethodFunc}}("{{.apiInfo.Path}}", tt.req, got)
			// 结果比对
			if (err != nil) != tt.wantErr {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (err.Code != tt.wantErrCode) {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErrCode %v", err, tt.wantErrCode)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
`

const apiTestRpcNoRspTpl = `
package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git/models/apim"
)

func Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}}(t *testing.T) {
	type DBInit struct{
	}

	tests := []struct {
		name string
		dbInit         *DBInit
		req  *apim.{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req
		wantErr        bool
		wantErrCode    int32
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 数据库初始化
			if tt.dbInit != nil {
			}
			// 发起请求
			err := sapi.{{.apiInfo.MethodFunc}}("{{.apiInfo.Path}}", tt.req, nil)
			t.Logf("api result err: %v, rsp: %v", err, nil)
			// 结果比对
			if (err != nil) != tt.wantErr {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (err.Code != tt.wantErrCode) {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErrCode %v", err, tt.wantErrCode)
				return
			}
			if !tt.wantErr {
			}
		})
	}
}
`

const apiTestNotifyAndQueueTpl = `
package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git/models/apim"
)

func Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}}(t *testing.T) {
	type DBInit struct{
	}

	tests := []struct {
		name string
		dbInit         *DBInit
		req  *apim.{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req
		wantErr        bool
		wantErrCode    int32
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 数据库初始化
			if tt.dbInit != nil {
			}
			// 发起请求
			err := sapi.{{.apiInfo.MethodFunc}}("{{.apiInfo.Path}}", tt.req, nil)
			t.Logf("api result err: %v, rsp: %v", err, nil)
			// 结果比对
			if (err != nil) != tt.wantErr {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (err.Code != tt.wantErrCode) {
				t.Errorf("Test{{.apiInfo.StructName}}{{.apiInfo.FuncName}} error = %v, wantErrCode %v", err, tt.wantErrCode)
				return
			}
			if !tt.wantErr {
				// 消息接口的结果比对通常在等待几秒后比对数据库
				<-time.After(time.Second * 1)
			}
		})
	}
}
`

const cppReqTpl = `
#ifndef __HTTP_REQUEST_H__
#define __HTTP_REQUEST_H__

#include "request.h"
#include "sapi_define.h"

namespace sapi {
    struct H{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req: public IHttpRequest
    {
        explicit H{{.apiInfo.StructName}}{{.apiInfo.FuncName}}Req(
			uint32_t userId, 
			uint32_t gameId,
		):userId(userId), 
			gameId(gameId),
			IHttpRequest(web::http::methods::POST){
			}

		const web::uri FormatedPath() const
		{
			return web::uri_builder(UserSlotURI)
			.append_query(U("user_id"), userId)
			.append_query(U("game_id"), gameId).to_uri();
		}

	{{- range $header := .apiInfo.HeaderFileds }}
    	{{$header.FieldName}} {{$header.Type}} {{ .Dot }}http:"{{$header.Name}},header"{{ .Dot }}
  	{{- end }}
  	{{- range $query := .apiInfo.QueryFileds }}
    	{{$query.CppType}} {{$query.FieldName}};
  	{{- end }}

    };
}
`

const cppRspTpl = `
`

const benchmarkRunAllTpl = `
#!/bin/sh
name={{.serviceName}}
version=1.1.0

rm -Rf bin
mkdir bin

echo "
<div>
<textarea rows=4 cols=40 readonly style=\"resize: none; border: none;\">
${name}: ${version}
{{.dot}}date{{.dot}}
</textarea>
</div>
" >> bin/report.html

RunTestCase(){
    ./${1}.sh 0
    echo "
<div>
<a href=\"./${1}.html\">${1}</a>
</div>" >> bin/report.html
}
{{- range $apiInfo := .apiInfos }}
RunTestCase {{.FileName}}
{{- end }}
`

const benchmarkGetTpl = `
#!/bin/sh
if [ ! -d "./bin" ]; then
    mkdir ./bin
fi
cd bin
name={{.apiInfo.FileName}}
rm ${name}.*.html ${name}.*.bin ${name}.*.txt ${name}.html

Benchmark(){
    echo "{{.apiInfo.Method}} http://test-papi-gw.svc.qipai007cs.com{{.apiInfo.Path}}" \
        | vegeta attack -name=${1}qps -rate=$1 -duration=5s -timeout=5s > ${name}.${1}.bin
    vegeta plot --title ${name}.${1}.qps ${name}.${1}.bin > ${name}.${1}.html
	vegeta report  -type=text ${name}.${1}.bin > ${name}.${1}.txt
    echo "<div>
<iframe width=\"100%\" height=\"100%\" src=\"${name}.${1}.html\"  frameborder=\"no\"></iframe>
<iframe width=\"100%\" src=\"${name}.${1}.txt\" seamless></iframe>
</div>" >> ${name}.html
}

# 此处设定要测试的qps
Benchmark 100

rm ${name}.*.bin
cd ../
`

const benchmarkPostTpl = `
#!/bin/sh
if [ ! -d "./bin" ]; then
    mkdir ./bin
fi
cd bin
name={{.apiInfo.FileName}}
rm ${name}.*.html ${name}.*.bin ${name}.*.txt ${name}.html
url=$(cd $(dirname $0); pwd)
workspace=${url:0:${#url}-4}

Benchmark(){
    echo "{{.apiInfo.Method}} http://test-papi-gw.svc.qipai007cs.com{{.apiInfo.Path}}" \
        | vegeta attack -name=${1}qps -rate=$1 -duration=5s -timeout=5s -body=${workspace}/${name}.txt > ${name}.${1}.bin
    vegeta plot --title ${name}.${1}.qps ${name}.${1}.bin > ${name}.${1}.html
	vegeta report  -type=text ${name}.${1}.bin > ${name}.${1}.txt
    echo "<div>
<iframe width=\"100%\" height=\"100%\" src=\"${name}.${1}.html\"  frameborder=\"no\"></iframe>
<iframe width=\"100%\" src=\"${name}.${1}.txt\" seamless></iframe>
</div>" >> ${name}.html
}

# 此处设定要测试的qps
Benchmark 100

rm ${name}.*.bin
cd ../
`
