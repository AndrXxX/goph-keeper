package postgressql

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vingarcia/ksql"
)

type ksqlProvider struct {
	mock.Mock
}

func (p *ksqlProvider) Insert(ctx context.Context, table ksql.Table, record interface{}) error {
	args := p.Called(ctx, table, record)
	return args.Error(0)
}

func (p *ksqlProvider) Patch(ctx context.Context, table ksql.Table, record interface{}) error {
	args := p.Called(ctx, table, record)
	return args.Error(0)
}
func (p *ksqlProvider) Delete(ctx context.Context, table ksql.Table, idOrRecord interface{}) error {
	args := p.Called(ctx, table, idOrRecord)
	return args.Error(0)
}

func (p *ksqlProvider) Query(ctx context.Context, records interface{}, query string, params ...interface{}) error {
	args := p.Called(ctx, records, query, params)
	return args.Error(0)
}
func (p *ksqlProvider) QueryOne(ctx context.Context, record interface{}, query string, params ...interface{}) error {
	args := p.Called(ctx, record, query, params)
	return args.Error(0)
}
func (p *ksqlProvider) QueryChunks(ctx context.Context, parser ksql.ChunkParser) error {
	args := p.Called(ctx, parser)
	return args.Error(0)
}

func (p *ksqlProvider) Exec(ctx context.Context, query string, params ...interface{}) (ksql.Result, error) {
	args := p.Called(ctx, query, params)
	return args.Get(0).(ksql.Result), args.Error(1)
}
func (p *ksqlProvider) Transaction(ctx context.Context, fn func(ksql.Provider) error) error {
	args := p.Called(ctx, fn)
	return args.Error(0)
}
