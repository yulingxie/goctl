package confluence

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"html"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	goconfluence "github.com/virtomize/confluence-go-api"
	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
	"gitlab.kaiqitech.com/nitro/nitro/v3/util/templatex"
)

const (
	CONFLUENCE_ADDR            = "http://confluence.kaiqitech.com"
	CONFLUENCE_SPACE           = "K7GameServerManual"
	SNG_ERRORS_PAGE_TITLE      = "service错误码（自动更新）"
	CAC_ERRORS_PAGE_TITLE      = "cac错误码（自动更新）"
	SNG_ERRORS_PAGE_TYPE       = "page"
	ERRORS_CLOUMN_NUMS         = 3
	SNG_TEST_PAGE_TITLE        = "service测试数据汇总（自动更新）"
	SNG_TEST_PAGE_TYPE         = "page"
	TEST_CLOUMN_NUMS           = 7
	SNG_CHANGELOG_TITLE        = "%s(CHANGELOG)"
	SNG_CHANGELOG_PARENT_TITLE = "版本更新日志(CHANGELOG)"
	SNG_CHANGELOG_PAGE_TYPE    = "page"
)

type Error struct {
	Code int
	Name string
	Msgs []string
}

type ErrorArray []*Error

func (self ErrorArray) Len() int {
	return len(self)
}

func (self ErrorArray) Less(i, j int) bool {
	return self[i].Code < self[j].Code
}

