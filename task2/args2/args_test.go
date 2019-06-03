package args2

import (
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
		assertNil(t, sr)
		assertNoError(t, err)

		gotFlag := sr.getFlag()
		assertStrings(t, gotFlag, srt.flag)

		gotTypeCode := sr.getTypeCode()
		assertStrings(t, gotTypeCode, srt.typeCode)

		gotDefaultValue := sr.getDefaultValue()
		assertStrings(t, gotDefaultValue, srt.defaultValue)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got an error but didn't want one")
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

func assertNil(t *testing.T, got interface{}) {
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
