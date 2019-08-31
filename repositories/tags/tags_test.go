package tags

import (
	"database/sql"
	"fmt"
	"go-twitter-test/sqlite"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsRepository(t *testing.T) {
	const dbDsn = "./testdata/test1.db"
	db, err := sqlite.New(dbDsn)
	require.Nil(t, err)

	defer tearDown(t, []string{dbDsn})
	defer closeDB(db)

	err = sqlite.LoadSchema(db)
	require.Nil(t, err)

	repo := New(db)
	tagID, err := repo.Put("A Nice Tag") // expected: a-nice-tag
	require.Nil(t, err)
	require.EqualValues(t, 1, tagID)

	tagID, err = repo.Put(" A-NICE tag   ") // expected: a-nice-tag
	require.Nil(t, err)
	require.EqualValues(t, 1, tagID) // same tag as before, ID should still be 1

	tagID, err = repo.Put("a different tag") // expected: a-different-tag
	require.Nil(t, err)
	require.EqualValues(t, 2, tagID)

	tagID, err = repo.GetID("a nice tag")
	require.Nil(t, err)
	require.EqualValues(t, 1, tagID)

	tagID, err = repo.GetID("a different tag")
	require.Nil(t, err)
	require.EqualValues(t, 2, tagID)
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(fmt.Errorf("could not close db: %v", err))
	}
}

func tearDown(t *testing.T, filenames []string) {
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