func (self ErrorArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type Test struct {
	Name               string
	Version            string
	UnitTestReportUrl  string
	Coverage           string
	HasApiTest         string
	HasConnserverTest  string
	HasGameserverTest  string
	HasBenchmark       string
	BenchmarkReportUrl string
}

type TestArray []*Test

func (self TestArray) Len() int {
	return len(self)
}

func (self TestArray) Less(i, j int) bool {
	return self[i].Name < self[j].Name
}

func (self TestArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type ChangeLog struct {
	Version     string
	ReleaseDate string
	Title       string
	Msgs        []string
}
type LogArray []*ChangeLog

func (self LogArray) Len() int {
	return len(self)
}

func (self LogArray) Less(i, j int) bool {
	return self[i].Version < self[j].Version
}

func (self LogArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type ChangeLogV2 struct {
	LatestVersion      string
	LatestReleaseDate  string
	LatestDownloadLink string
	LatestMsgs         []string
	HistoricalVersion  []*ChangeLog
}

type Confluence struct {
	goconfluenceapi *goconfluence.API
	api             *API
}

func NewConfluence(username, pass string) *Confluence {
	confluenceIns, _ := goconfluence.NewAPI(CONFLUENCE_ADDR+"/rest/api", username, pass)
	api := NewAPI(CONFLUENCE_ADDR, username, pass)
	confluece := &Confluence{
		goconfluenceapi: confluenceIns,
		api:             api,
	}
	return confluece
}

func (self *Confluence) UpdateSngErrors(filename string) error {
	page, _ := self.api.FindPage(CONFLUENCE_SPACE, SNG_ERRORS_PAGE_TITLE, SNG_ERRORS_PAGE_TYPE)
	content, _ := templatex.With("").Parse(sngErrorsTpl).Execute(map[string]interface{}{
		"errors": mergeErrors(self.getOldErrors(), self.getNewErrors(filename)),
	})
	str := content.String()
	return self.api.UpdatePage(page, str, false, []string{})
}

func (self *Confluence) UpdateCacErrors(filename string) error {
	page, _ := self.api.FindPage(CONFLUENCE_SPACE, CAC_ERRORS_PAGE_TITLE, SNG_ERRORS_PAGE_TYPE)
	errors := self.getNewErrors(filename)
	content, _ := templatex.With("").Parse(sngErrorsTpl).Execute(map[string]interface{}{
		"errors": errors,
	})
	str := content.String()
	return self.api.UpdatePage(page, str, false, []string{})
}

func (self *Confluence) getOldErrors() ErrorArray {
	oldPage, _ := self.goconfluenceapi.GetContent(goconfluence.ContentQuery{
		Expand:   []string{"body.view"},
		SpaceKey: CONFLUENCE_SPACE,
		Title:    SNG_ERRORS_PAGE_TITLE,
		Type:     SNG_ERRORS_PAGE_TYPE,
	})
	// logger.Info(oldPage.Results[0].Body.View.Value)
	reg1 := regexp.MustCompile(`<td [a-zA-Z]+.*?>([\s\S]*?)</td>`)
	reg2 := regexp.MustCompile(`<td [a-zA-Z]+.*?>`)
	reg3 := regexp.MustCompile(`</td>`)
	results := reg1.FindAllString(oldPage.Results[0].Body.View.Value, -1)
	for index, result := range results {
		bytes := reg2.ReplaceAll([]byte(result), []byte(""))
		bytes = reg3.ReplaceAll(bytes, []byte(""))
		results[index] = string(bytes)
	}
	errors := ErrorArray{}
	reg4 := regexp.MustCompile(`<p>([\s\S]*?)</p>`)
	for i := 0; i < len(results)/ERRORS_CLOUMN_NUMS; i++ {
		code, _ := strconv.Atoi(results[ERRORS_CLOUMN_NUMS*i])
		msgs := reg4.FindAllString(results[ERRORS_CLOUMN_NUMS*i+2], 100)
		for index, msg := range msgs {
			msg = strings.ReplaceAll(msg, "<p>", "")
			msg = strings.ReplaceAll(msg, "</p>", "")
			msgs[index] = msg
		}
		errors = append(errors, &Error{
			Code: code,
			Name: results[ERRORS_CLOUMN_NUMS*i+1],
			Msgs: msgs,
		})
	}
	return errors
}

func (self *Confluence) createCodeFile() error {
	oldErrors := self.getOldErrors()
	for _, err := range oldErrors {
		if len(err.Msgs) == 0 {
			err.Msgs = []string{"未知"}
		}
	}
	data, err := templatex.With("").Parse(codeJsonTpl).Execute(map[string]interface{}{
		"codes": oldErrors,
	})
	if err != nil {
		logger.Error(err)
		return err
	}
	result := strings.ReplaceAll(data.String(), "zh-CN: ", "")
	result = strings.ReplaceAll(result, "zh_CN: ", "")
	logger.Info(result)
	return nil
}

func mustParseFile(fset *token.FileSet, filename, src string) *ast.File {
	data, err := ioutil.ReadFile(filename)
	var f *ast.File
	if strings.HasSuffix(filename, ".proto") {
		newfilename := strings.Replace(filename, ".proto", ".go", 1)
		os.Rename(filename, newfilename)
		dataString := strings.ReplaceAll(strings.ReplaceAll(string(data), "syntax = \"proto3\";", ""), "option go_package = \"./commonpb\";", "")
		dataString = strings.ReplaceAll(strings.ReplaceAll(dataString, "enum ErrorCode {", "const ("), "}", ")")
		f, err = parser.ParseFile(fset, newfilename, dataString, parser.ParseComments)
	} else {
		f, err = parser.ParseFile(fset, filename, string(data), parser.ParseComments)
	}

	if err != nil {
		panic(any(err))
	}
	return f
}

func (self *Confluence) getNewErrors(filename string) ErrorArray {
	fset := token.NewFileSet()
	files := []*ast.File{
		mustParseFile(fset, filename, ""),
	}
	pkg, err := doc.NewFromFiles(fset, files, "")
	if err != nil {
		panic(any(err))
	}
	errors := ErrorArray{}
	for _, spec := range pkg.Consts[0].Decl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		value := valueSpec.Values[0].(*ast.BasicLit)
		code, _ := strconv.Atoi(value.Value)
		msgs := []string{}
		if valueSpec.Doc != nil {
			for _, comment := range valueSpec.Doc.List {
				msgs = append(msgs, strings.ReplaceAll(html.EscapeString(comment.Text), "//", ""))
			}
		}
		errors = append(errors, &Error{
			Code: code,
			Name: valueSpec.Names[0].Name,
			Msgs: msgs,
		})
	}
	return errors
}

// 若有重复的code，后传的会覆盖之前的
func mergeErrors(errArrays ...ErrorArray) ErrorArray {
	errors := map[int]*Error{}
	for _, errArray := range errArrays {
		for _, err := range errArray {
			errors[err.Code] = err
		}
	}
	result := ErrorArray{}
	for _, err := range errors {
		result = append(result, err)
	}
	sort.Sort(result)
	return result
}

func getVersionAndReleaseDate(rawStr string) *ChangeLog {
	out := &ChangeLog{}
	data := strings.ReplaceAll(rawStr, "# ", "")
	data = strings.ReplaceAll(data, "[", "")
	data = strings.ReplaceAll(data, "]", "")
	vr := strings.Split(data, " ")
	if len(vr) == 2 {
		out.Version = strings.ReplaceAll(vr[0], "#", "")
		out.ReleaseDate = vr[1]
	} else if len(vr) == 1 {
		out.Version = strings.ReplaceAll(vr[0], "#", "")
	} else {
		out.Version = strings.ReplaceAll(vr[0], "#", "")
		out.ReleaseDate = vr[len(vr)-1]
	}
	return out
}

func GetConfigByLog(logStr string) (log, cfg string) {
	//changelog中包含配置, 将changelog和yaml配置分离开并返回
	var newlogs []string
	logs := strings.Split(logStr, "```")
	for _, v := range logs {
		if len(v) > 0 {
			newlogs = append(newlogs, v)
		}
	}
	if len(newlogs) == 2 {
		conf := strings.ReplaceAll(newlogs[len(newlogs)-1], "\n", "<br/>")
		return newlogs[0], "```" + conf + "```"
	}
	return "", ""
}

func GetChangeLogV2(logStr string) *ChangeLog {
	var newlogs, cfg string
	var logs []string
	if strings.Contains(logStr, "```") {
		newlogs, cfg = GetConfigByLog(logStr)
	}
	if newlogs != "" {
		logs = strings.Split(newlogs, "\n")
	} else {
		logs = strings.Split(logStr, "\n")
	}
	var logArr []string
	for _, v := range logs {
		if len(v) > 0 {
			logArr = append(logArr, v)
		}
	}
	if len(logArr) == 0 {
		return nil
	}
	log := getVersionAndReleaseDate(logArr[0])
	if len(log.ReleaseDate) > 0 {
		log.Title = fmt.Sprintf("%s[%s]", log.Version, log.ReleaseDate)
	} else {
		log.Title = log.Version
	}
	for i, v := range logArr {
		if i > 0 {
			log.Msgs = append(log.Msgs, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v, "*", ""), "#", ""), "-", ""))
		}
	}
	if len(cfg) > 0 {
		log.Msgs = append(log.Msgs, cfg)
	}
	return log
}

