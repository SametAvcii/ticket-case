package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrTransactionInProgress = errors.New("transaction already in progress")
	ErrTransactionNotStarted = errors.New("transaction not started")
)

type Transactor interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}

type transactor struct {
	db *gorm.DB
	tx *gorm.DB
}

func (t *transactor) Begin(ctx context.Context) error {
	if t.tx != nil {
		return ErrTransactionInProgress
	}
	tx := t.db.Begin()
	t.tx = tx
	return nil
}

func (t *transactor) Commit() error {
	if t.tx == nil {
		return ErrTransactionNotStarted
	}

	t.tx.Commit()
	return nil
}

func (t *transactor) Rollback() error {
	if t.tx == nil {
		return ErrTransactionNotStarted
	}

	t.tx.Rollback()
	return nil
}
