package model

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
)

const indexPri = "PRIMARY"

type (
	// InformationSchemaModel defines information schema model
	InformationSchemaModel struct {
		db *gorm.DB
	}

	Table struct {
		TableName string `gorm:"column:TABLE_NAME"`
	}

	Column struct {
		Name            string      `gorm:"column:COLUMN_NAME"`
		DataType        string      `gorm:"column:DATA_TYPE"`
		Extra           string      `gorm:"column:EXTRA"`
		Comment         string      `gorm:"column:COLUMN_COMMENT"`
		ColumnDefault   interface{} `gorm:"column:COLUMN_DEFAULT"`
		IsNullAble      string      `gorm:"column:IS_NULLABLE"`
		OrdinalPosition int         `gorm:"column:ORDINAL_POSITION"`
		ColumnType      string      `gorm:"column:COLUMN_TYPE"`
	}

	Statistic struct {
		IndexName  string `gorm:"column:INDEX_NAME"`
		NonUnique  int    `gorm:"column:NON_UNIQUE"`
		SeqInIndex int    `gorm:"column:SEQ_IN_INDEX"`
	}

	ColumnData struct {
		*Column
		*Statistic
	}

	// Column defines column in table
	// Column struct {
	// 	*DbColumn
	// 	Index *DbIndex
	// }

	TableColumnsData struct {
		Db          string
		Table       string
		ColumnDatas []*ColumnData
	}

	// 根据数据库信息组合而成的完整table信息
	TableData struct {
		Db          string
		Table       string
		ColumnDatas []*ColumnData
		// Primary key not included
		PrimaryKey  *ColumnData
		UniqueIndex map[string][]*ColumnData
		NormalIndex map[string][]*ColumnData
	}

	// IndexType describes an alias of string
	IndexType string

	// Index describes a column index
	Index struct {
		IndexType IndexType
		Columns   []*Column
	}
)

// NewInformationSchemaModel creates an instance for InformationSchemaModel
func NewInformationSchemaModel(db *gorm.DB) *InformationSchemaModel {
	return &InformationSchemaModel{db: db}
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *InformationSchemaModel) GetAllTables(database string) ([]*Table, error) {
	var tables []*Table
	if err := m.db.Where("TABLE_SCHEMA = ?", database).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

// 根据db、table名获取该table的所有列信息
func (m *InformationSchemaModel) FindTableColumnsData(db, table string) (*TableColumnsData, error) {
	var columns []*Column
	err := m.db.Table("columns").Where("TABLE_SCHEMA = ? and TABLE_NAME = ?", db, table).Find(&columns).Error
	if err != nil {
		return nil, err
	}

	tableColumnsData := &TableColumnsData{
		Db:    db,
		Table: table,
	}
	var list []*Column
	for _, column := range columns {
		statistics, err := m.FindStatistics(db, table, column.Name)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
			continue
		}

		if len(statistics) > 0 {
			for _, statistic := range statistics {
				tableColumnsData.ColumnDatas = append(tableColumnsData.ColumnDatas, &ColumnData{
					Column:    column,
					Statistic: statistic,
				})
			}
		} else {
			tableColumnsData.ColumnDatas = append(tableColumnsData.ColumnDatas, &ColumnData{
				Column: column,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].OrdinalPosition < list[j].OrdinalPosition
	})

	return tableColumnsData, nil
}

// FindIndex finds index with given db, table and column.
func (m *InformationSchemaModel) FindStatistics(db, table, column string) ([]*Statistic, error) {
	var statistics []*Statistic
	err := m.db.Where("TABLE_SCHEMA = ? and TABLE_NAME = ? and COLUMN_NAME = ?", db, table, column).Find(&statistics).Error
	if err != nil {
		return nil, err
	}
	return statistics, nil
}

// 将一个table的列数据对象转换为一个table
func (c *TableColumnsData) Convert() (*TableData, error) {
	tableData := &TableData{
		Table:       c.Table,
		Db:          c.Db,
		ColumnDatas: c.ColumnDatas,
		UniqueIndex: map[string][]*ColumnData{},
		NormalIndex: map[string][]*ColumnData{},
	}

	m := make(map[string][]*ColumnData)
	for _, columnData := range c.ColumnDatas {
		if columnData.Statistic != nil {
			m[columnData.IndexName] = append(m[columnData.IndexName], columnData)
		}
	}

	primaryColumns := m[indexPri]
	if len(primaryColumns) == 0 {
		return nil, fmt.Errorf("db:%s, table:%s, missing primary key", c.Db, c.Table)
	}

	// todo: 暂时删除对联合主键的限制
	// if len(primaryColumns) > 1 {
	// 	return nil, fmt.Errorf("db:%s, table:%s, joint primary key is not supported", c.Db, c.Table)
	// }

	tableData.PrimaryKey = primaryColumns[0]
	for indexName, ColumnDatas := range m {
		if indexName == indexPri {
			continue
		}

		for _, columnData := range ColumnDatas {
			if columnData.Statistic != nil {
				if columnData.NonUnique == 0 {
					tableData.UniqueIndex[indexName] = ColumnDatas
				} else {
					tableData.NormalIndex[indexName] = ColumnDatas
				}
			}
		}
	}

	return tableData, nil
}
