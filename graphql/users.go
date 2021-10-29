package graphql

import "github.com/graphql-go/graphql"

func (s *Server) userQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to get user")
	return "asdf", nil
}