package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/queue"
	"reflect"
	"testing"
)

func TestScheduler_AddDomain(t *testing.T) {
	scheduler, err := NewScheduler()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&Scheduler{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&Scheduler{}))
	}

	domainID := uuid.MustParse("5c664d04-5d00-4a7b-be37-fa62538985a1")

	err = scheduler.AddDomain(domainID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if len(scheduler.Jobs[queue.QueueDNS]) != 1 {
		t.Errorf("unexpected number of jobs in queue %s, got: %d, want: 1", queue.QueueDNS, len(scheduler.Jobs[queue.QueueDNS]))
		return
	}
	if len(scheduler.Jobs[queue.QueueDNS][0]) != 2 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(scheduler.Jobs[queue.QueueDNS][0]))
		return
	}

	if scheduler.Jobs[queue.QueueDNS][0][0].(string) != queue.JobAddDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", scheduler.Jobs[queue.QueueDNS][0][0], queue.JobAddDomain)
	}
	if scheduler.Jobs[queue.QueueDNS][0][1].(string) != domainID.String() {
		t.Errorf("unexpected domain id, got: %s, want: %s", scheduler.Jobs[queue.QueueDNS][0][1], domainID.String())
	}
}

func TestScheduler_RemoveDomain(t *testing.T) {
	scheduler, err := NewScheduler()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&Scheduler{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&Scheduler{}))
	}

	domainID := uuid.MustParse("5c664d04-5d00-4a7b-be37-fa62538985a1")

	err = scheduler.RemoveDomain(domainID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if len(scheduler.Jobs[queue.QueueDNS]) != 1 {
		t.Errorf("unexpected number of jobs in queue %s, got: %d, want: 1", queue.QueueDNS, len(scheduler.Jobs[queue.QueueDNS]))
		return
	}
	if len(scheduler.Jobs[queue.QueueDNS][0]) != 2 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(scheduler.Jobs[queue.QueueDNS][0]))
		return
	}

	if scheduler.Jobs[queue.QueueDNS][0][0].(string) != queue.JobRemoveDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", scheduler.Jobs[queue.QueueDNS][0][0], queue.JobRemoveDomain)
		return
	}
	if scheduler.Jobs[queue.QueueDNS][0][1].(string) != domainID.String() {
		t.Errorf("unexpected domain id, got: %s, want: %s", scheduler.Jobs[queue.QueueDNS][0][1], domainID.String())
		return
	}

}
