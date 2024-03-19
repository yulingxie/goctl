package gen

import (
	"fmt"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
)

type Generator struct {
	console.Console
	dir string
}

func NewGenerator(dir string) (*Generator, error) {
	if dir == "" {
		dir = "./"
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	if err := util.MkdirIfNotExist(dir); err != nil {
		return nil, err
	}
	return &Generator{dir: dir, Console: console.NewColorConsole()}, nil
}

func (self *Generator) Dir() string {
	return self.dir
}

func (self *Generator) GenItems(items []*Item) error {
	for _, item := range items {
		self.Info(fmt.Sprintf("生成: %s", item.Name))
		switch item.Type {
		case Dir:
			if err := util.MkdirIfNotExist(self.dir + item.Name); err != nil {
				return err
			}
		case File:
			file, err := util.CreateIfNotExist(self.dir + item.Name)
			if err != nil {
				return err
			}
			defer file.Close()

			tplStr, err := util.LoadTemplate("", "", item.Template)
			if err != nil {
				return err
			}

			// 如果不需要使用go template解析，则直接返回模版文本
			if item.NotParseTpl {
				if _, err := file.Write([]byte(tplStr)); err != nil {
					return err
				}
				continue
			}

			// 返回使用go template解析后的文本
			out, err := util.With(item.Name).Parse(tplStr).Execute(item.Data)
			if err != nil {
				return err
			}
			if _, err := file.Write(out.Bytes()); err != nil {
				return err
			}
		}
	}
	return nil
}
