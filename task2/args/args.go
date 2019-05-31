package args

import (
	"errors"
	"strings"
)

var WrongSchemaRuleError = errors.New("can't create schema rule, wrong schema rule string")

type SchemaRule struct {
	flag        string
	typeCode    string
	defautValue string
}

func (sr *SchemaRule) getFlag() string {
	return sr.flag
}

func (sr *SchemaRule) getTypeCode() string {
	return sr.typeCode
}

func (sr *SchemaRule) getDefaultValue() string {
	if sr.defautValue != "" {
		return sr.defautValue
	}

	switch sr.getTypeCode() {
	case "bool":
		return "true"
	case "int":
		return "0"
	default:
		return ""
	}
}

func newSchemaRule(aSchemaRuleString string) (*SchemaRule, error) {
	sr := new(SchemaRule)
	splitSchemaRule := strings.Split(aSchemaRuleString, ":")
	switch len(splitSchemaRule) {
	case 2:
		sr.flag = splitSchemaRule[0]
		sr.typeCode = splitSchemaRule[1]
	case 3:
		sr.flag = splitSchemaRule[0]
		sr.typeCode = splitSchemaRule[1]
		sr.defautValue = splitSchemaRule[2]
	default:
		return nil, WrongSchemaRuleError
	}
	return sr, nil
}

type Schema struct {
	schemaRules map[string]*SchemaRule
}

func (s *Schema) getSchemaRule(flag string) *SchemaRule {
	return s.schemaRules[flag]
}

func (s *Schema) count() int {
	return len(s.schemaRules)
}

func newSchema(aSchemaString string) *Schema {
	aSchema := new(Schema)
	aSchema.schemaRules = make(map[string]*SchemaRule, 0)
	for _, schemaRuleString := range strings.Split(aSchemaString, " ") {
		aSchemaRule, _ := newSchemaRule(schemaRuleString)
		aSchema.schemaRules[aSchemaRule.getFlag()] = aSchemaRule
	}
	return aSchema
}
