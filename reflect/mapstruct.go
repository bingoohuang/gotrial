package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"reflect"
	"time"
)

// Map2Struct map m to struct by its type
func Map2Struct(m map[string]interface{}, result interface{}) error {
	structTypePtr := reflect.TypeOf(result)
	if structTypePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}
	v := reflect.ValueOf(result).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}

	structFields := CachedStructFields(v.Type(), "iql")
	for i, sf := range structFields {
		fillField(m, sf, v.Field(i))
	}

	return nil
}

func fillField(m map[string]interface{}, sf StructField, f reflect.Value) {
	for _, tag := range sf.Tag {
		if v, ok := m[tag]; ok {
			setFieldValue(sf, f, v)
			return
		}
	}

	name := strcase.ToSnake(sf.Name)
	if v, ok := m[name]; ok {
		setFieldValue(sf, f, v)
	}
}

var timeKind = reflect.TypeOf(time.Time{}).Kind()

func setFieldValue(sf StructField, f reflect.Value, v interface{}) {
	fk := sf.Kind
	vk := reflect.TypeOf(v).Kind()

	if fk == vk {
		f.Set(reflect.ValueOf(v))
	} else if fk == timeKind && vk == reflect.String {
		t, _ := time.Parse(time.RFC3339, v.(string))
		f.Set(reflect.ValueOf(t))
	} else if fk == reflect.Float32 && vk == reflect.Float64 {
		f.SetFloat(v.(float64))
	} else if fk == reflect.String {
		f.SetString(fmt.Sprintf("%v", v))
	} else {
		fmt.Printf("fk:%v, vk:%v\n", fk, vk)
	}
}
