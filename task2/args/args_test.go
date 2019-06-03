package args

import (
	"reflect"
	"testing"
)

type schemaRuleTest struct {
	schemaRuleString string
	flag             string
	typeCode         string
	defaultValue     string
	err              error
}

func TestSchemaRule(t *testing.T) {
	schemaRuleTests := []schemaRuleTest{
		{"l:bool:true", "l", "bool", "true", nil},
		{"p:int:80", "p", "int", "80", nil},
		{"p:int:", "p", "int", "0", nil},
		{"p:int", "p", "int", "0", nil},
		{"p", "", "", "", WrongSchemaRuleError},
	}

	for _, tt := range schemaRuleTests {
		sr, err := newSchemaRule(tt.schemaRuleString)
		assertNil(t, sr)
		if tt.err != nil {
			assertError(t, err, tt.err)
			continue
		}
		assertNoError(t, err)
		assertSchemaRule(t, tt, sr)
	}
}

func TestSchema(t *testing.T) {
	schemaTests := []struct {
		schemaString    string
		schemaRuleTests []schemaRuleTest
		count           int
	}{
		{
			"l:bool:true p:int:80 d:string:/usr/logs",
			[]schemaRuleTest{
				{"l:bool:true", "l", "bool", "true", nil},
				{"p:int:80", "p", "int", "80", nil},
				{"d:string:/usr/logs", "d", "string", "/usr/logs", nil},
			},
			3,
		},
		{
			"l:bool:true p:int:80",
			[]schemaRuleTest{
				{"l:bool:true", "l", "bool", "true", nil},
				{"p:int:80", "p", "int", "80", nil},
				{"d:string:/usr/logs", "d", "", "", FlagNotExistError},
			},
			2,
		},
	}

	for _, tt := range schemaTests {
		aSchema := newSchema(tt.schemaString)
		assertNil(t, aSchema)
		wantSchemaRuleCount := tt.count
		got := aSchema.count()
		if got != wantSchemaRuleCount {
			t.Errorf("got %d, want %d", got, wantSchemaRuleCount)
		}
		testSchemaRules(t, aSchema, tt.schemaRuleTests)
	}
}

type flagTest struct {
	flag     string
	typeCode string
	value    interface{}
}

func TestParser(t *testing.T) {
	aSchemaString := "l:bool:false p:int:80 d:string:./logs"
	t.Run("test full arg pair", func(t *testing.T) {
		argString := "-l true -p 8080 -d /usr/logs"
		flagTests := []flagTest{
			{"l", "bool", true},
			{"p", "int", 8080},
			{"d", "string", "/usr/logs"},
		}
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		aParser.parse(argString)
		testGetArgValue(t, aParser, flagTests)
	})

	t.Run("test common arg pair", func(t *testing.T) {
		argString := "-l -p 8080 -d /usr/logs"
		flagTests := []flagTest{
			{"l", "bool", true},
			{"p", "int", 8080},
			{"d", "string", "/usr/logs"},
		}
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		aParser.parse(argString)
		testGetArgValue(t, aParser, flagTests)
	})

	t.Run("test default arg pair", func(t *testing.T) {
		argString := "-d /usr/logs"
		flagTests := []flagTest{
			{"l", "bool", false},
			{"p", "int", 80},
			{"d", "string", "/usr/logs"},
		}
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		aParser.parse(argString)
		testGetArgValue(t, aParser, flagTests)
	})

	t.Run("test not exist arg pair", func(t *testing.T) {
		argString := "-h -p 8080 -d /usr/logs"
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		err := aParser.parse(argString)
		assertError(t, err, FlagNotExistError)
	})

	t.Run("test int list arg pair", func(t *testing.T) {
		aSchemaString := "l:bool:true p:int:80 g:[int]"
		argString := "-l -p 8080 -g 1,3,5,7"
		flagTests := []flagTest{
			{"l", "bool", true},
			{"p", "int", 8080},
			{"g", "[int]", []int{1, 3, 5, 7}},
		}
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		aParser.parse(argString)
		testGetArgValue(t, aParser, flagTests)
	})

	t.Run("test string list arg pair", func(t *testing.T) {
		aSchemaString := "l:bool:true p:int:80 g:[string]"
		argString := "-l -p 8080 -g this,is,a,list"
		flagTests := []flagTest{
			{"l", "bool", true},
			{"p", "int", 8080},
			{"g", "[string]", []string{"this", "is", "a", "list"}},
		}
		aParser := newParser(aSchemaString)
		assertNil(t, aParser)
		aParser.parse(argString)
		testGetArgValue(t, aParser, flagTests)
	})
}

func testGetArgValue(t *testing.T, aParser *Parser, flagTests []flagTest) {
	for _, tt := range flagTests {
		want := tt.value
		switch tt.typeCode {
		case "bool":
			got := aParser.GetBoolArg(tt.flag)
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		case "int":
			got := aParser.GetIntArg(tt.flag)
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		case "string":
			got := aParser.GetStringArg(tt.flag)
			assertStrings(t, got, want.(string))
		case "[int]":
			got := aParser.GetIntListArg(tt.flag)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		case "[string]":
			got := aParser.GetStringListArg(tt.flag)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		default:
			t.Errorf("not support type")
		}
	}
}

func testSchemaRules(t *testing.T, aSchema *Schema, srts []schemaRuleTest) {
	for _, srt := range srts {
		sr, err := aSchema.getSchemaRule(srt.flag)
		if err != nil {
			assertError(t, err, srt.err)
			continue
		}
		assertNoError(t, err)
		assertSchemaRule(t, srt, sr)
	}
}

func assertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Fatalf("got nil but didn't want nil")
	}
}

func assertSchemaRule(t *testing.T, tt schemaRuleTest, sr *SchemaRule) {
	t.Helper()
	gotFlag := sr.getFlag()
	assertStrings(t, gotFlag, tt.flag)

	gotTypeCode := sr.getTypeCode()
	assertStrings(t, gotTypeCode, tt.typeCode)

	gotDefaultValue := sr.getDefaultValue()
	assertStrings(t, gotDefaultValue, tt.defaultValue)
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got error but didn't want one")
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatalf("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
