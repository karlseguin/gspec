package gspec

import (
	"reflect"
	"testing"
)

type S struct {
	t *testing.T
}

type SR struct {
	t        *testing.T
	actual interface{}
}

type SRB struct {
	t        *testing.T
	actual []byte
}

func New(t *testing.T) *S {
	return &S{t: t}
}

func (s *S) Expect(actual interface{}, garbage ...interface{}) (sr *SR) {
	return &SR{s.t, actual}
}

func (s *S) ExpectBytes(actual []byte, garbage ...interface{}) (sr *SRB) {
	return &SRB{s.t,actual}
}

func (sr *SR) ToEqual(expected interface{}) {
	if sr.actual != expected {
		sr.t.Errorf("expected %+v to equal %+v", expected, sr.actual)
	}
}

func (sr *SRB) ToEqual(expected ...byte) {
	if len(sr.actual) != len(expected) {
		sr.t.Errorf("expected %d byte values, got %d", len(expected), len(sr.actual))
	}
	for index, b := range sr.actual {
		if b != expected[index] {
			sr.t.Errorf("Byte %d mismatch, expected %d got %d", index, expected[index], sr.actual[b])
		}
	}
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
