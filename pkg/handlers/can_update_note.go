package handlers

import (
	"context"
	"errors"
	subscription_pb "github.com/in-rich/proto/proto-go/subscription"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/pkg/models"
	"github.com/in-rich/uservice-subscription/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CanUpdateNoteHandler struct {
	subscription_pb.CanUpdateNoteServer
	service services.CanUpdateNoteService
}

func (h *CanUpdateNoteHandler) CanUpdateNote(ctx context.Context, in *subscription_pb.CanUpdateNoteRequest) (*subscription_pb.CanUpdateNoteResponse, error) {
	remainingEdits, err := h.service.Exec(ctx, &models.CanUpdateNoteRequest{
		Target:           in.GetTarget(),
		PublicIdentifier: in.GetPublicIdentifier(),
		AuthorID:         in.GetAuthorId(),
		ReadOnly:         in.GetReadOnly(),
	}, config.App.FreeTier, time.Now())

	if err != nil {
		if errors.Is(err, services.ErrNoteEditsExhausted) {
			return nil, status.Error(codes.ResourceExhausted, "note edits exhausted")
		}
		if errors.Is(err, services.ErrInvalidRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to check if note can be updated: %v", err)
	}

	return &subscription_pb.CanUpdateNoteResponse{
		RemainingEdits: int32(remainingEdits),
	}, nil
}

func NewCanUpdateNoteHandler(service services.CanUpdateNoteService) *CanUpdateNoteHandler {
	return &CanUpdateNoteHandler{
		service: service,
	}
}
