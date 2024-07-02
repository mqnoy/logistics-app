package transaction

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type TransactionManager interface {
	AcquireTx(ctx context.Context) *gorm.DB
	CommitOrRollback(ctx context.Context, tx *gorm.DB, isRollback bool) error
}

type transactionManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (t *transactionManager) AcquireTx(ctx context.Context) *gorm.DB {
	log.Println("TransactionManager begin transactions")
	return t.db.WithContext(ctx).Begin()
}

func (t *transactionManager) CommitOrRollback(ctx context.Context, tx *gorm.DB, isRollback bool) error {
	if isRollback {
		log.Printf("TransactionManager rollback transactions")
		return t.Rollback(ctx, tx)
	}

	return t.Commit(ctx, tx)
}

func (t *transactionManager) Rollback(ctx context.Context, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Rollback().Error; err != nil {
		log.Printf("TransactionManager rollback error: %v", err)
	}

	log.Println("TransactionManager rollback success")
	return nil
}

func (t *transactionManager) Commit(ctx context.Context, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		log.Printf("TransactionManager commit error: %v", err)
	}

	log.Println("TransactionManager commit success")
	return nil
}
