package postgres

import (
	"context"
	"errors"

	"github.com/dagulv/screener/internal/env"
	"github.com/webmafia/fast"
	"github.com/webmafia/pg"
)

func NewDB(ctx context.Context, env *env.Environment) (db *pg.DB, err error) {
	return pg.New(ctx, env.PostgresConn)
}

type db struct {
	*pg.DB
}

func (db db) AcquireContext(ctx context.Context, readOnly ...bool) (newCtx context.Context, err error) {
	return db.Transaction(ctx, readOnly...)
}

func (db db) CommitContext(ctx context.Context) error {
	if tx, ok := ctx.(*pg.Tx); ok {
		return tx.Commit(ctx)
	}

	return errors.New("no transaction context")
}

func (db db) ReleaseContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}

	if tx, ok := ctx.(*pg.Tx); ok {
		return tx.Release(ctx)
	}

	return errors.New("no transaction context")
}

type columns struct {
	cols []pg.ChainedIdentifier
}

// EncodeQuery implements QueryEncoder.
func (c columns) EncodeQuery(buf *fast.StringBuffer, _ *[]any) {
	writeIdentifiers(buf, c.cols)
}

func writeIdentifiers(b *fast.StringBuffer, ids []pg.ChainedIdentifier) {
	for i := range ids {
		if i != 0 {
			b.WriteByte(',')
		}

		ids[i].EncodeString(b)
	}
}