func GetChangeLog(logStr string) *ChangeLog {
	logArr := strings.SplitN(logStr, "\n", 2)
	if len(logArr) == 2 {
		//已成功获取第一行版本信息和changelog消息
		log := getVersionAndReleaseDate(logArr[0])
		if len(log.ReleaseDate) > 0 {
			log.Title = fmt.Sprintf("%s[%s]", log.Version, log.ReleaseDate)
		} else {
			log.Title = log.Version
		}
		for _, v := range strings.Split(logArr[1], "\n") {
			log.Msgs = append(log.Msgs, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v, "*", ""), "#", ""), "-", ""))
		}
		return log
	}
	return nil
}

func (self *Confluence) getNewChangeLog(filename string) LogArray {
	var out LogArray
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	winStr := strings.ReplaceAll(string(data), "\r\n", "\n")
	logArr := strings.Split(winStr, "\n\n")
	for _, v := range logArr {
		tmp := GetChangeLog(v)
		if tmp != nil {
			out = append(out, tmp)
		}
	}
	return out
}

func (self *Confluence) UpdateChangeLog(filename, serviceName string) error {
	page, _ := self.api.FindPage(CONFLUENCE_SPACE, fmt.Sprintf(SNG_CHANGELOG_TITLE, strings.ToLower(serviceName)), SNG_CHANGELOG_PAGE_TYPE)
	content, _ := templatex.With("").Parse(sngChangeLogTpl).Execute(map[string]interface{}{
		"changelogs": self.getNewChangeLog(filename),
	})
	str := content.String()
	if page == nil {
		parentPage, _ := self.api.FindPage(CONFLUENCE_SPACE, SNG_CHANGELOG_PARENT_TITLE, SNG_CHANGELOG_PAGE_TYPE)
		page, _ = self.api.CreatePage(CONFLUENCE_SPACE, SNG_CHANGELOG_PAGE_TYPE, parentPage, fmt.Sprintf(SNG_CHANGELOG_TITLE, strings.ToLower(serviceName)), str)
		return nil
	}
	return self.api.UpdatePage(page, str, false, []string{})
}

