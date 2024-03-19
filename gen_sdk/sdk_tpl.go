package sdk

const sngSdkTpl = `package {{.pkg}}

import (
	"context"

	sngsdkgo "gitlab.kaiqitech.com/k7game/server/supports/sng-sdk-go.git"
)

{{.import}}

{{.req}}

{{.rsp}}

{{$rspLen := len .rsp}}

func {{.func}}(ctx context.Context, req *{{.handle}}{{.func}}Req {{if gt $rspLen 0}}, rsp *{{.handle}}{{.func}}Rsp{{end}}) error {
	return sngsdkgo.Default.Call(ctx, "k7game.server.{{.service}}", "{{.endpointName}}", req {{if gt $rspLen 0}}, rsp{{else}}, nil{{end}})
}
`
