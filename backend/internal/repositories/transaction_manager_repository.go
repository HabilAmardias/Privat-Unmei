package repositories

import (
	"context"
	"database/sql"
)

type TxKey struct{}

type RepoDriver interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func GetTransactionFromContext(ctx context.Context) *sql.Tx {
	val := ctx.Value(TxKey{})
	tx, ok := val.(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

type TransactionManagerRepositories struct {
	DB *sql.DB
}

func CreateTransactionManager(db *sql.DB) *TransactionManagerRepositories {
	return &TransactionManagerRepositories{db}
}

func (tr *TransactionManagerRepositories) WithTransaction(ctx context.Context, callable func(ctx context.Context) error) error {
	tx, err := tr.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	newCtx := context.WithValue(ctx, TxKey{}, tx)

	err = callable(newCtx)
	if err != nil {
		return err
	}
	return nil
}
