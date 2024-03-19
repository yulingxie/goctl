package converter

import (
	"fmt"
	"strings"
)

var commonMysqlDataTypeMap = map[string]string{
	// For consistency, all integer types are converted to int64
	// number
	"bool":      "int64",
	"boolean":   "int64",
	"tinyint":   "int64",
	"smallint":  "int64",
	"mediumint": "int64",
	"int":       "int64",
	"integer":   "int64",
	"bigint":    "int64",
	"float":     "float64",
	"double":    "float64",
	"decimal":   "float64",
	// date&time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "string",
	"year":      "int64",
	// string
	"char":       "string",
	"varchar":    "string",
	"binary":     "string",
	"varbinary":  "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"enum":       "string",
	"set":        "string",
	"json":       "string",
}

func getDataType(dataBaseType string, isUnsigned bool) string {
	switch strings.ToLower(dataBaseType) {
	case "bool", "boolean":
		return "bool"
	case "tinyint":
		if isUnsigned {
			return "uint8"
		} else {
			return "int8"
		}
	case "smallint":
		if isUnsigned {
			return "uint16"
		} else {
			return "int16"
		}
	case "mediumint", "int", "integer":
		if isUnsigned {
			return "uint32"
		} else {
			return "int32"
		}
	case "bigint":
		if isUnsigned {
			return "uint64"
		} else {
			return "int64"
		}
	case "float":
		return "float32"
	case "double", "decimal":
		return "float64"
	case "date", "datetime", "timestamp":
		return "time.Time"
	case "year":
		return "uint8"
	case "time":
		return "string"
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext", "enum", "set", "json":
		return "string"
	case "blob":
		return "[]byte"
	default:
		return ""
	}
}

func ConvertDataType(dataBaseType string, columnType string, isDefaultNull bool) (string, error) {
	isUnsigned := strings.Contains(columnType, "unsigned")
	tp := getDataType(dataBaseType, isUnsigned)
	if len(tp) == 0 {
		return "", fmt.Errorf("unexpected database type: %s", dataBaseType)
	}
	return mayConvertNullType(tp, isDefaultNull), nil
}

func ConvertDataType2(dataBaseType string, isUnsigned, isDefaultNull bool) (string, error) {
	tp := getDataType(dataBaseType, isUnsigned)
	if len(tp) == 0 {
		return "", fmt.Errorf("unexpected database type: %s", dataBaseType)
	}
	return mayConvertNullType(tp, isDefaultNull), nil
}

func mayConvertNullType(goDataType string, isDefaultNull bool) string {
	if !isDefaultNull {
		return goDataType
	}

	switch goDataType {
	case "int64":
		return "sql.NullInt64"
	case "int32":
		return "sql.NullInt32"
	case "float64":
		return "sql.NullFloat64"
	case "bool":
		return "sql.NullBool"
	case "string":
		return "sql.NullString"
	case "time.Time":
		return "sql.NullTime"
	default:
		return goDataType
	}
}
