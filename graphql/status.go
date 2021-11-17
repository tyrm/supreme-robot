package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/version"
)

type status struct {
	Version string `json:"version"`
}

func (s *Server) statusQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to get status")

	newStatus := status{
		Version: version.Version,
	}

	return newStatus, nil
}
