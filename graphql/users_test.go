package graphql

import (
	"context"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"net/http"
	"testing"
)

func TestMeQuery(t *testing.T) {
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

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `{
			me{
				id
				username
				groups
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
		t.Errorf("unexpected error, got: %s, want: nil.", result.Errors[0].Error())
		return
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("no data returned")
		return
	}

	me, meOk := data["me"].(map[string]interface{})
	if !meOk {
		t.Errorf("no me data returned")
		return
	}

	meID, meIDOk := me["id"].(string)
	if !meIDOk {
		t.Errorf("no me id data returned")
	}
	meIDUUID, err := uuid.Parse(meID)
	if err != nil {
		t.Errorf("can't parse me id %s, got: %s, want: nil.", meID, err.Error())
	}
	if meIDUUID != uuid.Must(uuid.Parse("44892097-2c97-4c16-b4d1-e8522586df48")) {
		t.Errorf("got wrong id, got: %s, want: '44892097-2c97-4c16-b4d1-e8522586df48'.", meID)
	}

	username, usernameOk := me["username"].(string)
	if !usernameOk {
		t.Errorf("no me username data returned")
	}
	if username != "admin" {
		t.Errorf("got wrong username, got: %s, want: 'admin'.", meID)
	}

	createdAt, createdAtOk := me["createdAt"].(int)
	if !createdAtOk {
		t.Errorf("no me createdAt data returned")
	}
	if createdAt <= 0 {
		t.Errorf("got invalid createdAt, got: %d, want: >0.", createdAt)
	}

	updatedAt, updatedAtOk := me["updatedAt"].(int)
	if !updatedAtOk {
		t.Errorf("no me updatedAt data returned")
	}
	if updatedAt <= 0 {
		t.Errorf("got invalid updatedAt, got: %d, want: >0.", updatedAt)
	}
}
