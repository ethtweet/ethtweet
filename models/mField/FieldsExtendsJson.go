package mField

import (
	"encoding/json"

	"github.com/ethtweet/ethtweet/global"

	"github.com/tidwall/gjson"
)

type FieldsExtendsJsonType struct {
	ExtraData string `gorm:"type:text;comment:扩展字段"`
}

func (e *FieldsExtendsJsonType) GetExtendsJson(key string) gjson.Result {
	return gjson.Get(e.ExtraData, key)
}

func (e *FieldsExtendsJsonType) SetExtendsJson(key string, value interface{}) {
	r := global.Json2Map(e.ExtraData)
	r[key] = value
	eJson, _ := json.Marshal(r)
	e.ExtraData = string(eJson)
}
