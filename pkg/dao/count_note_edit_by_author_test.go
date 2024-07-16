package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var countNoteEditByAuthorFixtures = []*entities.NoteEdit{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000005")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)),
	},
	// Different note
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetUser,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)),
	},
	// Different target
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetCompany,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)),
	},
	// Different author
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
		AuthorID:         "author-id-2",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)),
	},
}

func TestCountNoteEditByAuthor(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name      string
		authorID  string
		since     *time.Time
		expect    int
		expectErr error
	}{
		{
			name:     "CountNoteEditByAuthor",
			authorID: "author-id-1",
			since:    lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			expect:   4,
		},
		{
			name:     "CountNoteEditByAuthorSince",
			authorID: "author-id-1",
			since:    lo.ToPtr(time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)),
			expect:   2,
		},
		{
			name:     "CountNoteEditByAuthorNone",
			authorID: "author-id-1",
			since:    lo.ToPtr(time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC)),
			expect:   0,
		},
	}

	stx := BeginTX(db, getLatestNoteEditByAuthorFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCountNoteEditsByAuthorRepository(tx)
			count, err := repo.CountNoteEditsByAuthor(context.Background(), tt.authorID, tt.since)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, count)
		})
	}
}
