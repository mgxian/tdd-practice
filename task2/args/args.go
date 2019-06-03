package args

import (
	"errors"
	"strconv"
	"strings"
)

var WrongSchemaRuleError = errors.New("can't create schema rule, wrong schema rule string")
var FlagNotExistError = errors.New("not found such flag, flag not exist")
var ArgValueError = errors.New("argument value error")

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

func (s *Schema) getSchemaRule(flag string) (*SchemaRule, error) {
	if sr, ok := s.schemaRules[flag]; ok {
		return sr, nil
	}
	return nil, FlagNotExistError
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

type Parser struct {
	schema   *Schema
	argPairs map[string]string
}

func newParser(aSchemaString string) *Parser {
	parser := new(Parser)
	parser.schema = newSchema(aSchemaString)
	parser.argPairs = make(map[string]string, 0)
	return parser
}

func (p *Parser) parse(aArgString string) error {
	args := strings.Split(aArgString, " ")
	for i := 0; i < len(args); {
		flag := args[i][1:]
		sr, err := p.schema.getSchemaRule(flag)
		if err != nil {
			return FlagNotExistError
		}
		step := 2
		value := args[i+1]
		if sr.getTypeCode() == "bool" && value != "true" {
			step = 1
			value = "true"
		}
		p.argPairs[flag] = value
		i += step
	}
	return nil
}

func (p *Parser) GetStringArg(flag string) string {
	return p.argPairs[flag]
}

func (p *Parser) GetBoolArg(flag string) bool {
	argv := p.GetStringArg(flag)
	if argv == "true" {
		return true
	}
	return false
}

func (p *Parser) GetIntArg(flag string) int {
	argv := p.GetStringArg(flag)
	if v, err := strconv.Atoi(argv); err == nil {
		return v
	}
	return 0
}
