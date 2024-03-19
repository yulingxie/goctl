package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

type ItemType int

const (
	Dir ItemType = iota
	File
)

type Item struct {
	Name     string
	Type     ItemType
	Template string
	Data     interface{}
	// 不解析模板文本，而是直接使用原文
	NotParseTpl bool
}

func GenItems(dir string, data map[string]interface{}, items []*Item) error {
	for _, item := range items {
		switch item.Type {
		case Dir:
			if err := util.MkdirIfNotExist(dir + item.Name); err != nil {
				return err
			}
		case File:
			file, err := util.CreateIfNotExist(dir + item.Name)
			if err != nil {
				return err
			}
			defer file.Close()

			tplStr, err := util.LoadTemplate("", "", item.Template)
			if err != nil {
				return err
			}

			// 如果这个模板文本不需要使用go语法解析
			out, err := util.With(item.Name).Parse(tplStr).Execute(data)
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
