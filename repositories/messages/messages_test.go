package messages

import (
	"database/sql"
	"go-twitter-test/repositories/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMessagesRepository(t *testing.T) {
	const dbDsn = "./testdata/test1.db"
	db := testutils.SetUp(t, dbDsn)
	defer testutils.TearDown(t, db, []string{dbDsn})

	loadFixtures(t, db)

	now := time.Now()
	dateStart := now.Unix()
	dateEnd := now.Add(1 * time.Second).Unix()
	formattedDateStart := now.Format("2006-01-02T15:04:05")

	repo := New(db)
	err := repo.Create(MessageCreate{
		UserID:  1,
		TagID:   1,
		Message: "Message 1",
	})
	require.Nil(t, err)

	err = repo.Create(MessageCreate{
		UserID:  1,
		TagID:   2, // using tag 2 here to test the filter by tag later
		Message: "Message 2",
	})
	require.Nil(t, err)

	// sleeping 2 seconds to allow tests assertions by date range
	time.Sleep(2 * time.Second)

	err = repo.Create(MessageCreate{
		UserID:  1,
		TagID:   1,
		Message: "Message 3",
	})
	require.Nil(t, err)

	// testing CountMessages
	count, err := repo.CountMessages(0, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, 2, count)

	count, err = repo.CountMessages(1, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, 1, count)

	count, err = repo.CountMessages(2, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, 1, count)

	count, err = repo.CountMessages(0, 0, 0)
	require.Nil(t, err)
	require.EqualValues(t, 3, count)

	// testing GetMessages
	list, err := repo.GetMessages(0, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, []MessageList{
		{
			ID:        1,
			Message:   "Message 1",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-1",
		},
		{
			ID:        2,
			Message:   "Message 2",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-2",
		},
	}, list)

	list, err = repo.GetMessages(1, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, []MessageList{
		{
			ID:        1,
			Message:   "Message 1",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-1",
		},
	}, list)

	list, err = repo.GetMessages(2, dateStart, dateEnd)
	require.Nil(t, err)
	require.EqualValues(t, []MessageList{
		{
			ID:        2,
			Message:   "Message 2",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-2",
		},
	}, list)

	list, err = repo.GetMessages(0, 0, 0)
	require.Nil(t, err)
	require.EqualValues(t, []MessageList{
		{
			ID:        1,
			Message:   "Message 1",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-1",
		},
		{
			ID:        2,
			Message:   "Message 2",
			CreatedAt: formattedDateStart,
			UserEmail: "user@email.com",
			Tag:       "tag-2",
		},
		{
			ID:        3,
			Message:   "Message 3",
			CreatedAt: now.Add(2 * time.Second).Format("2006-01-02T15:04:05"),
			UserEmail: "user@email.com",
			Tag:       "tag-1",
		},
	}, list)
}

func loadFixtures(t *testing.T, db *sql.DB) {
	_, err := db.Exec("INSERT INTO users (id, email) VALUES (?, ?)", 1, "user@email.com")
	require.Nil(t, err)

	_, err = db.Exec("INSERT INTO main.tags (id, tag) VALUES (?, ?), (?, ?)", 1, "tag-1", 2, "tag-2")
	require.Nil(t, err)
}
