package args2

import (
	"errors"
	"strconv"
	"strings"
)

var ErrWrongSchemaRule = errors.New("wrong shcemule rule")
var ErrNotSupportArgumentType = errors.New("not support argument type")
var ErrorFlagNotExist = errors.New("flag not exist")

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

type Schema struct {
	schemaRules map[string]*SchemaRule
}

func newSchema(aSchemaString string) (*Schema, error) {
	aSchema := new(Schema)
	aSchema.schemaRules = make(map[string]*SchemaRule, 0)
	schemaData := strings.Split(aSchemaString, " ")
	for _, sd := range schemaData {
		sr, err := newSchemaRule(sd)
		if err != nil {
			return nil, err
		}
		aSchema.schemaRules[sr.getFlag()] = sr
	}
	return aSchema, nil
}

func (s *Schema) size() int {
	return len(s.schemaRules)
}

func (s *Schema) typeOf(flag string) (string, error) {
	if sr, ok := s.schemaRules[flag]; ok {
		return sr.getTypeCode(), nil
	}
	return "", ErrorFlagNotExist
}

func (s *Schema) defaultValueOf(flag string) (string, error) {
	if sr, ok := s.schemaRules[flag]; ok {
		return sr.getDefaultValue(), nil
	}
	return "", ErrorFlagNotExist
}

type Parser struct {
	schema    *Schema
	arguments map[string]string
}

func newParser(aSchemaString string) (*Parser, error) {
	aSchema, err := newSchema(aSchemaString)
	if err != nil {
		return nil, err
	}
	aParser := new(Parser)
	aParser.schema = aSchema
	aParser.arguments = make(map[string]string, 0)
	return aParser, nil
}

func (p *Parser) parse(aArgumentsString string) error {
	p.arguments = make(map[string]string, 0)
	argumentsData := strings.Split(aArgumentsString, " ")
	for i := 0; i < len(argumentsData); {
		flag := argumentsData[i][1:]
		typeCode, err := p.schema.typeOf(flag)
		if err != nil {
			return err
		}

		step := 2
		value := ""
		if typeCode == "bool" && i+1 < len(argumentsData) && argumentsData[i+1] != "true" {
			step = 1
			value = "true"
		} else {
			value = argumentsData[i+1]
		}
		i += step
		p.arguments[flag] = value
	}
	return nil
}

func (p *Parser) stringValueOf(flag string) (string, error) {
	if v, ok := p.arguments[flag]; ok {
		return v, nil
	}

	dv, err := p.schema.defaultValueOf(flag)
	if err == nil {
		return dv, nil
	}
	return "", err
}

func (p *Parser) boolValueOf(flag string) (bool, error) {
	v, err := p.stringValueOf(flag)
	if err != nil {
		return false, err
	}

	if v == "true" {
		return true, nil
	}

	return false, nil
}

func (p *Parser) intValueOf(flag string) (int, error) {
	v, err := p.stringValueOf(flag)
	if err != nil {
		return 0, err
	}

	intv, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return intv, nil
}
