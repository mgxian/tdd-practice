package args2

import (
	"errors"
	"strings"
)

var ErrWrongSchemaRule = errors.New("wrong shcemule rule")
var ErrNotSupportArgumentType = errors.New("not support argument type")

type SchemaRule struct {
	flag         string
	typeCode     string
	defaultValue string
}

func (sr *SchemaRule) getFlag() string {
	return sr.flag
}

func (sr *SchemaRule) getTypeCode() string {
	return sr.typeCode
}

func (sr *SchemaRule) getDefaultValue() string {
	return sr.defaultValue
}

func isSupportArgType(typeCode string) bool {
	switch typeCode {
	case "bool", "int", "string":
		return true
	default:
		return false
	}
}

func getDefaultValue(typeCode string) string {
	switch typeCode {
	case "bool":
		return "false"
	case "int":
		return "0"
	default:
		return ""
	}
}

func newSchemaRule(aSchemaRuleString string) (*SchemaRule, error) {
	srData := strings.Split(aSchemaRuleString, ":")
	if len(srData) > 3 || len(srData) < 2 {
		return nil, ErrWrongSchemaRule
	}

	flag := srData[0]
	typeCode := srData[1]
	if !isSupportArgType(typeCode) {
		return nil, ErrNotSupportArgumentType
	}

	defaultValue := ""
	if len(srData) == 3 {
		defaultValue = srData[2]
	}

	if defaultValue == "" {
		defaultValue = getDefaultValue(typeCode)
	}

	sr := new(SchemaRule)
	sr.flag = flag
	sr.typeCode = typeCode
	sr.defaultValue = defaultValue
	return sr, nil
}
