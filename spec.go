package gspec

import (
	"reflect"
	"strings"
	"testing"
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

func New(t *testing.T) *S {
	return &S{t: t}
}

func (s *S) Expect(actual interface{}, garbage ...interface{}) (sr *SR) {
	return &SR{s.t, actual}
}

func (s *S) ExpectBytes(actual []byte, garbage ...interface{}) (sr *SRB) {
	return &SRB{s.t, actual}
}

func (sr *SR) ToEqual(expected interface{}) {
	if sr.actual != expected {
		sr.t.Errorf("expected %+v to equal %+v", expected, sr.actual)
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
				sr.t.Errorf("Byte %d mismatch, expected %d got %d", index, x[index], sr.actual[b])
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
			t.Errorf("Expected %q to not contain %q", a, expected)
		}
	default:
		t.Errorf("trying to call [Not]Contains on an unsuported type")
	}
}
