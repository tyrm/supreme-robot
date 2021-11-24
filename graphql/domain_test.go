package graphql

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/queue"
	"net/http"
	"reflect"
	"testing"
)

func TestAddDomainMutator(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, _, err := testDoLoginAdmin(server)
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

	domain := "testadddomainmutator."
	soa := map[string]interface{}{
		"ttl":     300,
		"mbox":    "hostmaster.test.",
		"refresh": 22,
		"retry":   44,
		"expire":  33,
	}
	newID, newDomain, newRecords, err := testDoAddDomain(server, metadata, domain, soa)
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

	// check for job
	testServerQueue.Lock()
	defer testServerQueue.Unlock()

	if len(testServerQueue.Jobs[queue.QueueDNS]) != 1 {
		t.Errorf("unexpected number of jobs in queue %s, got: %d, want: 1", queue.QueueDNS, len(testServerQueue.Jobs[queue.QueueDNS]))
		return
	}
	if len(testServerQueue.Jobs[queue.QueueDNS][0]) != 2 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(testServerQueue.Jobs[queue.QueueDNS][0]))
		return
	}

	if testServerQueue.Jobs[queue.QueueDNS][0][0].(string) != queue.JobAddDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", testServerQueue.Jobs[queue.QueueDNS][0][0], queue.JobAddDomain)
	}
	if testServerQueue.Jobs[queue.QueueDNS][0][1].(string) != newID {
		t.Errorf("unexpected domain id, got: %s, want: %s", testServerQueue.Jobs[queue.QueueDNS][0][1], newID)
	}

	// reset queue
	testServerQueue.Jobs[queue.QueueDNS] = [][]interface{}{}
}

func TestDeleteDomainMutator(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, _, err := testDoLoginAdmin(server)
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
	domain := "testdeletedomainmutator."
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
	t.Logf("testDoDeleteDomain: %v, %v", newSuccess, err)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if newSuccess != true {
		t.Errorf("delete unsuccessful, got: %v, want: true.", newSuccess)
	}

	// check for job
	testServerQueue.Lock()
	defer testServerQueue.Unlock()

	if len(testServerQueue.Jobs[queue.QueueDNS]) != 2 {
		t.Errorf("unexpected number of jobs in queue %s, got: %d, want: 2", queue.QueueDNS, len(testServerQueue.Jobs[queue.QueueDNS]))
		return
	}
	if len(testServerQueue.Jobs[queue.QueueDNS][1]) != 2 {
		t.Errorf("unexpected number of parameters for job %s, got: %d, want: 2", queue.JobAddDomain, len(testServerQueue.Jobs[queue.QueueDNS][0]))
		return
	}

	if testServerQueue.Jobs[queue.QueueDNS][1][0].(string) != queue.JobRemoveDomain {
		t.Errorf("unexpected job type, got: %s, want: %s", testServerQueue.Jobs[queue.QueueDNS][0][0], queue.JobRemoveDomain)
	}
	if testServerQueue.Jobs[queue.QueueDNS][1][1].(string) != newID {
		t.Errorf("unexpected domain id, got: %s, want: %s", testServerQueue.Jobs[queue.QueueDNS][0][1], newID)
	}

	// reset queue
	testServerQueue.Jobs[queue.QueueDNS] = [][]interface{}{}
}

func TestDomainQuery(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, _, err := testDoLoginAdmin(server)
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

	domain := "testdomainquery."
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

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `query (
			$id: String!
		){
			domain(
				id: $id
			){
				id
				domain
				createdAt
				updatedAt
			}
		}`,
		Variables: map[string]interface{}{
			"id": newID,
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
		t.Fatalf("unexpected error, got: %s, want: nil.", result.Errors[0].Error())
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("can't cast data, got: %d, want: map[string]interface{}", reflect.TypeOf(result.Data))
	}

	receivedDomain, receivedDomainOK := data["domain"].(map[string]interface{})
	if !receivedDomainOK {
		t.Errorf("can't cast domain, got: %d, want: map[string]interface{}", reflect.TypeOf(data["domain"]))
	}

	receivedID, receivedIDOK := receivedDomain["id"].(string)
	if !receivedIDOK {
		t.Errorf("can't cast domain id, got: %d, want: string", reflect.TypeOf(data["id"]))
	}
	if receivedID != newID {
		t.Errorf("unexpected domain id, got: %s, want: %s", receivedID, newID)
	}

	receivedName, receivedNameOK := receivedDomain["domain"].(string)
	if !receivedNameOK {
		t.Errorf("can't cast domain id, got: %d, want: string", reflect.TypeOf(data["domain"]))
	}
	if receivedName != domain {
		t.Errorf("unexpected domain id, got: %s, want: %s", receivedName, domain)
	}

}

func TestMyDomainsQuery(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// create new user
	newUser := models.User{
		Username: "newuser",
	}
	err = newUser.SetPassword("newpassword")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	err = testServerDB.Create(&newUser)

	domains := []string{
		"testmydomainsquery1.",
		"testmydomainsquery2.",
		"testmydomainsquery3.",
	}

	// create domains
	for _, d := range domains {
		newDomain := models.Domain{
			Domain:  d,
			OwnerID: newUser.ID,
		}
		dbErr := testServerDB.Create(&newDomain)
		if dbErr != nil {
			t.Errorf("unexpected error, got: %s, want: nil.", dbErr.Error())
			return
		}
	}

	// do login
	accessToken, _, err := testDoLogin(server, newUser.Username, "newpassword")
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

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `{
			myDomains{
				id
				domain
				createdAt
				updatedAt
			}
		}`,
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
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Fatalf("no data returned")
	}

	myDomains, myDomainsOk := data["myDomains"].([]interface{})
	if !myDomainsOk {
		t.Fatalf("can't cast result, got: %d, want: []interface{}", reflect.TypeOf(data["myDomains"]))

	}
	if len(myDomains) != 3 {
		t.Errorf("invalid number of rows returned, got: %d, want: 3.", len(myDomains))
	}

	searchDomains := make(map[string]bool)
	for _, d := range domains {
		searchDomains[d] = false
	}
	for _, myDomain := range myDomains {
		d, dOK := myDomain.(map[string]interface{})
		if !dOK {
			t.Errorf("can't cast result, got: %d, want: map[string]interface{}", reflect.TypeOf(myDomain))
		} else {
			dName, dNameOK := d["domain"].(string)
			if !dNameOK {
				t.Errorf("can't cast result, got: %d, want: map[string]interface{}", reflect.TypeOf(myDomain))
			}
			for _, sd := range domains {
				if dName == sd {
					searchDomains[sd] = true
				}
			}
		}
	}

	for k, v := range searchDomains {
		if !v {
			t.Errorf("didn't find expected domain: %s", k)
		}
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
