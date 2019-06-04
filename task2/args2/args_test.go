package args2

import (
	"reflect"
	"testing"
)

func TestSchemaRule(t *testing.T) {
	schemaRuleTests := []struct {
		schemaRule   string
		flag         string
		typeCode     string
		defaultValue string
		err          error
	}{
		{"l:bool:false", "l", "bool", "false", nil},
		{"l:bool:", "l", "bool", "false", nil},
		{"l:bool", "l", "bool", "false", nil},
		{"l", "l", "bool", "false", ErrWrongSchemaRule},
		{"p:int", "p", "int", "0", nil},
	}

	for _, srt := range schemaRuleTests {
		sr, err := newSchemaRule(srt.schemaRule)
		if srt.err != nil {
			assertError(t, err, srt.err)
			continue
		}
		assertNotNil(t, sr)
		assertNoError(t, err)

		gotFlag := sr.getFlag()
		assertStrings(t, gotFlag, srt.flag)

		gotTypeCode := sr.getTypeCode()
		assertStrings(t, gotTypeCode, srt.typeCode)

		gotDefaultValue := sr.getDefaultValue()
		assertStrings(t, gotDefaultValue, srt.defaultValue)
	}
}

func TestSchema(t *testing.T) {
	aSchemaString := "l:bool p:int:80 d:string"
	aSchema, err := newSchema(aSchemaString)
	assertNoError(t, err)
	assertNotNil(t, aSchema)

	wantSize := 3
	wantSchemaRules := []struct {
		flag         string
		typeCode     string
		defaultValue string
		err          error
	}{
		{"l", "bool", "false", nil},
		{"p", "int", "80", nil},
		{"d", "string", "", nil},
		{"e", "int", "0", ErrorFlagNotExist},
	}

	gotSize := aSchema.size()
	if gotSize != wantSize {
		t.Errorf("got %d, want %d", gotSize, wantSize)
	}

	for _, tt := range wantSchemaRules {
		typeCode, err := aSchema.typeOf(tt.flag)
		if tt.err != nil {
			assertError(t, err, tt.err)
			continue
		}
		assertStrings(t, typeCode, tt.typeCode)

		defaultValue, err := aSchema.defaultValueOf(tt.flag)
		if tt.err != nil {
			assertError(t, err, tt.err)
			continue
		}
		assertStrings(t, defaultValue, tt.defaultValue)
	}
}

type argument struct {
	flag     string
	typeCode string
	value    interface{}
	err      error
}

func TestParser(t *testing.T) {
	aSchemaString := "l:bool p:int:80 d:string"
	aParser, err := newParser(aSchemaString)
	assertNoError(t, err)
	assertNotNil(t, aParser)

	wantArguments := []argument{
		{"d", "string", "/usr/logs", nil},
		{"l", "bool", true, nil},
		{"p", "int", 8080, nil},
		{"e", "string", "not_exist", ErrorFlagNotExist},
	}
	argumentTests := []struct {
		name           string
		argumentString string
	}{
		{"full argument parse", "-l true -p 8080 -d /usr/logs"},
		{"simple argument parse", "-l -p 8080 -d /usr/logs"},
	}

	for _, tt := range argumentTests {
		testParse(t, tt.name, tt.argumentString, aParser, wantArguments)
	}
}

func testParse(t *testing.T, name, argumentsString string, aParser *Parser, wantArguments []argument) {
	t.Run(name, func(t *testing.T) {
		err := aParser.parse(argumentsString)
		assertNoError(t, err)

		var v interface{}
		for _, tt := range wantArguments {
			switch tt.typeCode {
			case "string":
				v, err = aParser.stringValueOf(tt.flag)
			case "bool":
				v, err = aParser.boolValueOf(tt.flag)
			case "int":
				v, err = aParser.intValueOf(tt.flag)
			default:
				t.Errorf("not support type")
			}
			if tt.err != nil {
				assertError(t, err, tt.err)
			} else {
				assertNoError(t, err)
				assertEqual(t, v, tt.value)
			}
		}
	})
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one")
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatalf("did not get an error but want one")
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("got nil but didn't want nil")
	}
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
