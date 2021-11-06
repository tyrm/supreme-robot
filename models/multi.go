package models

// CreateDomainWRecords will create a domain and it's records in a single database transaction.
func (c *Client) CreateDomainWRecords(domain *Domain, records ...*Record) error {
	// start transaction
	logger.Tracef("starting transaction")
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	// add
	logger.Tracef("tx: add domain %s", domain.Domain)
	err = tx.
		QueryRow(`INSERT INTO "public"."domains"("domain", "owner_id")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, domain.Domain, domain.OwnerID).
		Scan(&domain.ID, &domain.CreatedAt, &domain.UpdatedAt)
	logger.Tracef("tx: domain %s added, got domain id %s", domain.Domain, domain.ID.String())

	// rollback on error
	if err != nil {
		logger.Errorf("tx error: %s", err.Error())
		rberr := tx.Rollback()
		if rberr != nil {
			logger.Errorf("rollback error: %s", rberr.Error())
			// something went REALLY wrong
			return rberr
		}
		return err
	}

	logger.Tracef("tx: add domain %s", domain.Domain)
	recordList := make([]Record, len(records))
	for i, r := range records {
		logger.Tracef("tx: add %s record %s", r.Type, r.Name)
		// add
		r.DomainID = domain.ID
		err = tx.
			QueryRow(`INSERT INTO "public"."domain_records"("name", "domain_id", "type", "value", "ttl",
            "priority", "port", "weight", "refresh", "retry", "expire", "mbox", "tag")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at, updated_at;`,
				r.Name, r.DomainID, r.Type, r.Value, r.TTL, r.Priority, r.Port, r.Weight, r.Refresh, r.Retry, r.Expire,
				r.MBox, r.Tag).
			Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
		logger.Tracef("tx: %s record %s added, got domain id %s", r.Type, r.Name, r.ID.String())

		// rollback on error
		if err != nil {
			logger.Errorf("tx error: %s", err.Error())
			rberr := tx.Rollback()
			if rberr != nil {
				logger.Errorf("rollback error: %s", rberr.Error())
				// something went REALLY wrong
				return rberr
			}
			return err
		}

		logger.Tracef("tx: add %s record %s to domain %s record list", r.Type, r.Name, domain.Domain)
		recordList[i] = *r
	}
	domain.Records = &recordList

	// commit transaction
	logger.Tracef("committing transaction")
	err = tx.Commit()
	if err != nil {
		logger.Errorf("can't commit transaction: %s", err.Error())
		return err
	}

	return nil
}
