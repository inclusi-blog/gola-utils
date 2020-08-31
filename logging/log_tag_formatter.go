package logging

import (
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/sirupsen/logrus"
	"reflect"
)

const (
	LogTag = "log"
)

type LogTagFormatter struct {
	defaultFormatter logrus.Formatter
	options          map[string]interface{}
}

func (f *LogTagFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formattedEntry := *entry

	for k, v := range entry.Data {
		entryDataType := reflect.TypeOf(v)
		entryDataValue := reflect.ValueOf(v)
		if entryDataType.Kind() == reflect.Ptr {
			entryDataType = entryDataType.Elem()
		}
		if entryDataType.Kind() == reflect.Struct {
			mp := reflectx.NewMapperFunc(LogTag, func(s string) string { return "" })

			root := mp.TypeMap(entryDataType).Tree

			formattedStruct := traverseFieldTree(entryDataValue, root)
			formattedEntry.Data[k] = formattedStruct
		}
	}

	if f.defaultFormatter == nil {
		f.defaultFormatter = &logrus.JSONFormatter{}
	}
	return f.defaultFormatter.Format(&formattedEntry)
}

func traverseFieldTree(val reflect.Value, root *reflectx.FieldInfo) map[string]interface{} {
	if root == nil {
		return nil
	}
	var childMap = map[string]interface{}{}
	for _, child := range root.Children {
		if child == nil || child.Name == "" {
			continue
		}
		if len(child.Children) == 0 {
			childMap[child.Name] = reflectx.FieldByIndexes(val, child.Index).Interface()
			continue
		}
		childMap[child.Name] = traverseFieldTree(val, child)
	}
	return childMap
}
