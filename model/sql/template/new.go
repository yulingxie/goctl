package template

var New = `
var (
	{{.lowerCamelTableName}}ModelOnce    sync.Once
	default{{.upperStartCamelObject}}Model *{{.upperStartCamelObject}}Model
	{{.lowerCamelTableName}}NodeName       = "{{.nodeName}}"
	{{.lowerCamelTableName}}DbName       = "{{.dbName}}"
	{{.lowerCamelTableName}}TableName = "{{.table}}"
)

func Default{{.upperStartCamelObject}}Model() *{{.upperStartCamelObject}}Model {
	{{.lowerCamelTableName}}ModelOnce.Do(func() {
		client, err := sqlx.GetClient({{.lowerCamelTableName}}NodeName, {{.lowerCamelTableName}}DbName)
		if err != nil {
			logger.Error(err)
			return
		}
		default{{.upperStartCamelObject}}Model = &{{.upperStartCamelObject}}Model{
			client: client,
		}
	})
	return default{{.upperStartCamelObject}}Model
}
`

var NewWithCache = `
var (
	{{.lowerCamelTableName}}ModelOnce    sync.Once
	default{{.upperStartCamelObject}}Model *{{.upperStartCamelObject}}Model
	{{.lowerCamelTableName}}NodeName       = "{{.nodeName}}"
	{{.lowerCamelTableName}}DbName       = "{{.dbName}}"
	{{.lowerCamelTableName}}TableName = "{{.table}}"
)

func Default{{.upperStartCamelObject}}Model() *{{.upperStartCamelObject}}Model {
	{{.lowerCamelTableName}}ModelOnce.Do(func() {
		client, err := sqlx.GetClient({{.lowerCamelTableName}}NodeName, {{.lowerCamelTableName}}DbName)
		if err != nil {
			logger.Error(err)
			return
		}
		redisClient, err := redisx.GetClient("cache")
		if err != nil {
			logger.Error(err)
			return
		}
		default{{.upperStartCamelObject}}Model = &{{.upperStartCamelObject}}Model{
			client: client,
			cache:  cache.NewNode(redisClient, syncx.NewSharedCalls()),
			cacheKey: "{{.dbName}}:{{.table}}:%v",
		}
	})
	return default{{.upperStartCamelObject}}Model
}
`

var NewWithGameId = `
var (
	{{.lowerCamelTableName}}Lock         sync.RWMutex
	default{{.upperStartCamelObject}}Models map[uint32]*{{.upperStartCamelObject}}Model = map[uint32]*{{.upperStartCamelObject}}Model{}
	{{.lowerCamelTableName}}NodeName       = "{{.nodeName}}"
	{{.lowerCamelTableName}}DbName       = "{{.dbName}}"
	{{.lowerCamelTableName}}TableName = "{{.table}}"
)

func Default{{.upperStartCamelObject}}Model(gameId uint32) *{{.upperStartCamelObject}}Model {
	{{.lowerCamelTableName}}Lock.RLock()
	model := default{{.upperStartCamelObject}}Models[gameId]
	{{.lowerCamelTableName}}Lock.RUnlock()
	if model == nil {
		game, err := seraphSqlm.DefaulGameModel().FindValidOne(gameId)
		if err != nil {
			logger.Error(err)
			return nil
		}
		client, err := sqlx.GetClient({{.lowerCamelTableName}}NodeName, game.DbName)
		if err != nil {
			logger.Error(err)
			return nil
		}
		model = &{{.upperStartCamelObject}}Model{
			client: client,
		}
		{{.lowerCamelTableName}}Lock.Lock()
		default{{.upperStartCamelObject}}Models[gameId] = model
		{{.lowerCamelTableName}}Lock.Unlock()
	}
	return model
}
`

var NewWitchCacheAndGameId = `
var (
	{{.lowerCamelTableName}}Lock         sync.RWMutex
	default{{.upperStartCamelObject}}Models map[uint32]*{{.upperStartCamelObject}}Model = map[uint32]*{{.upperStartCamelObject}}Model{}
	{{.lowerCamelTableName}}NodeName       = "{{.nodeName}}"
	{{.lowerCamelTableName}}DbName       = "{{.dbName}}"
	{{.lowerCamelTableName}}TableName = "{{.table}}"
)

func Default{{.upperStartCamelObject}}Model(gameId uint32) *{{.upperStartCamelObject}}Model {
	{{.lowerCamelTableName}}Lock.RLock()
	model := default{{.upperStartCamelObject}}Models[gameId]
	{{.lowerCamelTableName}}Lock.RUnlock()
	if model == nil {
		game, err := seraphSqlm.DefaulGameModel().FindValidOne(gameId)
		if err != nil {
			logger.Error(err)
			return nil
		}
		client, err := sqlx.GetClient({{.lowerCamelTableName}}NodeName, game.DbName)
		if err != nil {
			logger.Error(err)
			return nil
		}
		redisClient, err := redisx.GetClient("cache")
		if err != nil {
			logger.Error(err)
			return nil
		}
		model = &{{.upperStartCamelObject}}Model{
			client: client,
			cache:  cache.NewNode(redisClient, syncx.NewSharedCalls()),
			cacheKey: game.DbName + ":{{.table}}:%v",
		}
		{{.lowerCamelTableName}}Lock.Lock()
		default{{.upperStartCamelObject}}Models[gameId] = model
		{{.lowerCamelTableName}}Lock.Unlock()
	}
	return model
}
`
