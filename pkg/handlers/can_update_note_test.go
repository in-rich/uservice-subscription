package handlers_test

import (
	"context"
	"errors"
	subscription_pb "github.com/in-rich/proto/proto-go/subscription"
	"github.com/in-rich/uservice-subscription/pkg/handlers"
	"github.com/in-rich/uservice-subscription/pkg/services"
	servicesmocks "github.com/in-rich/uservice-subscription/pkg/services/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestCanUpdateNote(t *testing.T) {
	testData := []struct {
		name string

		in *subscription_pb.CanUpdateNoteRequest

		serviceResp int
		serviceErr  error

		expect     *subscription_pb.CanUpdateNoteResponse
		expectCode codes.Code
	}{
		{
			name: "CanUpdateNote",
			in: &subscription_pb.CanUpdateNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			serviceResp: 1,
			expect: &subscription_pb.CanUpdateNoteResponse{
				RemainingEdits: 1,
			},
		},
		{
			name: "NoteEditsExhausted",
			in: &subscription_pb.CanUpdateNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			serviceErr: services.ErrNoteEditsExhausted,
			expectCode: codes.ResourceExhausted,
		},
		{
			name: "InvalidRequest",
			in: &subscription_pb.CanUpdateNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			serviceErr: services.ErrInvalidRequest,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &subscription_pb.CanUpdateNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			serviceErr: errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockCanUpdateNoteService(t)
			service.On("Exec", context.TODO(), mock.Anything, mock.Anything, mock.Anything).Return(tt.serviceResp, tt.serviceErr)

			handler := handlers.NewCanUpdateNoteHandler(service)

			resp, err := handler.CanUpdateNote(context.TODO(), tt.in)

			RequireGRPCCodesEqual(t, err, tt.expectCode)
			require.Equal(t, tt.expect, resp)
		})
	}
}
