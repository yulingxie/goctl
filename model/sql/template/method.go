package template

var Methods = `
// 查找
func (self *{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	resp := &{{.upperStartCamelObject}}{}
	err := self.client.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).First(&resp).Error
	return resp, err
}

// 条件查找
func (self *{{.upperStartCamelObject}}Model) FindOneByConds(ctx context.Context, conds ...interface{}) (*{{.upperStartCamelObject}}, error) {
	resp := &{{.upperStartCamelObject}}{}
	err := self.client.WithContext(ctx).First(&resp, conds...).Error
	return resp, err
}

// 批量查找
func (self *{{.upperStartCamelObject}}Model) BatchFind(ctx context.Context, {{.lowerStartCamelPrimaryKey}}s []{{.dataType}}) ([]*{{.upperStartCamelObject}}, error) {
	resp := []*{{.upperStartCamelObject}}{}
	err := self.client.WithContext(ctx).Where("{{.originalPrimaryKey}} IN ?", {{.lowerStartCamelPrimaryKey}}s).Find(&resp).Error
	return resp, err
}

// 批量条件查找
func (self *{{.upperStartCamelObject}}Model) BatchFindByConds(ctx context.Context, conds ...interface{}) ([]*{{.upperStartCamelObject}}, error) {
	resp := []*{{.upperStartCamelObject}}{}
	err := self.client.WithContext(ctx).Find(&resp, conds...).Error
	return resp, err
}

// 只更新非零的filed
func (self *{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{.lowerStartCamelObject}} *{{.upperStartCamelObject}}) (int64, error) {
	result := self.client.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}}).Updates({{.lowerStartCamelObject}})
	return result.RowsAffected, result.Error
}

// 更新指定fileds
func (self *{{.upperStartCamelObject}}Model) UpdateColumns(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, columns map[string]interface{}) (int64, error) {
	result := self.client.WithContext(ctx).Table({{.lowerStartCamelObject}}TableName).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).UpdateColumns(columns)
	return result.RowsAffected, result.Error
}

// 更新所有field, 包括零值field, 并且自带upsert功能
func (self *{{.upperStartCamelObject}}Model) Save(ctx context.Context, {{.lowerStartCamelObject}} *{{.upperStartCamelObject}}) (int64, error) {
	result := self.client.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}}).Save({{.lowerStartCamelObject}})
	return result.RowsAffected, result.Error
}

// 插入
func (self *{{.upperStartCamelObject}}Model) Insert(ctx context.Context, {{.lowerStartCamelObject}} *{{.upperStartCamelObject}}) (int64,error) {
	result := self.client.WithContext(ctx).Create({{.lowerStartCamelObject}})
	return result.RowsAffected, result.Error
}

// 删除
func (self *{{.upperStartCamelObject}}Model) Delete(ctx context.Context, conds ...interface{}) (int64, error) {
	result := self.client.WithContext(ctx).Delete(&{{.upperStartCamelObject}}{}, conds...)
	return result.RowsAffected, result.Error
}
`

var MethodsWithCache = `
func (self *{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	resp := &{{.upperStartCamelObject}}{}
	cacheKey := fmt.Sprintf(self.cacheKey, {{.lowerStartCamelPrimaryKey}})
	err := self.cache.Take(resp, cacheKey, func(v interface{}) error {
		return self.client.WithContext(ctx).Table("userwealthbank").Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).First(&v).Error
	})
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (self *{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{.lowerStartCamelObject}} *{{.upperStartCamelObject}}) (int64, error) {
	cacheKey := fmt.Sprintf(self.cacheKey, {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}})
	result := self.client.WithContext(ctx).Table({{.lowerStartCamelObject}}TableName).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}}).Updates({{.lowerStartCamelObject}})
	if result.RowsAffected != 0 {
		if err := self.cache.Del(cacheKey); err != nil {
			return result.RowsAffected, err
		}
	}
	return result.RowsAffected, result.Error
}

func (self *{{.upperStartCamelObject}}Model) Save(ctx context.Context, {{.lowerStartCamelObject}} *{{.upperStartCamelObject}}) (int64, error) {
	cacheKey := fmt.Sprintf(self.cacheKey, {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}})
	result := self.client.WithContext(ctx).Table({{.lowerStartCamelObject}}TableName).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelObject}}.{{.upperStartCamelPrimaryKey}}).Save({{.lowerStartCamelObject}})
	if result.RowsAffected != 0 {
		if err := self.cache.Del(cacheKey); err != nil {
			return result.RowsAffected, err
		}
	}
	return result.RowsAffected, result.Error
}
`
