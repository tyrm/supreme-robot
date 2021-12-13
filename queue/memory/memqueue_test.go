package memory

import (
	"reflect"
	"testing"
)

func TestNewMemQueue(t *testing.T) {
	memqueue, err := NewMemQueue()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(memqueue) != reflect.TypeOf(&MemQueue{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(memqueue), reflect.TypeOf(&MemQueue{}))
	}
}
