package mongo

// Text provides the default template for model to generate
var modelTpl = `package mongom

import (
	"context"
	"sync"

	"gitlab.kaiqitech.com/nitro/nitro/v3/stores/mongox"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	{{.typeName}} struct {
		// 此处自定义结构体字段,记得加上bson tag
	}

	{{.typeName}}Model struct {
		client *mongox.Client
	}
)

var (
	{{.lowerTypeName}}Once     sync.Once
	{{.lowerTypeName}}Model    *{{.typeName}}Model
	{{.lowerTypeName}}NodeName = "{{.nodeName}}"
	{{.lowerTypeName}}DBName   = "{{.dbName}}"
	{{.lowerTypeName}}CollName = "{{.collName}}"
)

func Default{{.typeName}}Model() *{{.typeName}}Model {
	{{.lowerTypeName}}Once.Do(func() {
		client, _ := mongox.GetClient({{.lowerTypeName}}NodeName, {{.lowerTypeName}}DBName, {{.lowerTypeName}}CollName)
		{{.lowerTypeName}}Model = &{{.typeName}}Model{
			client: client,
		}
	})
	return {{.lowerTypeName}}Model
}

func (self *{{.typeName}}Model) FindOne(ctx context.Context, filter interface{}) (*{{.typeName}}, error) {
	{{.lowerTypeName}} := &{{.typeName}}{}
	err := self.client.Find(ctx, filter).One({{.lowerTypeName}})
	if err != nil {
		return nil, err
	}
	return {{.lowerTypeName}}, nil
}

func (self *{{.typeName}}Model) FindMany(ctx context.Context, filter interface{}) ([]*{{.typeName}}, error) {
	{{.lowerTypeName}}s := []*{{.typeName}}{}
	err := self.client.Find(ctx, filter).All(&{{.lowerTypeName}}s)
	if err != nil {
		return nil, err
	}
	return {{.lowerTypeName}}s, nil
}

func (self *{{.typeName}}Model) InsertOne(ctx context.Context, {{.lowerTypeName}} *{{.typeName}}) error {
	_, err := self.client.InsertOne(ctx, {{.lowerTypeName}})
	return err
}

func (self *{{.typeName}}Model) InsertMany(ctx context.Context, {{.lowerTypeName}}s []*{{.typeName}}) error {
	_, err := self.client.InsertMany(ctx, {{.lowerTypeName}}s)
	return err
}

func (self *{{.typeName}}Model) Upsert(ctx context.Context, filter interface{}, {{.lowerTypeName}} *{{.typeName}}) error {
	_, err := self.client.Upsert(ctx, filter, {{.lowerTypeName}})
	return err
}

func (self *{{.typeName}}Model) UpdateOne(ctx context.Context, filter interface{}, {{.lowerTypeName}} *{{.typeName}}) error {
	return self.client.UpdateOne(ctx, filter, bson.M{"$set": {{.lowerTypeName}}})
}

func (self *{{.typeName}}Model) UpdateAll(ctx context.Context, filter interface{}, {{.lowerTypeName}} *{{.typeName}}) error {
	_, err := self.client.UpdateAll(ctx, filter, bson.M{"$set": {{.lowerTypeName}}})
	return err
}

func (self *{{.typeName}}Model) RemoveOne(ctx context.Context, filter interface{}) error {
	return self.client.Remove(ctx, filter)
}

func (self *{{.typeName}}Model) RemoveAll(ctx context.Context, filter interface{}) error {
	_, err := self.client.RemoveAll(ctx, filter)
	return err
}
`
