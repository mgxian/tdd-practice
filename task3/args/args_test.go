package args

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSchemaRule(t *testing.T) {
	assertSchemaRule := func(t *testing.T, schemaRuleString string, wantedSchemaRule interface{}) {
		t.Helper()
		got := newSchemaRule(schemaRuleString)
		if !reflect.DeepEqual(got, wantedSchemaRule) {
			t.Errorf("got %v, want %v", got, wantedSchemaRule)
		}
	}

	t.Run("return bool schema rule", func(t *testing.T) {
		aBoolSchemaRuleString := "l:bool:false"
		wantedSchemaRule := boolSchemaRule{baseSchemaRule{"l", "bool", "false"}}
		assertSchemaRule(t, aBoolSchemaRuleString, wantedSchemaRule)
	})

	t.Run("return int schema rule", func(t *testing.T) {
		aIntSchemaRuleString := "p:int:80"
		wantedSchemaRule := intSchemaRule{baseSchemaRule{"p", "int", "80"}}
		assertSchemaRule(t, aIntSchemaRuleString, wantedSchemaRule)
	})

	t.Run("return string schema rule", func(t *testing.T) {
		aStringSchemaRuleString := "d:string:./logs"
		wantedSchemaRule := stringSchemaRule{baseSchemaRule{"d", "string", "./logs"}}
		assertSchemaRule(t, aStringSchemaRuleString, wantedSchemaRule)
	})
}

func TestSchema(t *testing.T) {
	schemaString := "l:bool:false p:int:80 d:string:./logs"
	aSchema := newSchema(schemaString)
	wantedSchemaRules := map[string]SchemaRule{
		"l": boolSchemaRule{baseSchemaRule{"l", "bool", "false"}},
		"p": intSchemaRule{baseSchemaRule{"p", "int", "80"}},
		"d": stringSchemaRule{baseSchemaRule{"d", "string", "./logs"}},
	}
	for flag, sr := range wantedSchemaRules {
		got := aSchema.getSchemaRule(flag)
		if !reflect.DeepEqual(got, sr) {
			t.Errorf("got %v, want %v", got, sr)
		}
	}
}

type flagTest struct {
	flag     string
	value    string
	typeCode string
}

func TestParse(t *testing.T) {
	schemaString := "l:bool:false p:int:80 d:string:./logs"
	argsString := "-l -p 8080 -d /usr/logs"
	fullArgsString := "-l true -p 8080 -d /usr/logs"
	partialArgsString := "-l -p 8080"
	flagTests := []flagTest{
		{"l", "true", "bool"},
		{"p", "8080", "int"},
		{"d", "/usr/logs", "string"},
	}

	partialFlagTests := []flagTest{
		{"l", "true", "bool"},
		{"p", "8080", "int"},
		{"d", "./logs", "string"},
	}

	aSchema := newSchema(schemaString)

	t.Run("return number of args for full args string", func(t *testing.T) {
		if err := aSchema.Parse(fullArgsString); err != nil {
			t.Errorf(err.Error())
		}
		got := aSchema.Size()
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("return string value of full args string", func(t *testing.T) {
		if err := aSchema.Parse(fullArgsString); err != nil {
			t.Errorf(err.Error())
		}
		for _, tt := range flagTests {
			got := aSchema.getArg(tt.flag)
			assertArgValueString(t, got, tt.value)
		}
	})

	t.Run("return number of args for args string", func(t *testing.T) {
		if err := aSchema.Parse(argsString); err != nil {
			t.Errorf(err.Error())
		}
		got := aSchema.Size()
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("return string value of args string", func(t *testing.T) {
		if err := aSchema.Parse(argsString); err != nil {
			t.Errorf(err.Error())
		}
		for _, tt := range flagTests {
			got := aSchema.getArg(tt.flag)
			assertArgValueString(t, got, tt.value)
		}
	})

	name := "return right value of args string"
	testParseArgs(t, name, aSchema, argsString, flagTests)

	name = "return right value of partial args string"
	testParseArgs(t, name, aSchema, partialArgsString, partialFlagTests)
}

func testParseArgs(t *testing.T, name string, aSchema *Schema, argString string, tests []flagTest) {
	t.Run(name, func(t *testing.T) {
		if err := aSchema.Parse(argString); err != nil {
			t.Errorf(err.Error())
		}
		for _, tt := range tests {
			switch tt.typeCode {
			case "bool":
				assertArgBoolValue(t, aSchema, tt.flag, true)
			case "int":
				assertArgIntValue(t, aSchema, tt.flag, tt.value)
			case "string":
				argV := aSchema.GetStringArg(tt.flag)
				assertArgValueString(t, argV, tt.value)
			default:
				t.Errorf("unkown arg type")
			}
		}
	})
}

func assertArgBoolValue(t *testing.T, aSchema *Schema, flag string, want bool) {
	argV := aSchema.GetBoolArg(flag)
	if argV != want {
		t.Errorf("got %v, want %v", argV, want)
	}
}

func assertArgIntValue(t *testing.T, aSchema *Schema, flag, value string) {
	argV, _ := aSchema.GetIntArg(flag)
	want, _ := strconv.Atoi(value)
	if argV != want {
		t.Errorf("got %v, want %v", argV, want)
	}
}

func assertArgValueString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
