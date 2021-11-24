package worker

import (
	"context"
	"errors"
	"fmt"
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
)

func (w *Worker) addDomainHandler(ctx context.Context, args ...interface{}) error {
	help := faktory.HelperFor(ctx)
	idStr := args[0].(string)
	logger.Tracef("%s adding domain %s", help.Jid(), idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		msg := fmt.Sprintf("%s can't parsing uuid %s: %s", help.Jid(), idStr, err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}

	domain, err := w.db.ReadDomain(id)
	if err != nil {
		msg := fmt.Sprintf("%s error getting domain id %s: %s", help.Jid(), idStr, err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}
	if domain == nil {
		msg := fmt.Sprintf("%s domain not found for id %s", help.Jid(), idStr)
		logger.Warningf(msg)
		return errors.New(msg)
	}

	logger.Tracef("%s adding domain %s to domain list", help.Jid(), domain.Domain)
	err = w.kv.AddDomain(domain.Domain)
	if err != nil {
		msg := fmt.Sprintf("%s error adding domain to redis: %s", help.Jid(), err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}

	logger.Tracef("%s finished adding domain %s", help.Jid(), idStr)
	return nil
}

func (w *Worker) removeDomainHandler(ctx context.Context, args ...interface{}) error {
	help := faktory.HelperFor(ctx)
	idStr := args[0].(string)
	logger.Tracef("%s removing domain %s", help.Jid(), idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		msg := fmt.Sprintf("%s can't parsing uuid %s: %s", help.Jid(), idStr, err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}

	domain, err := w.db.ReadDomainZ(id)
	if err != nil {
		msg := fmt.Sprintf("%s error getting domain id %s: %s", help.Jid(), idStr, err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}
	if domain == nil {
		msg := fmt.Sprintf("%s domain not found for id %s", help.Jid(), idStr)
		logger.Warningf(msg)
		return errors.New(msg)
	}

	logger.Tracef("%s adding domain %s to domain list", help.Jid(), domain.Domain)
	err = w.kv.AddDomain(domain.Domain)
	if err != nil {
		msg := fmt.Sprintf("%s error adding domain to redis: %s", help.Jid(), err.Error())
		logger.Warningf(msg)
		return errors.New(msg)
	}

	logger.Tracef("%s finished removing domain %s", help.Jid(), idStr)
	return nil
}

func (w *Worker) updateSubDomainHandler(ctx context.Context, args ...interface{}) error {
	help := faktory.HelperFor(ctx)
	id := args[0].(string)
	subdomain := args[1].(string)
	logger.Tracef("%s updating sub domain %s for domain %s", help.Jid(), subdomain, id)

	return nil
}
