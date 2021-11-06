package graphql

import (
	"errors"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

func (s *Server) addUserMutation(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to add user")

	// acl
	if params.Context.Value(metadataKey) == nil { // did user authenticate
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	if !util.ContainsOneOfUUIDs(&models.GroupsUserAdmin, &metadata.Groups) {
		// user is not user admin
		logger.Tracef("user is not user admin")
		return nil, errUnauthorized
	}

	// marshall and cast the argument values
	username, _ := params.Args["username"].(string)
	password, _ := params.Args["password"].(string)
	groups, groupsOk := params.Args["groups"].([]interface{})

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
	logger.Debugf("trying to get me")

	// acl
	if params.Context.Value(metadataKey) == nil { // did user authenticate
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	return s.db.ReadUser(metadata.UserID)
}

func (s *Server) userQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to get user")

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
	logger.Tracef("%v(%s) %v(%s)", idOk, id)

	// acl
	if params.Context.Value(metadataKey) == nil { // did user authenticate
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	if !util.ContainsOneOfUUIDs(&models.GroupsUserAdmin, &metadata.Groups) {
		// user is not user admin
		logger.Tracef("user is not user admin")
		return nil, errUnauthorized
	}

	// find user
	user, err := s.db.ReadUser(id)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}
	return user, nil
}
