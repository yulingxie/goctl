package jenkins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
	"gitlab.kaiqitech.com/nitro/nitro/v3/util/templatex"
)

type Crumb struct {
	Crumb string `json:"crumb,omitempty"`
}

type JenKins struct {
	console.Console
	httpClient *http.Client
	username   string
	pass       string
	crumb      string
	cookie     string
}

func NewJenkins(username, pass string) *JenKins {
	jenkins := &JenKins{
		Console: console.NewColorConsole(),
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost:     1,
				MaxIdleConns:        1,
				MaxIdleConnsPerHost: 1,
				IdleConnTimeout:     time.Minute,
			},
		},
		username: username,
		pass:     pass,
	}
	if err := jenkins.initCrumb(); err != nil {
		panic(err)
	}
	return jenkins
}

func (self *JenKins) newRequest(method, url, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(self.username, self.pass)
	req.Header.Set("Jenkins-Crumb", self.crumb)
	req.Header.Set("Cookie", self.cookie)
	req.Header.Set("Content-Type", "application/xml")
	return req, nil
}

func (self *JenKins) initCrumb() error {
	req, err := self.newRequest("GET", JENKINS_CRUMB_URL, "")
	if err != nil {
		return err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if rsp.StatusCode != 200 {
		return errors.New(string(body))
	}

	crumb := &Crumb{}
	if err := json.Unmarshal(body, crumb); err != nil {
		return err
	}

	self.crumb = crumb.Crumb
	self.cookie = strings.Split(rsp.Header.Get("Set-Cookie"), ";")[0]
	return nil
}

func (self *JenKins) getJob(url string) (string, error) {
	url = fmt.Sprintf("%s/config.xml", url)
	req, err := self.newRequest("GET", url, "")
	if err != nil {
		return "", err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}

func (self *JenKins) hasJob(url string) bool {
	url = fmt.Sprintf("%s/config.xml", url)
	req, err := self.newRequest("GET", url, "")
	if err != nil {
		self.Error(err.Error())
		return false
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		self.Error(err.Error())
		return false
	}
	defer rsp.Body.Close()
	return rsp.StatusCode == 200
}

func (self *JenKins) createJob(url, name, configXml string) error {
	url = fmt.Sprintf("%s/createItem?name=%s", url, name)
	req, err := self.newRequest("POST", url, configXml)
	if err != nil {
		return err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	self.Info("%s 创建完毕", url)
	return nil
}

func (self *JenKins) updateJob(url, configXml string) error {
	url = fmt.Sprintf("%s/config.xml", url)
	req, err := self.newRequest("POST", url, configXml)
	if err != nil {
		return err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	self.Info("%s 更新完毕", url)
	return nil
}

func (self *JenKins) deleteJob(url string) error {
	url = fmt.Sprintf("%s/doDelete", url)
	req, err := self.newRequest("POST", url, "")
	if err != nil {
		return err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	self.Info("%s 删除完毕", url)
	return nil
}

func (self *JenKins) buildJob(url string) error {
	url = fmt.Sprintf("%s/build", url)
	req, err := self.newRequest("POST", url, "")
	if err != nil {
		return err
	}

	rsp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	self.Info("%s 开始构建", url)
	return nil
}

// 取account-service的jenkins配置作为模版
func (self *JenKins) getSngServiceTemplate(branch string) string {
	tpl, err := self.getJob(fmt.Sprintf("http://pink-jenkins.kaiqitech.com/job/service/job/account-service/job/%s/", branch))
	if err != nil {
		logger.Error(err.Error())
	}

	tpl = strings.Replace(tpl, "<?xml version='1.1' encoding='UTF-8'?>", "", -1)
	tpl = strings.Replace(tpl, "account", "{{.name}}", -1)
	return tpl
}

// 取account-service的jenkins配置作为模版
func (self *JenKins) getSngGateWayTemplate(branch string) string {
	tpl, err := self.getJob(fmt.Sprintf("http://pink-jenkins.kaiqitech.com/job/gateway/job/sapi-gw/job/%s/", branch))
	if err != nil {
		logger.Error(err.Error())
	}

	tpl = strings.Replace(tpl, "<?xml version='1.1' encoding='UTF-8'?>", "", -1)
	tpl = strings.Replace(tpl, "sapi-gw", "{{.name}}", -1)
	return tpl
}

// 取sapi-gw的jenkins配置作为模版
func (self *JenKins) CreateSngService(names ...string) error {
	devTpl := self.getSngServiceTemplate("dev")
	testTpl := self.getSngServiceTemplate("test")
	masterTpl := self.getSngServiceTemplate("master")
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		name = strings.Split(name, "-")[0]
		serviceName := name + "-service"
		data := map[string]interface{}{"name": name}
		serviceJobUrl := SNG_SERVICE_URL + "/job/" + serviceName

		if self.hasJob(serviceJobUrl) {
			self.Info(fmt.Sprintf("服务%s已存在", serviceName))
			continue
		}

		// 创建文件夹
		devBuf, err := templatex.With("").Parse(ServiceFoldTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(SNG_SERVICE_URL, serviceName, devBuf.String()); err != nil {
			return err
		}

		// 创建dev
		devBuf, err = templatex.With("").Parse(devTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(serviceJobUrl, "dev", devBuf.String()); err != nil {
			return err
		}
		// 创建test
		qaBuf, err := templatex.With("").Parse(testTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(serviceJobUrl, "test", qaBuf.String()); err != nil {
			return err
		}
		// 创建master
		trunkBuf, err := templatex.With("").Parse(masterTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(serviceJobUrl, "master", trunkBuf.String()); err != nil {
			return err
		}
	}
	return nil
}

func (self *JenKins) UpdateSngService(names ...string) error {
	devTpl := self.getSngServiceTemplate("dev")
	testTpl := self.getSngServiceTemplate("test")
	masterTpl := self.getSngServiceTemplate("master")
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		name = strings.Split(name, "-")[0]
		data := map[string]interface{}{"name": name}
		serviceName := name + "-service"
		serviceJobUrl := SNG_SERVICE_URL + "/job/" + serviceName

		// 更新dev
		devBuf, err := templatex.With("").Parse(devTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(serviceJobUrl+"/job/dev", devBuf.String()); err != nil {
			return err
		}
		// 更新test
		qaBuf, err := templatex.With("").Parse(testTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(serviceJobUrl+"/job/test", qaBuf.String()); err != nil {
			return err
		}
		// 更新master
		trunkBuf, err := templatex.With("").Parse(masterTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(serviceJobUrl+"/job/master", trunkBuf.String()); err != nil {
			return err
		}
	}
	return nil
}

func (self *JenKins) DeleteSngService(names ...string) error {
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		// 删除dev
		if err := self.deleteJob(fmt.Sprintf("%s/job/%s", SNG_SERVICE_URL, "dev")); err != nil {
			return err
		}
		// 删除test
		if err := self.deleteJob(fmt.Sprintf("%s/job/%s", SNG_SERVICE_URL, "test")); err != nil {
			return err
		}
		// 删除trunk
		if err := self.deleteJob(fmt.Sprintf("%s/job/%s", SNG_SERVICE_URL, "master")); err != nil {
			return err
		}
	}
	return nil
}

func (self *JenKins) BuildSngServiceDev(names ...string) error {
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		self.buildJob(fmt.Sprintf("%s/job/%s/job/%s", SNG_SERVICE_URL, name, "dev"))
	}
	return nil
}

func (self *JenKins) BuildSngServiceTest(names ...string) error {
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		self.buildJob(fmt.Sprintf("%s/job/%s/job/%s", SNG_SERVICE_URL, name, "test"))
	}
	return nil
}

func (self *JenKins) CreateSngGateway(names ...string) error {
	devTpl := self.getSngGateWayTemplate("dev")
	testTpl := self.getSngGateWayTemplate("test")
	masterTpl := self.getSngGateWayTemplate("master")
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		data := map[string]interface{}{"name": name}
		gatewayJobUrl := SNG_GATEWAY_URL + "/job/" + name

		if self.hasJob(gatewayJobUrl) {
			self.Info(fmt.Sprintf("服务%s已存在", name))
			continue
		}

		// 创建文件夹
		devBuf, err := templatex.With("").Parse(GatewayFoldTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(SNG_GATEWAY_URL, name, devBuf.String()); err != nil {
			return err
		}

		// 创建dev
		devBuf, err = templatex.With("").Parse(devTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(gatewayJobUrl, "dev", devBuf.String()); err != nil {
			return err
		}
		// 创建test
		qaBuf, err := templatex.With("").Parse(testTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(gatewayJobUrl, "test", qaBuf.String()); err != nil {
			return err
		}
		// 创建master
		trunkBuf, err := templatex.With("").Parse(masterTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.createJob(gatewayJobUrl, "master", trunkBuf.String()); err != nil {
			return err
		}
	}
	return nil
}

func (self *JenKins) UpdateSngGateway(names ...string) error {
	devTpl := self.getSngGateWayTemplate("dev")
	testTpl := self.getSngGateWayTemplate("test")
	masterTpl := self.getSngGateWayTemplate("master")
	for _, name := range names {
		if len(name) == 0 {
			return errors.New("未指定服务名")
		}
		data := map[string]interface{}{"name": name}
		gatewayJobUrl := SNG_GATEWAY_URL + "/job/" + name

		// 更新dev
		devBuf, err := templatex.With("").Parse(devTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(gatewayJobUrl+"/job/dev", devBuf.String()); err != nil {
			return err
		}
		// 更新test
		qaBuf, err := templatex.With("").Parse(testTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(gatewayJobUrl+"/job/test", qaBuf.String()); err != nil {
			return err
		}
		// 更新master
		trunkBuf, err := templatex.With("").Parse(masterTpl).Execute(data)
		if err != nil {
			return err
		}
		if err := self.updateJob(gatewayJobUrl+"/job/master", trunkBuf.String()); err != nil {
			return err
		}
	}
	return nil
}
