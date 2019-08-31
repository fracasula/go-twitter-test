package testutils

import (
	"database/sql"
	"go-twitter-test/sqlite"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetUp(t *testing.T, dbDsn string) *sql.DB {
	db, err := sqlite.New(dbDsn)
	require.Nil(t, err)

	err = sqlite.LoadSchema(db)
	require.Nil(t, err)

	return db
}

func TearDown(t *testing.T, db *sql.DB, filenames []string) {
	if err := db.Close(); err != nil {
		t.Fatalf("Could not close db: %v", err)
	}

	removeFiles(t, filenames)
}

func removeFiles(t *testing.T, filenames []string) {
	for _, filename := range filenames {
		if _, err := os.Stat(filename); err == nil {
			if err = os.Remove(filename); err != nil {
				t.Errorf("Could not remove file %q on teardown: %v", filename, err)
			}
		} else if !os.IsNotExist(err) {
			t.Errorf("Teardown failed with filename %q: %v", filename, err)
		}
	}
}
