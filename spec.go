package gspec

import (
	"reflect"
	"strings"
	"testing"
	"regexp"
	"encoding/json"
)

var (
	jsonFlattenPattern = regexp.MustCompile("\\s")
)

type S struct {
	t *testing.T
}

type SR struct {
	t      *testing.T
	actual interface{}
}

type SRB struct {
	t      *testing.T
	actual []byte
}

type SRJ struct {
	t      *testing.T
	actual interface{}
}

func New(t *testing.T) *S {
	return &S{t: t}
}

func (s *S) Expect(actual interface{}, garbage ...interface{}) (sr *SR) {
	return &SR{s.t, actual}
}

func (s *S) ExpectBytes(actual []byte, garbage ...interface{}) (sr *SRB) {
	return &SRB{s.t, actual}
}

func (s *S) ExpectJson(actual interface{}, garbage ...interface{}) (sr *SRJ) {
	return &SRJ{s.t, actual}
}

func (sr *SR) ToEqual(expecteds ...interface{}) {
	kind := reflect.TypeOf(sr.actual).Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		if len(expecteds) == 1 {
			sr.compareArrays(expecteds[0])
		} else {
			sr.compareArrays(expecteds)
		}
		return
	}
	if sr.actual != expecteds[0] {
		sr.t.Errorf("expected %+v to equal %+v", expecteds[0], sr.actual)
	}
}

func (sr *SR) compareArrays(e interface{}) {
	actual := reflect.ValueOf(sr.actual)
	expected := reflect.ValueOf(e)
	length := actual.Len()
	if length != expected.Len() {
		sr.t.Errorf("expected an array of %d values, got %d", expected.Len(), length)
		return
	}
	for i := 0; i < length; i++ {
		if actual.Index(i).Interface() != expected.Index(i).Interface() {
			sr.t.Errorf("item at index %d should have been %q, got %q", i, expected.Index(i).Interface(), actual.Index(i).Interface())
		}
	}
}

func (sr *SR) ToContain(expected interface{}) {
	sr.contains(expected, true)
}

func (sr *SR) ToNotContain(expected interface{}) {
	sr.contains(expected, false)
}

func (sr *SR) contains(expected interface{}, b bool) {
	contains(sr.t, sr.actual, expected, b)
}

func (sr *SR) ToNotEqual(expected interface{}) {
	if sr.actual == expected {
		sr.t.Errorf("expected %+v to not equal %+v", expected, sr.actual)
	}
}

func (sr *SR) ToBeNil() {
	if sr.actual == nil {
		return
	}
	if reflect.ValueOf(sr.actual).IsNil() {
		return
	}
	sr.t.Errorf("expected %+v to be nil", sr.actual)
}

func (sr *SR) ToNotBeNil() {
	if !reflect.ValueOf(sr.actual).IsNil() {
		return
	}
	sr.t.Errorf("expected %+v to not be nil", sr.actual)
}

func (sr *SRB) ToEqual(expected interface{}) {
	switch x := expected.(type) {
	case string:
		if x != string(sr.actual) {
			sr.t.Errorf("Expected %q, got %q", x, string(sr.actual))
		}
	case []byte:
		if len(sr.actual) != len(x) {
			sr.t.Errorf("expected %d byte values, got %d", len(x), len(sr.actual))
		}
		for index, b := range sr.actual {
			if b != x[index] {
				sr.t.Errorf("Byte %d mismatch, expected %d got %d", index, x[index], b)
			}
		}
	}
}

func (srb *SRB) ToContain(expected interface{}) {
	srb.contains(expected, true)
}

func (srb *SRB) ToNotContain(expected interface{}) {
	srb.contains(expected, false)
}

func (srb *SRB) contains(expected interface{}, b bool) {
	contains(srb.t, string(srb.actual), expected, b)
}

func contains(t *testing.T, actual interface{}, expected interface{}, b bool) {
	switch a := actual.(type) {
	case string:
		if strings.Contains(a, expected.(string)) != b {
			if b {
				t.Errorf("Expected %q to contain %q", a, expected)
			} else {
				t.Errorf("Expected %q to not contain %q", a, expected)
			}

		}
	default:
		t.Errorf("trying to call [Not]Contains on an unsuported type")
	}
}

func (sr *SRJ) ToEqual(expected string) {
	b, ok := sr.actual.([]byte)
	if ok == false {
		if sb, ok := sr.actual.(string); ok == false {
			sr.t.Errorf("JSON must be a []byte or string")
			return
		} else {
			b = []byte(sb)
		}

	}
	actualMap := make(map[string]interface{})
	expectedMap := make(map[string]interface{})
	if err := json.Unmarshal(b, &actualMap); err != nil {
		sr.t.Error(err)
		return
	}
	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		sr.t.Error(err)
		return
	}
	if reflect.DeepEqual(actualMap, expectedMap) == false {
		sr.t.Errorf("Expected:\n%s\nto equal\n%s", jsonFlattenPattern.ReplaceAllLiteralString(expected, ""), jsonFlattenPattern.ReplaceAllLiteralString(string(b), ""))
	}
}
