package model

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	model *InformationSchemaModel
)

func init() {
	db, _ := gorm.Open(mysql.Open("k7game-server:jJ8VkXERg83D7z44@tcp(tidb-1001.db.qipai007cs.com:3306)/information_schema"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	model = NewInformationSchemaModel(db)
}

func TestInformationSchemaModel_GetAllTables(t *testing.T) {
	got, _ := model.GetAllTables("yygplatform")
	for _, table := range got {
		t.Logf("%v", table.TableName)
	}
}
