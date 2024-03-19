package version

var (
	Version   = "1.0.4"
	ChangeLog = map[string]string{
		"1.0.0": `
			初始版本
		`,
		"1.0.2": `
			支持changelog上传到confluence
		`,
		"1.0.3": `
			changelog支持下载链接参数
		`,
		"1.0.4": `
			changelog支持解析yaml配置文件内容
		`,
	}
)
