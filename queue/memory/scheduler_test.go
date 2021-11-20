package memory

import (
	"reflect"
	"testing"
)

func TestNewScheduler(t *testing.T) {
	scheduler, err := NewScheduler()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&Scheduler{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&Scheduler{}))
	}
}
