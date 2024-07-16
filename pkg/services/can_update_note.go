package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/in-rich/uservice-subscription/pkg/models"
	"github.com/samber/lo"
	"time"
)

var (
	// NoteEditBufferTime is the time window in which a user can edit a note, without it counting as a new edit.
	NoteEditBufferTime = 60 * time.Minute
)

type CanUpdateNoteService interface {
	Exec(
		ctx context.Context,
		canUpdateRequest *models.CanUpdateNoteRequest,
		tier config.TierInformation,
		now time.Time,
	) (int, error)
}

type canUpdateNoteServiceImpl struct {
	countEditsRepository    dao.CountNoteEditsByAuthorRepository
	createEditRepository    dao.CreateNoteEditRepository
	getLatestEditRepository dao.GetLatestNoteEditByAuthorRepository
}

func (s *canUpdateNoteServiceImpl) Exec(
	ctx context.Context,
	canUpdateRequest *models.CanUpdateNoteRequest,
	tier config.TierInformation,
	now time.Time,
) (int, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(canUpdateRequest); err != nil {
		return 0, errors.Join(ErrInvalidRequest, err)
	}

	editsSince := now.UTC().Add(-*tier.Notes.CountEditsOver)
	editsCount, err := s.countEditsRepository.CountNoteEditsByAuthor(ctx, canUpdateRequest.AuthorID, &editsSince)
	if err != nil {
		return 0, fmt.Errorf("count note edits: %w", err)
	}

	// Avoid discrepancies if the limit has been overflowed.
	remainingEdits := lo.Max([]int{tier.Notes.MaxEdits - editsCount, 0})
	// Don't throw in read only mode.
	if canUpdateRequest.ReadOnly {
		return remainingEdits, nil
	}

	latestEditForNote, err := s.getLatestEditRepository.GetLatestNoteEditByAuthor(
		ctx,
		canUpdateRequest.AuthorID,
		entities.Target(canUpdateRequest.Target),
		canUpdateRequest.PublicIdentifier,
	)
	if err != nil && !errors.Is(err, dao.ErrNoNoteEditFound) {
		return 0, fmt.Errorf("get latest note edit: %w", err)
	}

	// Edit is recent, nothing to do.
	if latestEditForNote != nil && latestEditForNote.CreatedAt.After(now.UTC().Add(-NoteEditBufferTime)) {
		return remainingEdits, nil
	}

	// No more edits remaining, throw.
	if remainingEdits == 0 {
		return 0, ErrNoteEditsExhausted
	}

	// Create a new note edit.
	_, err = s.createEditRepository.CreateNoteEdit(ctx, canUpdateRequest.AuthorID, &dao.CreateNoteEditData{
		Target:           entities.Target(canUpdateRequest.Target),
		PublicIdentifier: canUpdateRequest.PublicIdentifier,
	})
	if err != nil {
		return 0, fmt.Errorf("create note edit: %w", err)
	}

	return remainingEdits - 1, nil
}

func NewCanUpdateNoteService(
	countEditsRepository dao.CountNoteEditsByAuthorRepository,
	createEditRepository dao.CreateNoteEditRepository,
	getLatestEditRepository dao.GetLatestNoteEditByAuthorRepository,
) CanUpdateNoteService {
	return &canUpdateNoteServiceImpl{
		countEditsRepository:    countEditsRepository,
		createEditRepository:    createEditRepository,
		getLatestEditRepository: getLatestEditRepository,
	}
}
