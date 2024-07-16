package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/uptrace/bun"
)

type GetLatestNoteEditByAuthorRepository interface {
	GetLatestNoteEditByAuthor(ctx context.Context, author string, target entities.Target, publicIdentifier string) (*entities.NoteEdit, error)
}

type getLatestNoteEditByAuthorRepositoryImpl struct {
	db bun.IDB
}

func (r *getLatestNoteEditByAuthorRepositoryImpl) GetLatestNoteEditByAuthor(
	ctx context.Context, author string, target entities.Target, publicIdentifier string,
) (*entities.NoteEdit, error) {
	noteEdit := new(entities.NoteEdit)

	err := r.db.NewSelect().
		Model(noteEdit).
		Where("author_id = ?", author).
		Where("target = ?", target).
		Where("public_identifier = ?", publicIdentifier).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoNoteEditFound
		}

		return nil, err
	}

	return noteEdit, nil
}

func NewGetLatestNoteEditByAuthorRepository(db bun.IDB) GetLatestNoteEditByAuthorRepository {
	return &getLatestNoteEditByAuthorRepositoryImpl{
		db: db,
	}
}
