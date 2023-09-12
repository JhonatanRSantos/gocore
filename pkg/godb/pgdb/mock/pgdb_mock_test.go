package mock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/JhonatanRSantos/gocore/pkg/godb/pgdb"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewClientSingleResult(t *testing.T) {
	ctx := context.Background()
	pgClient, err := NewClient()
	assert.NoError(t, err, "failed to create new pg client")
	defer pgClient.Close()

	newConnection := ClientConnectionMock{
		PgxConnIface: pgClient.connPool.AsConn(),
	}

	type User struct {
		UserID         string     `db:"user_id"`
		Name           string     `db:"name"`
		Email          string     `db:"email"`
		ConfirmedEmail bool       `db:"confirmed_email"`
		CreatedAt      time.Time  `db:"created_at"`
		UpdatedAt      time.Time  `db:"updated_at"`
		DeletedAt      *time.Time `db:"deleted_at"`
	}

	now := time.Now().UTC()
	newID := func() string {
		id, err := uuid.NewV4()
		assert.NoError(t, err)
		return id.String()
	}

	oneUser := User{
		UserID:         newID(),
		Name:           newID(),
		Email:          fmt.Sprintf("%s@email.com", newID()),
		ConfirmedEmail: true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	mockedRows := pgxmock.NewRows([]string{
		"user_id", "name", "email", "confirmed_email", "created_at", "updated_at", "deleted_at",
	})
	mockedRows.AddRow(
		oneUser.UserID,
		oneUser.Name,
		oneUser.Email,
		oneUser.ConfirmedEmail,
		oneUser.CreatedAt,
		oneUser.UpdatedAt,
		oneUser.DeletedAt,
	)
	newConnection.PgxConnIface.
		ExpectQuery("SELECT * FROM users WHERE = @email;").
		WithArgs(pgx.NamedArgs{"email": oneUser.Email}).
		WillReturnRows(mockedRows).RowsWillBeClosed()

	pgClient.AddConnection(newConnection)
	conn, err := pgClient.GetConnection(ctx)
	assert.NoError(t, err, "failed to get connection from pool")

	rows, err := conn.Query(
		ctx,
		"SELECT * FROM users WHERE = @email;", pgx.NamedArgs{"email": oneUser.Email},
	)
	assert.NoError(t, err, "failed to select all user data")

	user, err := pgdb.ParseRowTo[User](ctx, rows)
	assert.NoError(t, err, "failed to parse row data")
	assert.Equal(t, oneUser, *user)
}

func TestNewClientMultiResult(t *testing.T) {
	ctx := context.Background()
	pgClient, err := NewClient()
	assert.NoError(t, err, "failed to create new pg client")
	defer pgClient.Close()

	newConnection := ClientConnectionMock{
		PgxConnIface: pgClient.connPool.AsConn(),
	}
	defer newConnection.Release()
	_, errt := pgClient.GetConnection(ctx)
	assert.Error(t, errt)

	type User struct {
		UserID         string     `db:"user_id"`
		Name           string     `db:"name"`
		Email          string     `db:"email"`
		ConfirmedEmail bool       `db:"confirmed_email"`
		CreatedAt      time.Time  `db:"created_at"`
		UpdatedAt      time.Time  `db:"updated_at"`
		DeletedAt      *time.Time `db:"deleted_at"`
	}

	now := time.Now().UTC()
	newID := func() string {
		id, err := uuid.NewV4()
		assert.NoError(t, err)
		return id.String()
	}

	userOne := User{
		UserID:         newID(),
		Name:           newID(),
		Email:          fmt.Sprintf("%s@email.com", newID()),
		ConfirmedEmail: true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	userTwo := User{
		UserID:         newID(),
		Name:           newID(),
		Email:          fmt.Sprintf("%s@email.com", newID()),
		ConfirmedEmail: true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	mockedRows := pgxmock.NewRows([]string{
		"user_id", "name", "email", "confirmed_email", "created_at", "updated_at", "deleted_at",
	})
	mockedRows.AddRow(
		userOne.UserID,
		userOne.Name,
		userOne.Email,
		userOne.ConfirmedEmail,
		userOne.CreatedAt,
		userOne.UpdatedAt,
		userOne.DeletedAt,
	).AddRow(
		userTwo.UserID,
		userTwo.Name,
		userTwo.Email,
		userTwo.ConfirmedEmail,
		userTwo.CreatedAt,
		userTwo.UpdatedAt,
		userTwo.DeletedAt,
	)
	newConnection.PgxConnIface.
		ExpectQuery("SELECT * FROM users WHERE = @email;").
		WithArgs(pgx.NamedArgs{"email": userOne.Email}).
		WillReturnRows(mockedRows).RowsWillBeClosed()

	pgClient.AddConnection(newConnection)
	conn, err := pgClient.GetConnection(ctx)
	assert.NoError(t, err, "failed to get connection from pool")

	rows, err := conn.Query(
		ctx,
		"SELECT * FROM users WHERE = @email;", pgx.NamedArgs{"email": userOne.Email},
	)
	assert.NoError(t, err, "failed to select all user data")

	expectedUsers := []*User{&userOne, &userTwo}
	users, err := pgdb.ParseRowsTo[User](ctx, rows)
	assert.NoError(t, err, "failed to parse row data")
	assert.Equal(t, expectedUsers, users)
}
