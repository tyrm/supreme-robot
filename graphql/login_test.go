package graphql

import (
	"context"
	"github.com/graphql-go/graphql"
	"net/http"
	"testing"
)

func TestLoginMutator_BadPassword(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
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
			"password": "badpassword",
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

	if !result.HasErrors() {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	// validate error
	err = result.Errors[0]
	if err.Error() != errBadLogin.Error() {
		t.Errorf("unexpected error, got: '%s', want: '%s'.", err.Error(), errBadLogin.Error())

	}
}

func TestLoginMutator_BadUsername(t *testing.T) {
	// create web server
	server, _, _, _, err := newTestServer()
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
			"username": "notauser",
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

	if !result.HasErrors() {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	// validate error
	err = result.Errors[0]
	if err.Error() != errBadLogin.Error() {
		t.Errorf("unexpected error, got: '%s', want: '%s'.", err.Error(), errBadLogin.Error())
	}
}

func TestLoginMutator_ValidLogin(t *testing.T) {
	// create web server
	server, _, _, _, err := newTestServer()
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

func TestLogoutMutator_NoMetadata(t *testing.T) {
	// create web server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// prepare query
	ctx := context.Background()
	pLogout := postData{
		Query: `mutation {
			logout{
				success
			}
		}`,
	}

	// do query
	logoutResult := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schema(),
		RequestString:  pLogout.Query,
		VariableValues: pLogout.Variables,
		OperationName:  pLogout.Operation,
	})
	if !logoutResult.HasErrors() {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	// validate error
	err = logoutResult.Errors[0]
	if err.Error() != errUnauthorized.Error() {
		t.Errorf("unexpected error, got: '%s', want: '%s'.", err.Error(), errUnauthorized.Error())
	}
}

func TestLogoutMutator_Valid(t *testing.T) {
	// create web server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// prepare query
	ctx := context.Background()
	pLogin := postData{
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
		RequestString:  pLogin.Query,
		VariableValues: pLogin.Variables,
		OperationName:  pLogin.Operation,
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

	accessToken, accessTokenOK := login["accessToken"].(string)
	if !accessTokenOK {
		t.Errorf("no accessToken returned.")
		return
	}

	// extract metadata
	req := http.Request{}
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	metadata, err := server.extractTokenMetadata(&req)
	if err != nil {
		t.Errorf("unexpected error, got: %#v, want: nil.", err.Error())
	}

	// add metadata to context
	ctx = context.WithValue(ctx, metadataKey, metadata)
	pLogout := postData{
		Query: `mutation {
			logout{
				success
			}
		}`,
	}

	// do query
	logoutResult := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schema(),
		RequestString:  pLogout.Query,
		VariableValues: pLogout.Variables,
		OperationName:  pLogout.Operation,
	})
	if logoutResult.HasErrors() {
		for _, e := range logoutResult.Errors {
			t.Errorf("unexpected error, got: %#v, want: nil.", e.Error())
		}
		return
	}

	// validate data
	logoutData, logoutDataOk := logoutResult.Data.(map[string]interface{})
	if !logoutDataOk {
		t.Errorf("no data returned.")
		return
	}

	logout, logoutOK := logoutData["logout"].(map[string]interface{})
	if !logoutOK {
		t.Errorf("no logout data returned.")
		return
	}

	success, successOK := logout["success"].(bool)
	if !successOK {
		t.Errorf("no success returned.")
		return
	}
	if !success {
		t.Errorf("logout failed.")
	}
}

func TestRefreshAccessTokenMutator_Valid(t *testing.T) {
	// create web server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// prepare query
	ctx := context.Background()
	pLogin := postData{
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
		RequestString:  pLogin.Query,
		VariableValues: pLogin.Variables,
		OperationName:  pLogin.Operation,
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

	accessToken, accessTokenOK := login["accessToken"].(string)
	if !accessTokenOK {
		t.Errorf("no accessToken returned.")
		return
	}

	refreshToken, refreshTokenOK := login["refreshToken"].(string)
	if !refreshTokenOK {
		t.Errorf("no refreshToken returned.")
		return
	}

	// extract metadata
	req := http.Request{}
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	metadata, err := server.extractTokenMetadata(&req)
	if err != nil {
		t.Errorf("unexpected error, got: %#v, want: nil.", err.Error())
	}

	// add metadata to context
	ctx = context.WithValue(ctx, metadataKey, metadata)
	pRefreshAccessToken := postData{
		Query: `mutation (
			$refreshToken: String!
		){
			refreshAccessToken(
				refreshToken: $refreshToken
			){
				accessToken
				refreshToken
			}
		}`,
		Variables: map[string]interface{}{
			"refreshToken": refreshToken,
		},
	}

	// do query
	refreshResult := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schema(),
		RequestString:  pRefreshAccessToken.Query,
		VariableValues: pRefreshAccessToken.Variables,
		OperationName:  pRefreshAccessToken.Operation,
	})
	if refreshResult.HasErrors() {
		for _, e := range result.Errors {
			t.Errorf("unexpected error, got: %#v, want: nil.", e.Error())
		}
		return
	}

	// validate data
	refreshData, refreshDataOk := refreshResult.Data.(map[string]interface{})
	if !refreshDataOk {
		t.Errorf("no data returned.")
		return
	}

	refreshAccessToken, refreshAccessTokenOK := refreshData["refreshAccessToken"].(map[string]interface{})
	if !refreshAccessTokenOK {
		t.Errorf("no refreshAccessToken data returned.")
		return
	}

	_, ratAccessTokenOK := refreshAccessToken["accessToken"].(string)
	if !ratAccessTokenOK {
		t.Errorf("no accessToken returned.")
		return
	}

	_, ratRefreshTokenOK := refreshAccessToken["refreshToken"].(string)
	if !ratRefreshTokenOK {
		t.Errorf("no refreshToken returned.")
		return
	}
}
