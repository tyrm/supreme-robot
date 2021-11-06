package graphql

import (
	"database/sql"
	"errors"
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
		OwnerID: metadata.UserId,
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
	ttl, _ := soaObj["ttl"].(int)
	newSoaRecord.TTL = sql.NullInt32{
		Int32: int32(ttl),
		Valid: true,
	}
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
	newSoaRecord.Retry = sql.NullInt32{
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
		return nil, errors.New(fmt.Sprintf("domain %s exists", newDomain.Domain))
	}

	// add domain to database
	err = s.db.CreateDomainWRecords(&newDomain, &newSoaRecord)
	if err != nil {
		logger.Errorf("db: %s", err.Error())
		return nil, err
	}

	return newDomain, nil
}

func (s *Server) domainQuery(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to get domain")

	// get id
	idStr, _ := params.Args["id"].(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// did user authenticate
	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)
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
	domains, err := s.db.ReadDomainsForUser(metadata.UserId)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}

	return domains, nil
}
