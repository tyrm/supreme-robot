package memory

import (
	"context"
	faktoryClient "github.com/contribsys/faktory/client"
)

// internal keys for context value storage
type valueKey int

const (
	poolKey valueKey = 2
	jobKey  valueKey = 3
)

func makePseudoContext(jid string) context.Context {
	fakeJob := faktoryClient.Job{
		Jid: jid,
	}
	return context.WithValue(context.Background(), jobKey, &fakeJob)
}
