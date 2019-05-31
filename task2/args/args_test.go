package args

import (
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
	}{
		{
			"l:bool:true p:int:80 d:string:/usr/logs",
			[]schemaRuleTest{
				{"l:bool:true", "l", "bool", "true", nil},
				{"p:int:80", "p", "int", "80", nil},
				{"d:string:/usr/logs", "d", "string", "/usr/logs", nil},
			},
		},
		{
			"l:bool:true p:int:80",
			[]schemaRuleTest{
				{"l:bool:true", "l", "bool", "true", nil},
				{"p:int:80", "p", "int", "80", nil},
			},
		},
	}

	for _, tt := range schemaTests {
		aSchema := newSchema(tt.schemaString)
		if aSchema == nil {
			t.Errorf("got nil but didn't want nil")
		}

		wantSchemaRuleCount := len(tt.schemaRuleTests)
		got := aSchema.count()
		if got != wantSchemaRuleCount {
			t.Errorf("got %d, want %d", got, wantSchemaRuleCount)
		}

		for _, srt := range tt.schemaRuleTests {
			sr := aSchema.getSchemaRule(srt.flag)
			if sr == nil {
				t.Errorf("got nil but didn't want nil")
			}
			assertSchemaRule(t, srt, sr)
		}
	}
}

func assertSchemaRule(t *testing.T, tt schemaRuleTest, sr *SchemaRule) {
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
