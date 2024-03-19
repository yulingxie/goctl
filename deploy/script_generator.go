package deploy

import (
	"fmt"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

type ScriptGenerator struct {
	cfg *generatorConfig
}

func (s *ScriptGenerator) Generate() (string, error) {
	fmt.Println("脚本生成中 ...")
	buf, err := util.With("script").Parse(deployScript).Execute(map[string]interface{}{
		"name":       s.cfg.Name,
		"lower_name": strings.ToLower(s.cfg.Name),
		"env":        s.cfg.Env,
		"file":       s.cfg.RemoteFile,
		"url":        s.cfg.URL,
		"verbose":    s.cfg.Verbose,
	})
	if s.cfg.Verbose {
		fmt.Println(buf.String())
	}
	return buf.String(), err
}
