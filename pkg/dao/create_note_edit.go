package dao

import (
	"context"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/uptrace/bun"
)

type CreateNoteEditData struct {
	Target           entities.Target
	PublicIdentifier string
}

type CreateNoteEditRepository interface {
	CreateNoteEdit(ctx context.Context, author string, data *CreateNoteEditData) (*entities.NoteEdit, error)
}

type createNoteEditRepositoryImpl struct {
	db bun.IDB
}

func (r *createNoteEditRepositoryImpl) CreateNoteEdit(
	ctx context.Context, author string, data *CreateNoteEditData,
) (*entities.NoteEdit, error) {
	noteEdit := &entities.NoteEdit{
		PublicIdentifier: data.PublicIdentifier,
		Target:           data.Target,
		AuthorID:         author,
	}

	if _, err := r.db.NewInsert().Model(noteEdit).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return noteEdit, nil
}

func NewCreateNoteEditRepository(db bun.IDB) CreateNoteEditRepository {
	return &createNoteEditRepositoryImpl{
		db: db,
	}
}
