package args

import (
	"reflect"
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

func TestParse(t *testing.T) {
	schemaString := "l:bool:false p:int:80 d:string:./logs"
	argsString := "-l -p 8080 -d /usr/logs"
	fullArgsString := "-l true -p 8080 -d /usr/logs"
	flagTests := []struct {
		flag  string
		value string
	}{
		{"l", "true"},
		{"p", "8080"},
		{"d", "/usr/logs"},
	}
	aSchema := newSchema(schemaString)

	t.Run("return number of args for full args string", func(t *testing.T) {
		aSchema.Parse(fullArgsString)
		got := aSchema.Size()
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("return string value of full args string", func(t *testing.T) {
		aSchema.Parse(fullArgsString)
		for _, tt := range flagTests {
			got := aSchema.GetArg(tt.flag)
			if got != tt.value {
				t.Errorf("got '%s', want '%s'", got, tt.value)
			}
		}
	})

	t.Run("return number of args for args string", func(t *testing.T) {
		aSchema.Parse(argsString)
		got := aSchema.Size()
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("return string value of args string", func(t *testing.T) {
		aSchema.Parse(argsString)
		for _, tt := range flagTests {
			got := aSchema.GetArg(tt.flag)
			if got != tt.value {
				t.Errorf("got '%s', want '%s'", got, tt.value)
			}
		}
	})
}
