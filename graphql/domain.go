package graphql

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
)

func (s *Server) addDomainMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to add domain")

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	// collect vars
	newDomain := models.Domain{
		OwnerID: metadata.UserID,
	}
	newDomain.Domain, _ = params.Args["domain"].(string)

	newSoaRecord := models.Record{
		Name:  "@",
		Type:  "SOA",
		Value: s.primaryNS,
	}
	soaObj, _ := params.Args["soa"].(map[string]interface{})
	mbox, _ := soaObj["mbox"].(string)
	newSoaRecord.MBox = sql.NullString{
		String: mbox,
		Valid:  true,
	}
	newSoaRecord.TTL, _ = soaObj["ttl"].(int)
	refresh, _ := soaObj["refresh"].(int)
	newSoaRecord.Refresh = sql.NullInt32{
		Int32: int32(refresh),
		Valid: true,
	}
	retry, _ := soaObj["retry"].(int)
	newSoaRecord.Retry = sql.NullInt32{
		Int32: int32(retry),
		Valid: true,
	}
	expire, _ := soaObj["expire"].(int)
	newSoaRecord.Expire = sql.NullInt32{
		Int32: int32(expire),
		Valid: true,
	}
	logger.Tracef("domain: %s soa: %v", newDomain.Domain, soaObj)

	// check for domain
	d, err := s.db.ReadDomainByDomain(newDomain.Domain)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}
	if d != nil {
		return nil, fmt.Errorf("domain %s exists", newDomain.Domain)
	}

	// add domain to database
	err = s.db.CreateDomainWRecords(&newDomain, &newSoaRecord)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	err = s.scheduler.AddDomain(newDomain.ID)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	return newDomain, nil
}

func (s *Server) deleteDomainMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to delete domain")

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	// get id
	idStr, _ := params.Args["id"].(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// do query
	domain, err := s.db.ReadDomain(id)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}
	if domain == nil {
		return nil, nil
	}

	// acl
	isOwner := domain.OwnerID == metadata.UserID
	isDNSAdmin := util.ContainsOneOfUUIDs(&metadata.Groups, &models.GroupsDNSAdmin)
	if !isOwner && !isDNSAdmin {
		return nil, errUnauthorized
	}

	err = s.db.Delete(domain)
	if err != nil {
		return nil, err
	}

	err = s.scheduler.RemoveDomain(domain.ID)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	return success{Success: true}, nil
}

func (s *Server) domainQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to get domain")

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	// get id
	idStr, _ := params.Args["id"].(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// do query
	domain, err := s.db.ReadDomain(id)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}
	if domain == nil {
		return nil, nil
	}

	domain.Records, err = s.db.ReadRecordsForDomain(domain.ID, "name", true)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}

	// does user own domain
	if domain.OwnerID == metadata.UserID {
		return domain, nil
	}

	// is user a dns admin
	if util.ContainsOneOfUUIDs(&metadata.Groups, &models.GroupsDNSAdmin) {
		return domain, nil
	}

	return nil, errUnauthorized
}

func (s *Server) myDomainsQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to get my domains")

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
	logger.Tracef("metadata: %v", metadata)

	// do query
	domains, err := s.db.ReadDomainsForUser(metadata.UserID)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}

	return domains, nil
}