func (self *Confluence) getNewMarkdownChangeLog(filename, link string) *ChangeLogV2 {
	var out = &ChangeLogV2{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	winStr := strings.ReplaceAll(string(data), "\r\n", "\n")
	logArr := strings.Split(winStr, "\n\n")
	for i, v := range logArr {
		tmp := GetChangeLogV2(v)
		if tmp != nil {
			//logger.Debugf("v: %+v\ntmp: %+v\n", v, tmp)
			if i == 0 {
				out.LatestVersion = tmp.Version
				out.LatestReleaseDate = tmp.ReleaseDate
				out.LatestMsgs = tmp.Msgs
			} else {
				out.HistoricalVersion = append(out.HistoricalVersion, tmp)
			}
		}
	}
	if len(link) > 0 {
		out.LatestDownloadLink = link
	}
	return out
}

func (self *Confluence) UpdateChangeLogMarkdown(filename, serviceName, ParentPage, link, ConfluenceSpace string) error {
	if len(ConfluenceSpace) == 0 {
		ConfluenceSpace = CONFLUENCE_SPACE
	}
	page, _ := self.api.FindPage(ConfluenceSpace, fmt.Sprintf(SNG_CHANGELOG_TITLE, strings.ToUpper(serviceName)), SNG_CHANGELOG_PAGE_TYPE)
	content, _ := templatex.With("").Parse(markdownChangeLogTpl).Execute(self.getNewMarkdownChangeLog(filename, link))
	str := content.String()
	if page == nil {
		var PARENT_TITLE string
		serviceName = strings.ToUpper(serviceName)
		if len(ParentPage) > 0 {
			PARENT_TITLE = ParentPage
		} else {
			if strings.HasSuffix(serviceName, "-SERVICE") {
				PARENT_TITLE = "SNG"
			} else if strings.HasPrefix(serviceName, "CAC-") {
				PARENT_TITLE = "CAC"
			} else if strings.Contains(serviceName, "SERVER") {
				PARENT_TITLE = "SERVER"
			} else {
				PARENT_TITLE = SNG_CHANGELOG_PARENT_TITLE
			}
		}
		parentPage, _ := self.api.FindPage(ConfluenceSpace, PARENT_TITLE, SNG_CHANGELOG_PAGE_TYPE)
		_, err := self.api.CreatePage(ConfluenceSpace, SNG_CHANGELOG_PAGE_TYPE, parentPage, fmt.Sprintf(SNG_CHANGELOG_TITLE, serviceName), str)
		if err != nil {
			logger.Errorf("创建页面(父页面: %v)失败: %+v", PARENT_TITLE, err)
		}
		return nil
	}
	return self.api.UpdatePage(page, str, false, []string{})
}

func (self *Confluence) UpdateSngDevBuild(name, version, coverage, reportUrl string) error {
	page, _ := self.api.FindPage(CONFLUENCE_SPACE, SNG_TEST_PAGE_TITLE, SNG_TEST_PAGE_TYPE)
	oldTests := self.getOldTests()
	if oldTests[name] == nil {
		oldTests[name] = &Test{
			Name:               name,
			Version:            version,
			Coverage:           coverage,
			UnitTestReportUrl:  reportUrl,
			HasApiTest:         "no",
			HasConnserverTest:  "no",
			HasGameserverTest:  "no",
			HasBenchmark:       "no",
			BenchmarkReportUrl: "-",
		}
	} else {
		oldTests[name].Name = name
		oldTests[name].Version = version
		oldTests[name].Coverage = coverage
		oldTests[name].UnitTestReportUrl = reportUrl
	}
	content, err := templatex.With("").Parse(sngTestTpl).Execute(map[string]interface{}{
		"tests": oldTests,
	})
	if err != nil {
		logger.Error(err.Error())
	}
	return self.api.UpdatePage(page, content.String(), false, []string{})
}

func (self *Confluence) UpdateSngTestBuild(name, hasApiTest, hasConnserverTest, hasGameserverTest, hasBenchmark, benchmarkReportUrl string) error {
	page, _ := self.api.FindPage(CONFLUENCE_SPACE, SNG_TEST_PAGE_TITLE, SNG_TEST_PAGE_TYPE)
	oldTests := self.getOldTests()
	if oldTests[name] == nil {
		oldTests[name] = &Test{
			Name:               name,
			Version:            "-",
			Coverage:           "-",
			UnitTestReportUrl:  "-",
			HasApiTest:         hasApiTest,
			HasConnserverTest:  hasConnserverTest,
			HasGameserverTest:  hasGameserverTest,
			HasBenchmark:       hasBenchmark,
			BenchmarkReportUrl: benchmarkReportUrl,
		}
	} else {
		oldTests[name].HasApiTest = hasApiTest
		oldTests[name].HasConnserverTest = hasConnserverTest
		oldTests[name].HasGameserverTest = hasGameserverTest
		oldTests[name].HasBenchmark = hasBenchmark
		oldTests[name].BenchmarkReportUrl = benchmarkReportUrl
	}
	content, err := templatex.With("").Parse(sngTestTpl).Execute(map[string]interface{}{
		"tests": oldTests,
	})
	if err != nil {
		logger.Error(err.Error())
	}
	return self.api.UpdatePage(page, content.String(), false, []string{})
}

func (self *Confluence) getOldTests() map[string]*Test {
	oldPage, _ := self.goconfluenceapi.GetContent(goconfluence.ContentQuery{
		Expand:   []string{"body.view"},
		SpaceKey: CONFLUENCE_SPACE,
		Title:    SNG_TEST_PAGE_TITLE,
		Type:     SNG_TEST_PAGE_TYPE,
	})
	reg1 := regexp.MustCompile(`<td [a-zA-Z]+.*?>([\s\S]*?)</td>`)
	reg2 := regexp.MustCompile(`<td [a-zA-Z]+.*?>`)
	reg3 := regexp.MustCompile(`</td>`)
	results := reg1.FindAllString(oldPage.Results[0].Body.View.Value, 200)
	for index, result := range results {
		bytes := reg2.ReplaceAll([]byte(result), []byte(""))
		bytes = reg3.ReplaceAll(bytes, []byte(""))
		results[index] = string(bytes)
	}

	reg4 := regexp.MustCompile(`href="([\s\S]*?)"`)
	reg5 := regexp.MustCompile(`<a [a-zA-Z]+.*?>`)
	reg6 := regexp.MustCompile(`</a>`)
	tests := map[string]*Test{}
	for i := 0; i < len(results)/TEST_CLOUMN_NUMS; i++ {
		test := &Test{
			Name:              results[TEST_CLOUMN_NUMS*i],
			Version:           results[TEST_CLOUMN_NUMS*i+1],
			HasApiTest:        results[TEST_CLOUMN_NUMS*i+3],
			HasConnserverTest: results[TEST_CLOUMN_NUMS*i+4],
			HasGameserverTest: results[TEST_CLOUMN_NUMS*i+5],
		}
		bytes := reg5.ReplaceAll([]byte(results[TEST_CLOUMN_NUMS*i+2]), []byte(""))
		bytes = reg6.ReplaceAll(bytes, []byte(""))
		test.Coverage = string(bytes)
		test.UnitTestReportUrl = reg4.FindString(results[TEST_CLOUMN_NUMS*i+2])
		test.UnitTestReportUrl = strings.Replace(test.UnitTestReportUrl, "href=\"", "", 100)
		test.UnitTestReportUrl = strings.Replace(test.UnitTestReportUrl, "\"", "", 100)

		bytes = reg5.ReplaceAll([]byte(results[TEST_CLOUMN_NUMS*i+6]), []byte(""))
		bytes = reg6.ReplaceAll(bytes, []byte(""))
		test.HasBenchmark = string(bytes)
		test.BenchmarkReportUrl = reg4.FindString(results[TEST_CLOUMN_NUMS*i+6])
		test.BenchmarkReportUrl = strings.Replace(test.BenchmarkReportUrl, "href=\"", "", 100)
		test.BenchmarkReportUrl = strings.Replace(test.BenchmarkReportUrl, "\"", "", 100)
		tests[test.Name] = test
	}
	return tests
}

func mergeTests(testArrays ...TestArray) TestArray {
	tests := map[string]*Test{}
	for _, testArray := range testArrays {
		for _, test := range testArray {
			tests[test.Name] = test
		}
	}
	result := TestArray{}
	for _, test := range tests {
		result = append(result, test)
	}
	sort.Sort(result)
	return result
}
