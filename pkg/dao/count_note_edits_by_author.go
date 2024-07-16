package dao

import (
	"context"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/uptrace/bun"
	"time"
)

type CountNoteEditsByAuthorRepository interface {
	CountNoteEditsByAuthor(ctx context.Context, author string, since *time.Time) (int, error)
}

type countNoteEditsByAuthorRepositoryImpl struct {
	db bun.IDB
}

func (r *countNoteEditsByAuthorRepositoryImpl) CountNoteEditsByAuthor(ctx context.Context, author string, since *time.Time) (int, error) {
	var count int

	count, err := r.db.NewSelect().
		Model((*entities.NoteEdit)(nil)).
		Where("author_id = ?", author).
		Where("created_at >= ?", since).
		Count(ctx)

	return count, err
}

func NewCountNoteEditsByAuthorRepository(db bun.IDB) CountNoteEditsByAuthorRepository {
	return &countNoteEditsByAuthorRepositoryImpl{
		db: db,
	}
}
