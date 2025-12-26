package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const usersTableName = "users"

func newUsersQ(db *pgdb.DB) data.UsersQ {
	return &usersQ{
		db:  db,
		sql: squirrel.StatementBuilder,
	}
}

type usersQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *usersQ) GetByUsername(ctx context.Context, username string) (*data.User, error) {
	query := m.sql.Select("*").
		From(usersTableName).
		Where("username = ?", username).
		PlaceholderFormat(squirrel.Dollar)

	var result data.User
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *usersQ) Insert(ctx context.Context, user data.User) (*data.User, error) {
	clauses := structs.Map(user)
	query := m.sql.Insert(usersTableName).
		SetMap(clauses).
		Suffix("RETURNING *")

	var result data.User
	err := m.db.GetContext(ctx, &result, query)
	return &result, err
}

func (m *usersQ) UpdateRefreshToken(ctx context.Context, userID int64, refreshToken string) error {
	query := m.sql.Update(usersTableName).
		Set("refresh_token", refreshToken).
		Where("id = ?", userID).
		PlaceholderFormat(squirrel.Dollar)

	return m.db.ExecContext(ctx, query)
}