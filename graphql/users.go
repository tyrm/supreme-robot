package graphql

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

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