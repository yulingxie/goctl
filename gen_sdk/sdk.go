package sdk

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"

	"gitlab.kaiqitech.com/nitro/nitro/v3/api"
	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
)

func mustParseFile(fset *token.FileSet, filename, src string) *ast.File {
	data, err := ioutil.ReadFile(filename)
	f, err := parser.ParseFile(fset, filename, string(data), parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}

func getApiEndpoints(filename string) []*api.Endpoint {
	endpoints := []*api.Endpoint{}

	fset := token.NewFileSet()
	files := []*ast.File{
		mustParseFile(fset, filename, ""),
	}
	pkg, err := doc.NewFromFiles(fset, files, "")
	if err != nil {
		panic(err)
	}
	for _, spec := range pkg.Vars[0].Decl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		for _, val := range valueSpec.Values {
			compositeLit := val.(*ast.CompositeLit)
			for _, elt := range compositeLit.Elts {
				endpoint := &api.Endpoint{
					Name:    strings.ReplaceAll(elt.(*ast.CompositeLit).Elts[0].(*ast.KeyValueExpr).Value.(*ast.BasicLit).Value, `"`, ""),
					Path:    []string{},
					Method:  []string{},
					Handler: strings.ReplaceAll(elt.(*ast.CompositeLit).Elts[3].(*ast.KeyValueExpr).Value.(*ast.BasicLit).Value, `"`, ""),
				}
				for _, pathElt := range elt.(*ast.CompositeLit).Elts[1].(*ast.KeyValueExpr).Value.(*ast.CompositeLit).Elts {
					endpoint.Path = append(endpoint.Path, strings.ReplaceAll(pathElt.(*ast.BasicLit).Value, `"`, ""))
				}
				for _, methodElt := range elt.(*ast.CompositeLit).Elts[2].(*ast.KeyValueExpr).Value.(*ast.CompositeLit).Elts {
					endpoint.Method = append(endpoint.Method, strings.ReplaceAll(methodElt.(*ast.BasicLit).Value, `"`, ""))
				}
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	return endpoints
}

func getReqOrRspContent(filenames ...string) string {
	var (
		data []byte
		err  error
	)
	for _, filename := range filenames {
		if data, err = ioutil.ReadFile(filename); err == nil {
			break
		}
	}
	if len(data) == 0 {
		logger.Infof("can not find files: %+v", filenames)
		return ""
	}
	content := string(data)

	regs := []*regexp.Regexp{
		regexp.MustCompile(`type \(([\s\S]+?)\)`),
		regexp.MustCompile(`type [^\(]([\s\S]+?)struct \{([\s\S]+?)\}`),
	}

	contents := []string{}
	for _, regexp := range regs {
		findContents := regexp.FindAllString(content, -1)
		contents = append(contents, findContents...)
	}

	content = strings.Join(contents, "\n")
	return content
}

func getImportContent(filenames ...string) string {
	var (
		data []byte
		err  error
	)
	for _, filename := range filenames {
		if data, err = ioutil.ReadFile(filename); err == nil {
			break
		}
	}
	if len(data) == 0 {
		logger.Infof("can not find files: %+v", filenames)
		return ""
	}
	content := string(data)

	regs := []*regexp.Regexp{
		regexp.MustCompile(`import \(([\s\S]+?)\)`),
		regexp.MustCompile(`import [^\(]([\s\S]+?)\n`),
	}

	contents := []string{}
	for _, regexp := range regs {
		findContents := regexp.FindAllString(content, -1)
		contents = append(contents, findContents...)
	}

	content = strings.Join(contents, "\n")
	return content
}
