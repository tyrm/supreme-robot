package graphql

import (
	"context"
	"github.com/graphql-go/graphql"
	"net/http"
	"testing"
)

func TestAddRecordAMutator(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, _, err := testDoLogin(server, "admin", "password")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// extract metadata
	req := http.Request{}
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	metadata, err := server.extractTokenMetadata(&req)
	if err != nil {
		t.Errorf("unexpected error, got: %#v, want: nil.", err.Error())
	}

	domain := "test."
	soa := map[string]interface{}{
		"ttl":     300,
		"mbox":    "hostmaster.test.",
		"refresh": 22,
		"retry":   44,
		"expire":  33,
	}
	domainID, _, _, err := testDoAddDomain(server, metadata, domain, soa)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `mutation (
			$domainId: String!
			$name: String!
			$ip: String!
			$ttl: Int!
		){
			addRecordA(
				domainId: $domainId
				name: $name
				ip: $ip
				ttl: $ttl
			){
				id
				name
				value
				ttl
			}
		}`,
		Variables: map[string]interface{}{
			"domainId": domainID,
			"name":     "@",
			"ip":       "10.1.64.1",
			"ttl":      300,
		},
	}

	// do query
	result := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schema(),
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if result.HasErrors() {
		for _, e := range result.Errors {
			t.Errorf("unexpected error, got: %#v, want: nil.", e.Error())
		}
		return
	}
}

func TestAddRecordAAAAMutator(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, _, err := testDoLogin(server, "admin", "password")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// extract metadata
	req := http.Request{}
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	metadata, err := server.extractTokenMetadata(&req)
	if err != nil {
		t.Errorf("unexpected error, got: %#v, want: nil.", err.Error())
	}

	domain := "test."
	soa := map[string]interface{}{
		"ttl":     300,
		"mbox":    "hostmaster.test.",
		"refresh": 22,
		"retry":   44,
		"expire":  33,
	}
	domainID, _, _, err := testDoAddDomain(server, metadata, domain, soa)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `mutation (
			$domainId: String!
			$name: String!
			$ip: String!
			$ttl: Int!
		){
			addRecordAAAA(
				domainId: $domainId
				name: $name
				ip: $ip
				ttl: $ttl
			){
				id
				name
				value
				ttl
			}
		}`,
		Variables: map[string]interface{}{
			"domainId": domainID,
			"name":     "@",
			"ip":       "2001:db8::1",
			"ttl":      300,
		},
	}

	// do query
	result := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schema(),
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if result.HasErrors() {
		for _, e := range result.Errors {
			t.Errorf("unexpected error, got: %#v, want: nil.", e.Error())
		}
		return
	}
}
