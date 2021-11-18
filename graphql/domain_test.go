package graphql

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
	"testing"
)

func TestAddDomainMutator(t *testing.T) {
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
	_, newDomain, newRecords, err := testDoAddDomain(server, metadata, domain, soa)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if newDomain != domain {
		t.Errorf("invalid domain in reponse, got: %s, want: %s.", newDomain, domain)
	}

	for _, r := range newRecords {
		record, recordOk := r.(map[string]interface{})
		if !recordOk {
			t.Errorf("record %#v is not map[string]interface{}", record)
			continue
		}

		fmt.Printf("%#v\n", r)

		switch record["type"].(string) {
		case models.RecordTypeSOA:
			ttl, ttlOk := record["ttl"].(int)
			if !ttlOk {
				t.Errorf("soa record missing ttl")
			} else {
				if ttl != 300 {
					t.Errorf("soa record had invalid ttl, got: %d, want: 300.", ttl)
				}
			}
			value, valueOk := record["value"].(string)
			if !valueOk {
				t.Errorf("soa record missing value")
			} else {
				if value != "ns1.example.com." {
					t.Errorf("soa record had invalid value, got: '%s', want: 'ns1.example.com.'.", value)
				}
			}
			mbox, mboxOk := record["mbox"].(string)
			if !mboxOk {
				t.Errorf("soa record missing mbox")
			} else {
				if mbox != "hostmaster.test." {
					t.Errorf("soa record had invalid mbox, got: '%s', want: 'hostmaster.test.'.", mbox)
				}
			}
			refresh, refreshOk := record["refresh"].(int)
			if !refreshOk {
				t.Errorf("soa record missing refresh")
			} else {
				if refresh != 22 {
					t.Errorf("soa record had invalid refresh, got: %d, want: 22.", refresh)
				}
			}
			retry, retryOk := record["retry"].(int)
			if !retryOk {
				t.Errorf("soa record missing retry")
			} else {
				if retry != 44 {
					t.Errorf("soa record had invalid retry, got: %d, want: 44.", retry)
				}
			}
			expire, expireOk := record["expire"].(int)
			if !expireOk {
				t.Errorf("soa record missing expire")
			} else {
				if expire != 33 {
					t.Errorf("soa record had invalid expire, got: %d, want: 44.", expire)
				}
			}
		}
	}
}

func TestDeleteDomainMutator(t *testing.T) {
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

	// add domain
	domain := "test."
	soa := map[string]interface{}{
		"ttl":     300,
		"mbox":    "hostmaster.test.",
		"refresh": 22,
		"retry":   44,
		"expire":  33,
	}
	newID, _, _, err := testDoAddDomain(server, metadata, domain, soa)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// delete domain
	newSuccess, err := testDoDeleteDomain(server, metadata, newID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if newSuccess != true {
		t.Errorf("delete unsuccessful, got: %v, want: true.", newSuccess)
	}
}

func testDoAddDomain(server *Server, metadata *accessDetails, d string, soa map[string]interface{}) (string, string, []interface{}, error) {
	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `mutation (
			$domain: String!
			$soa: SOA!
		){
			addDomain(
				domain: $domain
				soa: $soa
			){
				id
				domain
				createdAt
				updatedAt
				records{
					id
					name
					type
					value
					ttl
					mbox
					refresh
					retry
					expire
					createdAt
					updatedAt
				}
			}
		}`,
		Variables: map[string]interface{}{
			"domain": d,
			"soa":    soa,
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
		return "", "", nil, result.Errors[0]
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		return "", "", nil, fmt.Errorf("no data returned")
	}

	addDomain, addDomainOK := data["addDomain"].(map[string]interface{})
	if !addDomainOK {
		return "", "", nil, fmt.Errorf("no addDomain data returned")
	}

	id, idOK := addDomain["id"].(string)
	if !idOK {
		return "", "", nil, fmt.Errorf("no id returned")
	}

	domain, domainOK := addDomain["domain"].(string)
	if !domainOK {
		return "", "", nil, fmt.Errorf("no domain returned")
	}

	records, recordsOK := addDomain["records"].([]interface{})
	if !recordsOK {
		return "", "", nil, fmt.Errorf("no records returned")
	}

	return id, domain, records, nil
}

func testDoDeleteDomain(server *Server, metadata *accessDetails, id string) (bool, error) {
	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `mutation (
			$id: String!
		){
			deleteDomain(
				id: $id
			){
				success
			}
		}`,
		Variables: map[string]interface{}{
			"id": id,
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
		return false, result.Errors[0]
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		return false, fmt.Errorf("no data returned")
	}

	deleteDomain, deleteDomainOK := data["deleteDomain"].(map[string]interface{})
	if !deleteDomainOK {
		return false, fmt.Errorf("no deleteDomain data returned")
	}

	newSuccess, newSuccessOK := deleteDomain["success"].(bool)
	if !newSuccessOK {
		return false, fmt.Errorf("no success data returned")
	}

	return newSuccess, nil
}
