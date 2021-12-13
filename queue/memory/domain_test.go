package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/queue"
	"reflect"
	"testing"
	"time"
)

func TestScheduler_AddDomain(t *testing.T) {
	scheduler, err := NewMemQueue()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&MemQueue{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&MemQueue{}))
	}

	domainID := uuid.MustParse("5c664d04-5d00-4a7b-be37-fa62538985a1")

	err = scheduler.AddDomain(domainID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// wait for job
	var incomingJob []interface{}
	select {
	case incomingJob = <-scheduler.Queues[queue.QueueDNS]:
	case <-time.After(60 * time.Second):
		t.Errorf("job not queued")
		return
	}

	if len(incomingJob) != 3 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(incomingJob))
		return
	}

	if incomingJob[0].(string) != queue.JobAddDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", incomingJob[0], queue.JobAddDomain)
	}
	if len(incomingJob[1].(string)) != 16 {
		t.Errorf("unexpected jid, got: %d, want: 24", len(incomingJob[1].(string)))
	}
	if incomingJob[2].(string) != domainID.String() {
		t.Errorf("unexpected domain id, got: %s, want: %s", incomingJob[1], domainID.String())
	}
}

func TestScheduler_RemoveDomain(t *testing.T) {
	scheduler, err := NewMemQueue()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&MemQueue{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&MemQueue{}))
	}

	domainID := uuid.MustParse("5c664d04-5d00-4a7b-be37-fa62538985a1")

	err = scheduler.RemoveDomain(domainID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// wait for job
	var incomingJob []interface{}
	select {
	case incomingJob = <-scheduler.Queues[queue.QueueDNS]:
	case <-time.After(60 * time.Second):
		t.Errorf("job not queued")
		return
	}

	if len(incomingJob) != 3 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(incomingJob))
		return
	}

	if incomingJob[0].(string) != queue.JobRemoveDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", incomingJob[0], queue.JobAddDomain)
	}
	if len(incomingJob[1].(string)) != 16 {
		t.Errorf("unexpected jid, got: %d, want: 24", len(incomingJob[1].(string)))
	}
	if incomingJob[2].(string) != domainID.String() {
		t.Errorf("unexpected domain id, got: %s, want: %s", incomingJob[1], domainID.String())
	}
}

func TestScheduler_UpdateSubDomain(t *testing.T) {
	scheduler, err := NewMemQueue()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
	}
	if reflect.TypeOf(scheduler) != reflect.TypeOf(&MemQueue{}) {
		t.Errorf("unexpected scheduler type, got: %s, want: %s", reflect.TypeOf(scheduler), reflect.TypeOf(&MemQueue{}))
	}

	domainID := uuid.MustParse("5c664d04-5d00-4a7b-be37-fa62538985a1")
	name := "@"

	err = scheduler.UpdateSubDomain(domainID, name)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// wait for job
	var incomingJob []interface{}
	select {
	case incomingJob = <-scheduler.Queues[queue.QueueDNS]:
	case <-time.After(60 * time.Second):
		t.Errorf("job not queued")
		return
	}

	if len(incomingJob) != 4 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(incomingJob))
		return
	}

	if incomingJob[0].(string) != queue.JobUpdateSubDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", incomingJob[0], queue.JobAddDomain)
	}
	if len(incomingJob[1].(string)) != 16 {
		t.Errorf("unexpected jid, got: %d, want: 24", len(incomingJob[1].(string)))
	}
	if incomingJob[2].(string) != domainID.String() {
		t.Errorf("unexpected domain id, got: %s, want: %s", incomingJob[1], domainID.String())
	}
	if incomingJob[3].(string) != name {
		t.Errorf("unexpected name, got: %s, want: %s", incomingJob[1], domainID.String())
	}
}
