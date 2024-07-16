package models

type CanUpdateNoteRequest struct {
	Target           string `json:"target" validate:"omitempty,required_without=ReadOnly,oneof=company user"`
	PublicIdentifier string `json:"publicIdentifier" validate:"omitempty,required_without=ReadOnly,max=255"`
	AuthorID         string `json:"authorID" validate:"required,max=255"`
	ReadOnly         bool   `json:"read_only"`
}
