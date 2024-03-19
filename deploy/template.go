package deploy

var deployScript = `
TmpPkgPath=/tmp/{{.name}}.tar.gz
ExtractFiles=""
{{if eq .env "dev"}}
WorkDir=/data/apps/Server/DEVServer/{{.name}}
{{else}}
WorkDir=/data/apps/Server/Server/{{.name}}
{{end}}
echo "================================================="
echo "  发布环境: {{.env}}"
echo "  发布服务: {{.name}}"
echo "  工作目录: $WorkDir"
echo "    服务名: {{.env}}.{{.lower_name}}.service"
echo "   取包URL: {{.url}}"
echo "本地包路径: {{.file}}"
echo "================================================="

log()
{
        {{if .verbose}}echo $1{{end}}
        return 0
}

downloadPkg()
{
        log "正在下载安装包到 $TmpPkgPath ..."
        wget {{.url}} -O $TmpPkgPath -o /dev/null
        if [ $? != 0 ]; then
                log "下载安装包失败"
                return 1
        fi
}

preClean()
{
        sudo rm -rf "$WorkDir/{{.name}} " "$WorkDir/{{.name}}.pdb" Build.properties Changelog
}

postClean()
{
        sudo rm -rf $TmpPkgPath
}

getExtractableFiles()
{
        pkg=$1
        log $pkg
        Files=$(tar -tf $pkg {{.name}} Build.properties ChangeLog config/ServerString.json 2>/dev/null)
        log $Files | grep -q "{{.name}}"
        if [ $? -eq 0 ]; then
                ExtractFiles=$(echo $Files | tr -s "\n" " ")
        else
                log "安装包中缺少必要文件{{.name}}"
        fi
        return $?
}

deploy()
{
        PreDir=$(pwd)
        log "正在部署 {{.name}} ..."
        cd $WorkDir
        if [ $? != 0 ]; then
                log "目录 [$WorkDir] 不存在"
                return 1
        fi
        log "正在从安装包 $TmpPkgPath 中获取文件列表 ..."
        getExtractableFiles $TmpPkgPath
        if [ $? != 0 ]; then
                return 1
        fi
        preClean $WorkDir
        log "从安装包 $TmpPkgPath 中提取文件 $ExtractFiles ..."
        tar xf $TmpPkgPath $ExtractFiles
        if [ $? != 0 ]; then
                log "解压 $TmpPkgPath 失败"
                return 1
        fi
        log "修改权限 {{.name}}"
        chmod ugoa+rwx {{.name}}
        if [ $? != 0 ]; then
                log "文件 {{.name}} 权限修改失败"
                return 1
        fi
        log "正在重启服务 {{.env}}.{{.lower_name}}.service ..."
        sudo systemctl restart {{.env}}.{{.lower_name}}.service
        if [ $? != 0 ]; then
                log "服务 {{.env}}.{{.lower_name}}.service 重启失败"
                return 1
        fi
        log "服务 {{.env}}.{{.lower_name}}.service 发布完成"
        systemctl status {{.env}}.{{.lower_name}}.service
        cd $PreDir
        return 0
}
{{if gt (len .url) 0}}
downloadPkg
{{else}}
mv {{.file}} $TmpPkgPath
{{end}}
if [ "$?" -ne 0 ]; then
        exit 0
else
        deploy
fi

postClean $WorkDir
exit 0
`
