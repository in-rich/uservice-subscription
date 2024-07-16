package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type NoteEdit struct {
	bun.BaseModel `bun:"table:note_edits"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	AuthorID string `bun:"author_id,notnull"`

	PublicIdentifier string `bun:"public_identifier,notnull"`
	Target           Target `bun:"target,notnull"`

	CreatedAt *time.Time `bun:"created_at,notnull"`
}
