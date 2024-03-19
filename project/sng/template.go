package sng

const vscodeTemplate = `
{
     // 使用 IntelliSense 了解相关属性。 
     // 悬停以查看现有属性的描述。
     // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
     "version": "0.2.0",
     "configurations": [
         {
             "name": "Launch",
             "type": "go",
             "request": "launch",
             "mode": "auto",
             "program": "${workspaceFolder}/main.go",
             "env": {},
             "args": [
                "--nacos-v2address=dev-nacos-v2.k7.cn:80",
                "--nacos-address=dev-nacos-web.k7.cn:80",
                "--nacos-namespace=SERVER-DEV",
                "--password=password",
                "--stage=local",
                "--log-level=debug",
                "--port=:8000",
                "--mgmt-port=:8080",
                "--pprof-port=:6060",
                "--trace=true"
             ]
         }
     ]
}
`

const gitignorTemplate = `
cache
log
cover.*
`

const changeLogTemplate = `
# 1.0.0
feat: 初始版本
`

const dockerfileTemplate = `
FROM alpine:3.14
LABEL app="sng"
LABEL appname="{{.serviceName}}"

COPY  {{.serviceName}}  /usr/local/bin/{{.serviceName}}

COPY docker-entrypoint.sh /usr/local/bin/
RUN  chmod a+rx  /usr/local/bin/docker-entrypoint.sh && chmod a+rx /usr/local/bin/{{.serviceName}}
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
     apk add --no-cache tini &&\
     apk add --no-cache curl &&\
     apk add --no-cache tzdata &&\
     apk add --no-cache tcpdump &&\
     cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
     echo "Asia/Shanghai" > /etc/timezone &&\
     rm -rf /var/cache/apk/*

ENTRYPOINT ["/sbin/tini", "--", "docker-entrypoint.sh"]
`

const dockerEntrypointShTpl = `#!/bin/sh
if [[ $TZ ]]; then
  cp /usr/share/zoneinfo/${TZ} /etc/localtime
  echo ${TZ} > /etc/timezone
fi
exec "$@"  
`

const goModTemplate = `
module gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git

go 1.17

require (
	gitlab.kaiqitech.com/k7game/server/components/seraph.git dev
	gitlab.kaiqitech.com/nitro/nitro/v3 dev
)
`

const mainTemplate = `
package main

import (
	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/app"
	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/app/sng/service"
	"gitlab.kaiqitech.com/k7game/server/components/seraph.git/config"
	"gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git/logic"
)

// 此处变量会用于jenkens解析服务名与版本号，开发时不可修改变量名或将其删除
const (
	APP_NAME    = "{{.serviceName}}"
	APP_VERSION = "1.0.0"
)

func main() {
	conf := config.Default(APP_NAME, APP_VERSION)
	service.Start(
		logic.New{{.upperName}}(),
		logic.Endpoints,
		app.Name(APP_NAME),
		app.Version(APP_VERSION),
		app.Config(conf),
	)
}
`

const sngApiTestGoMod = `
module gitlab.kaiqitech.com/k7game/server/test/{{.testName}}.git

go 1.17

require (
	github.com/stretchr/testify v1.7.0
	gitlab.kaiqitech.com/k7game/server/components/seraph.git dev
	gitlab.kaiqitech.com/k7game/server/services/{{.serviceName}}.git dev
	gitlab.kaiqitech.com/nitro/nitro/v3 dev
)
`
