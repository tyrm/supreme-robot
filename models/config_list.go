package models

type ConfigList []Config

// Model Functions

func (cl *ConfigList) Upsert(c *Client) error {
	// get current state
	k := make([]string, len(*cl))
	for i, cfg := range *cl {
		k[i] = cfg.Key
	}
	currentConfig, err := c.ReadConfigsByKeys(&k)

	// start transaction
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	var (
		index int
		found bool
	)
	for _, cfg := range *cl {
		index = 0
		found = false
		for i, cc := range *currentConfig {
			if cfg.Key == cc.Key {
				index = i
				found = true
				break
			}
		}

		var err error
		if found {
			// update
			logger.Tracef("updating config key %s: %s", cfg.Value, cfg.Key)
			_, err = tx.
				Exec(`UPDATE public.config SET value=$1 , updated_at=CURRENT_TIMESTAMP
				WHERE id=$2;`, cfg.Value, (*currentConfig)[index].ID)
		} else {
			// create
			logger.Tracef("adding config key %s: %s", cfg.Value, cfg.Key)
			_, err = tx.
				Exec(`INSERT INTO "public"."config"("key", "value")
			VALUES ($1, $2)`, cfg.Key, cfg.Value)
		}

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
	}

	// commit transaction
	logger.Tracef("committing config updated")
	err = tx.Commit()
	if err != nil {
		logger.Errorf("commit transaction: %s", err.Error())
		return err
	}

	return nil
}
