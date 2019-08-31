package users

import (
	"database/sql"
	"go-twitter-test/repositories/testutils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserRepository_Get(t *testing.T) {
	const dbDsn = "./testdata/test1.db"
	db := testutils.SetUp(t, dbDsn)
	defer testutils.TearDown(t, db, []string{dbDsn})

	loadFixtures(t, db)

	repo := New(db)
	user, err := repo.Get(1)

	require.Nil(t, err)
	require.Equal(t, &User{
		ID:    1,
		Email: "test@email.com",
	}, user)
}

func loadFixtures(t *testing.T, db *sql.DB) {
	_, err := db.Exec("INSERT INTO users (id, email) VALUES (?, ?)", 1, "test@email.com")
	require.Nil(t, err)
}
