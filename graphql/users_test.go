package graphql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"net/http"
	"testing"
)

func TestAddUserMutator(t *testing.T) {
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

	_, _, _, _, _, err = testDoAddUser(server, metadata, "testaddusermutator", "newpassword", []string{})
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// do login
	_, _, err = testDoLogin(server, "testaddusermutator", "newpassword")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
}

func TestChangePasswordMutator(t *testing.T) {
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

	_, _, _, _, _, err = testDoAddUser(server, metadata, "testchangepasswordmutator", "newpassword", []string{})
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// do login
	accessToken2, _, err := testDoLogin(server, "testchangepasswordmutator", "newpassword")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// extract metadata
	req2 := http.Request{}
	req2.Header = http.Header{}
	req2.Header.Set("Authorization", "Bearer "+accessToken2)
	metadata2, err := server.extractTokenMetadata(&req2)
	if err != nil {
		t.Errorf("unexpected error, got: %#v, want: nil.", err.Error())
	}

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata2)
	p := postData{
		Query: `mutation (
			$password: String!
		){
			changePassword(
				password: $password
			){
				success
			}
		}`,
		Variables: map[string]interface{}{
			"password": "aD1fferentPassword!",
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
		t.Errorf("unexpected error, got: %s, want: nil.", result.Errors[0].Error())
		return
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("no data returned")
		return
	}

	changePassword, changePasswordOk := data["changePassword"].(map[string]interface{})
	if !changePasswordOk {
		t.Errorf("no changePassword data returned")
		return
	}

	isSuccess, successOk := changePassword["success"].(bool)
	if !successOk {
		t.Errorf("no success data returned")
	}
	if isSuccess != true {
		t.Errorf("got invalid updatedAt, got: %v, want: true.", isSuccess)
	}

	// do login
	_, _, err = testDoLogin(server, "testchangepasswordmutator", "aD1fferentPassword!")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
}

func TestMeQuery(t *testing.T) {
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
	if meIDUUID != uuid.Must(uuid.Parse("8c504483-1e11-4243-b6c8-14499877a641")) {
		t.Errorf("got wrong id, got: %s, want: '8c504483-1e11-4243-b6c8-14499877a641'.", meID)
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

func TestUserQuery(t *testing.T) {
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

	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `query (
			$id: String!
		){
			user(
				id: $id
			){
				id
				username
				groups
				createdAt
				updatedAt
			}
		}`,
		Variables: map[string]interface{}{
			"id": "8c504483-1e11-4243-b6c8-14499877a641",
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
		t.Errorf("unexpected error, got: %s, want: nil.", result.Errors[0].Error())
		return
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		t.Errorf("no data returned")
		return
	}

	user, userOk := data["user"].(map[string]interface{})
	if !userOk {
		t.Errorf("no changePassword user returned")
		return
	}

	userID, userIDOk := user["id"].(string)
	if !userIDOk {
		t.Errorf("no id data returned")
	}
	if userID != "8c504483-1e11-4243-b6c8-14499877a641" {
		t.Errorf("got invalid updatedAt, got: %s, want: '8c504483-1e11-4243-b6c8-14499877a641'", userID)
	}

	userUsername, userUsernameOk := user["username"].(string)
	if !userUsernameOk {
		t.Errorf("no username data returned")
	}
	if userUsername != "admin" {
		t.Errorf("got invalid updatedAt, got: %s, want: 'user'", userUsername)
	}

	userGroups, userGroupsOk := user["groups"].([]interface{})
	if !userGroupsOk {
		t.Errorf("no groups data returned")
	}
	if len(userGroups) != 1 {
		t.Errorf("got invalid updatedAt, got: %d, want: 0'", len(userGroups))
	}

	createdAt, createdAtOk := user["createdAt"].(int)
	if !createdAtOk {
		t.Errorf("no me createdAt data returned")
	}
	if createdAt <= 0 {
		t.Errorf("got invalid createdAt, got: %d, want: >0.", createdAt)
	}

	updatedAt, updatedAtOk := user["updatedAt"].(int)
	if !updatedAtOk {
		t.Errorf("no me updatedAt data returned")
	}
	if updatedAt <= 0 {
		t.Errorf("got invalid updatedAt, got: %d, want: >0.", updatedAt)
	}
}

func testDoAddUser(server *Server, metadata *accessDetails, username, password string, groups []string) (string, string, []interface{}, int, int, error) {
	// prepare query
	ctx := context.WithValue(context.Background(), metadataKey, metadata)
	p := postData{
		Query: `mutation (
			$username: String!
			$password: String!
			$groups: [String!]
		){
			addUser(
				username: $username
				password: $password
				groups: $groups
			){
				id
				username
				groups
				createdAt
				updatedAt
			}
		}`,
		Variables: map[string]interface{}{
			"username": username,
			"password": password,
			"groups":   groups,
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
		return "", "", nil, 0, 0, result.Errors[0]
	}

	// validate data
	data, dataOk := result.Data.(map[string]interface{})
	if !dataOk {
		return "", "", nil, 0, 0, fmt.Errorf("no data returned")
	}

	addUser, addUserOK := data["addUser"].(map[string]interface{})
	if !addUserOK {
		return "", "", nil, 0, 0, fmt.Errorf("no addUser data returned")
	}

	id, idOK := addUser["id"].(string)
	if !idOK {
		return "", "", nil, 0, 0, fmt.Errorf("no id returned")
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return "", "", nil, 0, 0, fmt.Errorf("can't parse id %s: %s", id, err.Error())
	}

	newUsername, usernameOK := addUser["username"].(string)
	if !usernameOK {
		return "", "", nil, 0, 0, fmt.Errorf("no username returned")
	}

	newGroups, groupsOK := addUser["groups"].([]interface{})
	if !groupsOK {
		return "", "", nil, 0, 0, fmt.Errorf("no groups returned")
	}

	createdAt, createdAtOk := addUser["createdAt"].(int)
	if !createdAtOk {
		return "", "", nil, 0, 0, fmt.Errorf("no me createdAt data returned")
	}
	if createdAt <= 0 {
		return "", "", nil, 0, 0, fmt.Errorf("got invalid createdAt, got: %d, want: >0", createdAt)
	}

	updatedAt, updatedAtOk := addUser["updatedAt"].(int)
	if !updatedAtOk {
		return "", "", nil, 0, 0, fmt.Errorf("no me updatedAt data returned")
	}
	if updatedAt <= 0 {
		return "", "", nil, 0, 0, fmt.Errorf("got invalid updatedAt, got: %d, want: >0", updatedAt)
	}

	return id, newUsername, newGroups, createdAt, updatedAt, nil
}
