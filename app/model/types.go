package model

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// DatabaseMap 自定义类型，用于存储 map[string]any 到 JSON 字符串
type DatabaseMap map[string]any

// Value 实现了 sql.Scanner 接口，用于从数据库中扫描值
func (m DatabaseMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现了 sql.Scanner 接口，用于扫描数据库中的值
func (m *DatabaseMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(DatabaseMap)
		return nil
	}

	// 确保传入的是一个可扫描的类型
	b, ok := value.([]byte)
	if !ok {
		return gorm.ErrInvalidValue
	}

	// 解码 JSON 到 map
	return json.Unmarshal(b, m)
}

// DataStrings 自定义类型，用于存储字符串切片到 JSON 字符串
type DataStrings []string

// Value 实现了 gorm.Valuer 接口，用于将自定义类型转换为数据库值
func (s DataStrings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现了 sql.Scanner 接口，用于将数据库值扫描到自定义类型
func (s *DataStrings) Scan(value interface{}) error {
	if value == nil {
		*s = make(DataStrings, 0)
		return nil
	}

	// 确保传入的是一个可扫描的类型
	b, ok := value.([]byte)
	if !ok {
		return gorm.ErrInvalidValue
	}

	// 解码 JSON 到字符串切片
	return json.Unmarshal(b, s)
}
