package graphql

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/config"
	dbMem "github.com/tyrm/supreme-robot/db/memory"
	kvMem "github.com/tyrm/supreme-robot/kv/memory"
	queueMem "github.com/tyrm/supreme-robot/queue/memory"
	"testing"
	"time"
)

func TestLoginMutator(t *testing.T) {
	cnf := config.Config{
		AccessExpiration:  time.Hour * 24,
		AccessSecret:      "test",
		PrimaryNS:         "ns1.example.com.",
		RefreshExpiration: time.Hour * 24,
		RefreshSecret:     "test1234",
	}

	db, err := dbMem.NewClient()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	kv, err := kvMem.NewClient()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	qc, err := queueMem.NewScheduler()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	// create web server
	server, err := NewServer(&cnf, qc, db, kv)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// prepare query
	ctx := context.Background()
	p := postData{
		Query: `mutation (
    $username: String!
    $password: String!
){
    login(
        username: $username
        password: $password
    ){
        accessToken
        refreshToken
    }
}`,
		Variables: map[string]interface{}{
			"username": "admin",
			"password": "password",
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

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("no data returned.")
		return
	}

	login, loginOK := data["login"].(map[string]interface{})
	if !loginOK {
		t.Errorf("no login data returned.")
		return
	}

	_, accessTokenOK := login["accessToken"].(string)
	if !accessTokenOK {
		t.Errorf("no accessToken returned.")
		return
	}

	_, refreshTokenOK := login["refreshToken"].(string)
	if !refreshTokenOK {
		t.Errorf("no refreshToken returned.")
		return
	}
}
