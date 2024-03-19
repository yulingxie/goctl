package template

var (
	// Imports defines a import template for model in cache case
	ImportsWithCache = `import (
	"context"
	"fmt"
	"sync"
	{{if .time}}"time"{{end}}

	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/cache"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/redisx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/sqlx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/util/syncx"
	{{if .withGameId}}seraphSqlm "gitlab.kaiqitech.com/k7game/server/components/seraph.git/model/sqlm"{{end}}
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"context"
	"sync"
	{{if .time}}"time"{{end}}

	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/sqlx"
	{{if .withGameId}}seraphSqlm "gitlab.kaiqitech.com/k7game/server/components/seraph.git/model/sqlm/gamemodel"{{end}}
)
`
)
