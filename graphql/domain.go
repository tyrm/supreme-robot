package graphql

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

func (s *Server) addDomainMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to add domain")

	// get id
	domainStr, _ := params.Args["domain"].(string)
	soaObj, _ := params.Args["soa"].(map[string]interface{})

	logger.Tracef("%s: %v", domainStr, soaObj)

	// did user authenticate
	if params.Context.Value(MetadataKey) == nil {
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	return nil, nil
}

func (s *Server) domainQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to get domain")

	// get id
	idStr, _ := params.Args["id"].(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// did user authenticate
	if params.Context.Value(MetadataKey) == nil {
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	// do query
	domain, err := s.db.ReadDomain(id)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}
	if domain == nil {
		return nil, nil
	}

	domain.Records, err = domain.GetRecords(s.db)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}

	// does user own domain
	if domain.OwnerID == metadata.UserId {
		return domain, nil
	}

	// is user a dns admin
	if util.ContainsOneOfUUIDs(&metadata.Groups, &models.GroupsDnsAdmin) {
		return domain, nil
	}

	return nil, ErrUnauthorized
}

func (s *Server) myDomainsQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Tracef("trying to get my domains")

	// did user authenticate
	if params.Context.Value(MetadataKey) == nil {
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)
	logger.Tracef("metadata: %v", metadata)

	// do query
	domains, err := s.db.ReadDomainsForUser(metadata.UserId)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}

	return domains, nil
}