package memory

import (
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/tyrm/supreme-robot/queue"
	"sync"
)

var queues = []string{
	queue.QueueDNS,
}

// Terminate closes all channels
func (m *MemQueue) Terminate(_ bool) {
	for _, q := range queues {
		if _, ok := <-m.Queues[q]; ok {
			close(m.Queues[q])
		}
	}
}

// ProcessStrictPriorityQueues enables queues. should be called once before Run()
func (m *MemQueue) ProcessStrictPriorityQueues(queues ...string) {
	for _, q := range queues {
		m.QueuesEnabled[q] = true
	}
}

// Register a handler for the given jobtype.  It is expected that all jobtypes
// are registered upon process startup.
func (m *MemQueue) Register(job string, f faktory.Perform) {
	m.Callbacks[job] = f
}

// Run starts processing jobs.
func (m *MemQueue) Run() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	for _, q := range queues {
		if m.QueuesEnabled[q] {
			go func(m *MemQueue, wg *sync.WaitGroup) {
				for item := range m.Queues[q] {
					jobType, jobTypeOK := item[0].(string)
					jid, jidOK := item[0].(string)
					if jobTypeOK && jidOK {
						switch jobType {
						case queue.JobAddDomain:
							if m.Callbacks[queue.JobAddDomain] != nil {
								err := m.Callbacks[queue.JobAddDomain](makePseudoContext(jid), item[1].(string))
								if err != nil {
									logger.Errorf("[%s] callback error: %s", queue.JobAddDomain, err.Error())
									m.Queues[q] <- item
								}
							}
						}
					} else {
						logger.Errorf("[%s] can't cast job", queue.JobAddDomain)
						m.Queues[q] <- item
					}
				}
				wg.Done()
			}(m, &waitgroup)
		}
	}
	waitgroup.Wait()
	return
}
