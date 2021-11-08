package graphql

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

func (s *Server) addRecordAMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to add A record")

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	// get domain
	domainIDStr, _ := params.Args["domainId"].(string)
	domainID, err := uuid.Parse(domainIDStr)
	if err != nil {
		return nil, err
	}
	domain, err := s.db.ReadDomain(domainID)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}
	if domain == nil {
		return nil, errors.New("domain not found")
	}

	// acl
	isOwner := domain.OwnerID == metadata.UserID
	isDNSAdmin := util.ContainsOneOfUUIDs(&metadata.Groups, &models.GroupsDNSAdmin)
	if !isOwner && !isDNSAdmin {
		return nil, errUnauthorized
	}

	// create record
	newRecord := models.Record{
		DomainID: domainID,
		Type:     models.RecordTypeA,
	}
	newRecord.Name, _ = params.Args["name"].(string)
	newRecord.Value, _ = params.Args["ip"].(string)
	ttl, ttlOk := params.Args["ttl"].(int)
	if ttlOk {
		newRecord.TTL = sql.NullInt32{
			Int32: int32(ttl),
			Valid: true,
		}
	}

	err = newRecord.Validate()
	if err != nil {
		return nil, err
	}

	err = s.db.Create(&newRecord)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	return newRecord, nil
}
