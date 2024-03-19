package redis

// Text provides the default template for model to generate
var HmapModelTpl = `package redism

import (
    "context"
	"time"
	"sync"

	"github.com/go-redis/redis/v8"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/redisx"
)

type (
	{{.type}} struct{
		// 此处自定义结构体字段，记得加上redis tag
	}

	{{.type}}Model struct{
		client *redisx.Client
	}
)

var (
	{{.lowerType}}once           sync.Once
	{{.lowerType}}Model *{{.type}}Model
	{{.lowerType}}NodeName="{{.nodeName}}"
)

func Default{{.type}}Model() *{{.type}}Model {
	{{.lowerType}}once.Do(func() {
		client, _ := redisx.GetClient({{.lowerType}}NodeName)
		{{.lowerType}}Model = &{{.type}}Model{
			client: client,
		}
	})
	return {{.lowerType}}Model
}

func (self *{{.type}}Model) HGetAll(ctx context.Context, key string) (*{{.type}}, error) {
	{{.lowerType}} := &{{.type}}{}
	if err := self.client.HGetAll(ctx, key).Scan({{.lowerType}}); err != nil {
		return nil, err
	}
	return {{.lowerType}}, nil
}

// expiration: 指定过期时间，0表示没有过期时间
func (self *{{.type}}Model) HMSet(ctx context.Context, key string, {{.lowerType}} *{{.type}}, expiration time.Duration) error {
	if {{.lowerType}} == nil {
		return nil
	}
	if expiration == 0 {
		return self.client.HMSet(ctx, key, redisx.Marshal({{.lowerType}})).Err()
	}
	// 有过期时间，则将两个命令放在一个管道中处理
	_, err := self.client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.HMSet(ctx, key, redisx.Marshal({{.lowerType}}))
		pipeliner.Expire(ctx, key, expiration)
		return nil
	})
	return err
}

func (self *{{.type}}Model) Hset(ctx context.Context, key string, values ...interface{}) error {
	return self.client.HSet(ctx, key, values...).Err()
}

func (self *{{.type}}Model) Del(ctx context.Context, key string) error {
	return self.client.Del(ctx, key).Err()
}
`

var CommonModelTpl = `package redism

import (
	"sync"

	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/redisx"
)

type (
	CommonModel struct {
		client *redisx.Client
	}
)

var (
	commonOnce     sync.Once
	commonModel    *CommonModel
	commonNodeName = "common"
)

func DefaultCommonModel() *CommonModel {
	commonOnce.Do(func() {
		client, err := redisx.GetClient(commonNodeName)
		if err != nil {
			logger.Error(err.Error)
			return
		}
		commonModel = &CommonModel{
			client: client,
		}
	})
	return commonModel
}
`

// Error provides the default template for error definition in mongo code generation.
var errorTpl = `
package model

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidObjectId = errors.New("invalid objectId")
`
