package graphql

import (
	"errors"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

func (s *Server) addUserMutation(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to add user")

	// marshall and cast the argument values
	username, _ := params.Args["username"].(string)
	password, _ := params.Args["password"].(string)
	groups, groupsOk := params.Args["groups"].([]interface{})

	// acl
	if params.Context.Value(MetadataKey) == nil { // did user authenticate
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	if !util.ContainsOneOfUUIDs(&models.GroupsUserAdmin, &metadata.Groups) {
		// user is not user admin
		logger.Tracef("user is not user admin")
		return nil, ErrUnauthorized
	}

	// validate inputs
	groupUUIDs := make([]uuid.UUID, len(groups))
	if groupsOk {
		logger.Tracef("groups found")
		for i, str := range groups {
			u, err := uuid.Parse(str.(string))
			if err != nil {
				return nil, err
			}
			groupUUIDs[i] = u
		}
	}

	// check if user exists
	user, err := s.db.ReadUserByUsername(username)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}
	if user != nil {
		return nil, errors.New("username taken")
	}

	// create user
	newUser := models.User{Username: username}
	err = newUser.SetPassword(password)
	if err != nil {
		logger.Errorf("setting password: %s", err.Error())
		return nil, err
	}

	// add user to database
	err = newUser.Create(s.db)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	// add groups
	if groupsOk {
		err = newUser.AddGroup(s.db, groupUUIDs...)
		if err != nil {
			return nil, err
		}
	}

	return newUser, nil
}

func (s *Server) meQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to get me")

	// acl
	if params.Context.Value(MetadataKey) == nil { // did user authenticate
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	return s.db.ReadUser(metadata.UserId)
}

func (s *Server) userQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to get user")

	// marshall and cast the argument values
	idStr, idOk := params.Args["id"].(string)
	id := uuid.Nil
	if idOk {
		var err error
		id, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
	}
	username, usernameOk := params.Args["username"].(string)
	logger.Tracef("%v(%s) %v(%s)", idOk, id, usernameOk, username)

	// acl
	if params.Context.Value(MetadataKey) == nil { // did user authenticate
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	if !util.ContainsOneOfUUIDs(&models.GroupsUserAdmin, &metadata.Groups) {
		// user is not user admin
		logger.Tracef("user is not user admin")
		return nil, ErrUnauthorized
	}

	return nil, nil
}