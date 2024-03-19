package chart

import (
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/chart/gw_tpls"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/gen"
)

type Generator struct {
	*gen.Generator
	serviceName string
	foldName    string
}

func NewGenerator(serviceName, dir, foldName string) (*Generator, error) {
	generator, err := gen.NewGenerator(dir)
	if err != nil {
		return nil, err
	}
	if len(foldName) == 0 {
		// 如果未指定folaName，则foldName设置为serviceName
		foldName = serviceName
	}
	return &Generator{serviceName: serviceName, Generator: generator, foldName: foldName}, nil
}

func (self *Generator) GenSngServiceCharts() error {
	serviceNames := strings.Split(self.serviceName, ",")
	for _, name := range serviceNames {
		data := map[string]interface{}{
			"name": name,
		}
		foldName := self.foldName
		// 如果指定了多个服务，则文件夹必须使用服务名以作区分
		if len(serviceNames) > 1 {
			foldName = name
		}
		if err := self.GenItems([]*gen.Item{
			{
				Name: foldName,
				Type: gen.Dir,
			},
			{
				Name:     foldName + "/Chart.yaml",
				Type:     gen.File,
				Template: ChartYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values.yaml",
				Type:     gen.File,
				Template: valuesYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-prev.yaml",
				Type:     gen.File,
				Template: valuesPrevYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-test.yaml",
				Type:     gen.File,
				Template: valuesTestYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-dev.yaml",
				Type:     gen.File,
				Template: valuesDevYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-yeetown.yaml",
				Type:     gen.File,
				Template: valuesYeetownYamlTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/.helmignore",
				Type:     gen.File,
				Template: helmignoreTpl,
			},
			{
				Name: foldName + "/templates",
				Type: gen.Dir,
			},
			{
				Name:        foldName + "/templates/_helpers.tpl",
				Type:        gen.File,
				Template:    helpersTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/deployment.yaml",
				Type:        gen.File,
				Template:    deploymentYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/hpa.yaml",
				Type:        gen.File,
				Template:    hpaYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/ingress.yaml",
				Type:        gen.File,
				Template:    ingressYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/NOTES.txt",
				Type:        gen.File,
				Template:    notesTxtTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/service.yaml",
				Type:        gen.File,
				Template:    serviceYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/serviceaccount.yaml",
				Type:        gen.File,
				Template:    serviceAccountYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/serviceMonitor.yaml",
				Type:        gen.File,
				Template:    serviceMonitorYamlTpl,
				NotParseTpl: true,
			},
		}); err != nil {
			return err
		}
	}
	self.Info("生成完毕")
	return nil
}

func (self *Generator) GenSngGatewayCharts() error {
	serviceNames := strings.Split(self.serviceName, ",")
	for _, name := range serviceNames {
		data := map[string]interface{}{
			"name": name,
		}
		foldName := self.foldName
		// 如果指定了多个服务，则文件夹必须使用服务名以作区分
		if len(serviceNames) > 1 {
			foldName = name
		}
		if err := self.GenItems([]*gen.Item{
			{
				Name: foldName,
				Type: gen.Dir,
			},
			{
				Name:     foldName + "/.helmignore",
				Type:     gen.File,
				Template: gw_tpls.HelmignoreTpl,
			},
			{
				Name:     foldName + "/Chart.yaml",
				Type:     gen.File,
				Template: gw_tpls.ChartTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values.yaml",
				Type:     gen.File,
				Template: gw_tpls.ValueTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-test.yaml",
				Type:     gen.File,
				Template: gw_tpls.ValueTestTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-dev.yaml",
				Type:     gen.File,
				Template: gw_tpls.ValueDevTpl,
				Data:     data,
			},
			{
				Name:     foldName + "/values-yeetown.yaml",
				Type:     gen.File,
				Template: gw_tpls.ValueYeetownTpl,
				Data:     data,
			},
			{
				Name: foldName + "/templates",
				Type: gen.Dir,
			},
			{
				Name:        foldName + "/templates/_helpers.tpl",
				Type:        gen.File,
				Template:    gw_tpls.HelpersTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/deployment.yaml",
				Type:        gen.File,
				Template:    gw_tpls.DeploymentTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/hpa.yaml",
				Type:        gen.File,
				Template:    gw_tpls.HpaYamlTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/ingress-haproxy.yaml",
				Type:        gen.File,
				Template:    gw_tpls.IngeressHaproxyTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/ingress.yaml",
				Type:        gen.File,
				Template:    gw_tpls.IngressTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/NOTES.txt",
				Type:        gen.File,
				Template:    gw_tpls.NotesTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/service.yaml",
				Type:        gen.File,
				Template:    gw_tpls.ServiceTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/serviceaccount.yaml",
				Type:        gen.File,
				Template:    gw_tpls.ServiceaccountTpl,
				NotParseTpl: true,
			},
			{
				Name:        foldName + "/templates/serviceMonitor.yaml",
				Type:        gen.File,
				Template:    gw_tpls.ServiceMonitorTpl,
				NotParseTpl: true,
			},
			{
				Name: foldName + "/templates/tests/",
				Type: gen.Dir,
			},
			{
				Name:        foldName + "/templates/tests/test-connection.yaml",
				Type:        gen.File,
				Template:    gw_tpls.TestConnectionTpl,
				NotParseTpl: true,
			},
		}); err != nil {
			return err
		}
	}
	self.Info("生成完毕")
	return nil
}
