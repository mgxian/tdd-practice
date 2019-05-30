package args

import "strings"

type Schema struct {
	schemaRules map[string]SchemaRule
	argPairs    map[string]string
}

func (s *Schema) getSchemaRule(flag string) SchemaRule {
	return s.schemaRules[flag]
}

func (s *Schema) Parse(argsString string) {
	argPair := ""
	s.argPairs = make(map[string]string, 0)
	for i, ss := range strings.Split(argsString, " ") {
		argPair += ss
		if i%2 != 0 {
			flag := argPair[1]
			value := argPair[2:]
			s.argPairs[string(flag)] = value
			argPair = ""
		}
	}
}

func (s *Schema) Size() int {
	return len(s.argPairs)
}

func (s *Schema) GetArg(flag string) string {
	return s.argPairs[flag]
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
