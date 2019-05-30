package args

import (
	"errors"
	"strconv"
	"strings"
)

type Schema struct {
	schemaRules map[string]SchemaRule
	argPairs    map[string]string
}

func (s *Schema) getSchemaRule(flag string) SchemaRule {
	return s.schemaRules[flag]
}

func (s *Schema) Parse(argsString string) error {
	s.argPairs = make(map[string]string, 0)
	splitArgs := strings.Split(argsString, " ")
	for i := 0; i < len(splitArgs); {
		flag := splitArgs[i][1:]
		sr, ok := s.schemaRules[flag]
		if !ok {
			return errors.New("unknown arg: " + flag)
		}
		typeCode := sr.getTypeCode()
		if typeCode == "bool" {
			s.argPairs[flag] = "true"
			i++
			if splitArgs[i][0] != '-' {
				i++
			}
			continue
		}
		s.argPairs[flag] = splitArgs[i+1]
		i += 2
	}
	return nil
}

func (s *Schema) Size() int {
	return len(s.argPairs)
}

func (s *Schema) getArg(flag string) string {
	v, ok := s.argPairs[flag]
	if !ok {
		return s.schemaRules[flag].getValue()
	}
	return v
}

func (s *Schema) GetBoolArg(flag string) bool {
	if s.getArg(flag) == "true" {
		return true
	}
	return false
}

func (s *Schema) GetIntArg(flag string) (int, error) {
	v, err := strconv.Atoi(s.getArg(flag))
	if err != nil {
		return 0, errors.New("type error")
	}
	return v, err
}

func (s *Schema) GetStringArg(flag string) string {
	return s.getArg(flag)
}

func newSchema(schemaString string) *Schema {
	aSchema := new(Schema)
	aSchema.schemaRules = make(map[string]SchemaRule, 0)
	for _, s := range strings.Split(schemaString, " ") {
		sr := newSchemaRule(s)
		aSchema.schemaRules[sr.getFlag()] = sr
	}
	return aSchema
}

type SchemaRule interface {
	getFlag() string
	getValue() string
	getTypeCode() string
}

type baseSchemaRule struct {
	flag         string
	typeCode     string
	defaultValue string
}

func (bsr baseSchemaRule) getFlag() string {
	return bsr.flag
}

func (bsr baseSchemaRule) getValue() string {
	return bsr.defaultValue
}

func (bsr baseSchemaRule) getTypeCode() string {
	return bsr.typeCode
}

type boolSchemaRule struct {
	baseSchemaRule
}

type intSchemaRule struct {
	baseSchemaRule
}

type stringSchemaRule struct {
	baseSchemaRule
}

func newSchemaRule(s string) SchemaRule {
	schemaData := strings.Split(s, ":")
	typeCode := schemaData[1]
	switch typeCode {
	case "bool":
		return boolSchemaRule{
			baseSchemaRule{
				flag:         schemaData[0],
				typeCode:     schemaData[1],
				defaultValue: schemaData[2],
			},
		}
	case "int":
		return intSchemaRule{
			baseSchemaRule{
				flag:         schemaData[0],
				typeCode:     schemaData[1],
				defaultValue: schemaData[2],
			},
		}
	case "string":
		return stringSchemaRule{
			baseSchemaRule{
				flag:         schemaData[0],
				typeCode:     schemaData[1],
				defaultValue: schemaData[2],
			},
		}
	default:
		return nil
	}
}
