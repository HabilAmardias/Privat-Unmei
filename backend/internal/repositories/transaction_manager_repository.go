package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/db"
	"privat-unmei/internal/logger"
)

type TxKey struct{}

type CustomTx struct {
	tx     *sql.Tx
	logger logger.CustomLogger
}

func (cdb *CustomTx) Query(query string, args ...any) (*sql.Rows, error) {
	cdb.logger.Infoln(query)
	return cdb.tx.Query(query, args...)
}

func (cdb *CustomTx) QueryRow(query string, args ...any) *sql.Row {
	cdb.logger.Infoln(query)
	return cdb.tx.QueryRow(query, args...)
}

func (cdb *CustomTx) Exec(query string, args ...any) (sql.Result, error) {
	cdb.logger.Infoln(query)
	return cdb.tx.Exec(query, args...)
}

func (cdb *CustomTx) Commit() error {
	return cdb.tx.Commit()
}

func (cdb *CustomTx) Rollback() error {
	return cdb.tx.Rollback()
}

func wrapTX(tx *sql.Tx, logger logger.CustomLogger) *CustomTx {
	return &CustomTx{tx, logger}
}

type RepoDriver interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func GetTransactionFromContext(ctx context.Context) *CustomTx {
	val := ctx.Value(TxKey{})
	tx, ok := val.(*CustomTx)
	if !ok {
		return nil
	}
	return tx
}

type TransactionManagerRepositories struct {
	DB     *db.CustomDB
	logger logger.CustomLogger
}

func CreateTransactionManager(db *db.CustomDB, logger logger.CustomLogger) *TransactionManagerRepositories {
	return &TransactionManagerRepositories{db, logger}
}

func (tr *TransactionManagerRepositories) WithTransaction(ctx context.Context, callable func(ctx context.Context) error) error {
	tx, err := tr.DB.Begin()
	if err != nil {
		return err
	}
	wrappedTx := wrapTX(tx, tr.logger)
	defer wrappedTx.Rollback()

	newCtx := context.WithValue(ctx, TxKey{}, wrappedTx)

	err = callable(newCtx)
	if err != nil {
		return err
	}
	return wrappedTx.Commit()
}
