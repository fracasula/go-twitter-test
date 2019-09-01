package tags

import (
	"go-twitter-test/repositories/testutils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsRepository(t *testing.T) {
	const dbDsn = "./testdata/test1.db"
	db := testutils.SetUp(t, dbDsn)
	defer testutils.TearDown(t, db, []string{dbDsn})

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
