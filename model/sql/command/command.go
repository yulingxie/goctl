package command

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/gen"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/urfave/cli"
)

var errNotMatched = errors.New("sql not matched")

func CreateSqlModel(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String("url"))
	node := strings.TrimSpace(ctx.String("node"))
	db := strings.TrimSpace(ctx.String("db"))
	dir := strings.TrimSpace(ctx.String("dir"))
	cache := ctx.Bool("cache")
	tablePattern := strings.TrimSpace(ctx.String("table"))
	return fromDB(url, node, db, tablePattern, dir, cache)
}

func fromDB(url, nodeName, dbName, tablePattern, dir string, cache bool) error {
	if len(url) == 0 {
		// 默认使用内网的数据库
		url = "k7game-server:jJ8VkXERg83D7z44@tcp(tidb-1001.db.qipai007cs.com:3306)"
	}

	if len(dir) == 0 {
		dir = "./"
	}

	if len(nodeName) == 0 {
		nodeName = "common"
	}

	if len(dbName) == 0 {
		fmt.Printf("%v", "未指定dbName")
		return nil
	}

	if len(tablePattern) == 0 {
		fmt.Printf("%v", "未指定table")
		return nil
	}

	databaseSource := url + "/information_schema"
	db, err := gorm.Open(mysql.Open(databaseSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	if err != nil {
		return nil
	}
	informationSchemaModel := model.NewInformationSchemaModel(db)

	tables, err := informationSchemaModel.GetAllTables(dbName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.TableData)
	for _, item := range tables {
		match, err := filepath.Match(tablePattern, item.TableName)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := informationSchemaModel.FindTableColumnsData(dbName, item.TableName)
		if err != nil {
			return err
		}

		tableData, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item.TableName] = tableData
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, nodeName)
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, cache)
}
