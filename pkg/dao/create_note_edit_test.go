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

var createNoteEditFixtures = []*entities.NoteEdit{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestCreateNoteEdit(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name      string
		authorID  string
		data      *dao.CreateNoteEditData
		expect    *entities.NoteEdit
		expectErr error
	}{
		{
			name:     "CreateNoteEdit",
			authorID: "author-id-1",
			data: &dao.CreateNoteEditData{
				PublicIdentifier: "public-identifier-2",
				Target:           entities.TargetUser,
			},
			expect: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-2",
				Target:           entities.TargetUser,
			},
		},
		{
			name:     "CreateNoteEdit/SameNote",
			authorID: "author-id-1",
			data: &dao.CreateNoteEditData{
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetUser,
			},
			expect: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetUser,
			},
		},
	}

	stx := BeginTX(db, createNoteEditFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateNoteEditRepository(tx)
			noteEdit, err := repo.CreateNoteEdit(context.TODO(), tt.authorID, tt.data)

			if noteEdit != nil {
				// Since ID and CreatedAt are random, nullify them for comparison.
				noteEdit.ID = nil
				noteEdit.CreatedAt = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, noteEdit)
		})
	}
}
