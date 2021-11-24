package graphql

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
	"testing"
)

var testUserAdminAccessToken = ""
var testUserAdminRefreshToken = ""

func TestLoginMutator_BadPassword(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, refreshToken, err := testDoLogin(server, "admin", "badpassword")
	if err.Error() != errBadLogin.Error() {
		t.Errorf("unexpected error, got: %s, want: %s.", err.Error(), errBadLogin.Error())
	}
	if accessToken != "" {
		t.Errorf("returned access token bit shouldn't")
	}
	if refreshToken != "" {
		t.Errorf("returned refresh token bit shouldn't")
	}
}

func TestLoginMutator_BadUsername(t *testing.T) {
	// create web server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, refreshToken, err := testDoLogin(server, "notauser", "password")
	if err.Error() != errBadLogin.Error() {
		t.Errorf("unexpected error, got: %s, want: %s.", err.Error(), errBadLogin.Error())
	}
	if accessToken != "" {
		t.Errorf("returned access token bit shouldn't")
	}
	if refreshToken != "" {
		t.Errorf("returned refresh token bit shouldn't")
	}
}

func TestLoginMutator_ValidLogin(t *testing.T) {
	// create server
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, refreshToken, err := testDoLogin(server, "admin", "password")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if accessToken == "" {
		t.Errorf("no access token returned")
	}
	if refreshToken == "" {
		t.Errorf("no refresh token returned")
	}
}

func TestLogoutMutator_NoMetadata(t *testing.T) {
	// create web server
	server, err := newTestServer()
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
	server, err := newTestServer()
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

	// add metadata to context
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
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
	server, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// do login
	accessToken, refreshToken, err := testDoLogin(server, "admin", "password")
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

	// add metadata to context
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
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
		for _, e := range refreshResult.Errors {
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

func testDoLogin(server *Server, username, password string) (string, string, error) {
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
			"username": username,
			"password": password,
		},
	}

	// do query
	result := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         server.schemaUnauthorized(),
		RequestString:  pLogin.Query,
		VariableValues: pLogin.Variables,
		OperationName:  pLogin.Operation,
	})
	if result.HasErrors() {
		return "", "", result.Errors[0]
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		return "", "", fmt.Errorf("no data returned")
	}

	login, loginOK := data["login"].(map[string]interface{})
	if !loginOK {
		return "", "", fmt.Errorf("no login data returned")
	}

	accessToken, accessTokenOK := login["accessToken"].(string)
	if !accessTokenOK {
		return "", "", fmt.Errorf("no accessToken returned")
	}

	refreshToken, refreshTokenOK := login["refreshToken"].(string)
	if !refreshTokenOK {
		return "", "", fmt.Errorf("no refreshToken returned")
	}

	return accessToken, refreshToken, nil
}

func testDoLoginAdmin(server *Server) (string, string, error) {
	if testUserAdminAccessToken != "" && testUserAdminRefreshToken != "" {
		return testUserAdminAccessToken, testUserAdminRefreshToken, nil
	}

	newAccessToken, newRefreshToken, err := testDoLogin(server, "admin", "password")
	if err != nil {
		return "", "", err
	}
	testUserAdminAccessToken = newAccessToken
	testUserAdminRefreshToken = newRefreshToken

	return testUserAdminAccessToken, testUserAdminRefreshToken, nil
}
