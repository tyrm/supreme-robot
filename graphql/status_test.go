package graphql

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/version"
	"testing"
)

func TestStatusQuery(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// prepare query
	p := postData{
		Query: `{
			status{
				version
			}
		}`,
	}

	// do query
	result := graphql.Do(graphql.Params{
		Context:        context.Background(),
		Schema:         server.schema(),
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if result.HasErrors() {
		t.Errorf("unexpected error, got: %s, want: nil.", result.Errors[0].Error())
		return
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("no data returned")
		return
	}

	sts, statusOk := data["status"].(map[string]interface{})
	if !statusOk {
		t.Errorf("no status data returned")
		return
	}

	ver, versionOk := sts["version"].(string)
	if !versionOk {
		t.Errorf("no version data returned")
	}

	if ver != version.Version {
		t.Errorf("got invalid version, got: %s, want: %s.", ver, version.Version)
	}
}
