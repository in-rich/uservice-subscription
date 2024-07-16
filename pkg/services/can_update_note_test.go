package services_test

import (
	"context"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	daomocks "github.com/in-rich/uservice-subscription/pkg/dao/mocks"
	"github.com/in-rich/uservice-subscription/pkg/entities"
	"github.com/in-rich/uservice-subscription/pkg/models"
	"github.com/in-rich/uservice-subscription/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCanUpdateNote(t *testing.T) {
	testData := []struct {
		name string

		data *models.CanUpdateNoteRequest
		now  time.Time
		tier config.TierInformation

		shouldCallCountNote bool
		countNoteResponse   int
		countNoteErr        error

		shouldCallLatestNote bool
		latestNoteResponse   *entities.NoteEdit
		latestNoteErr        error

		shouldCallCreateNote bool
		createNoteErr        error

		expect    int
		expectErr error
	}{
		// Success cases.
		{
			name: "CanUpdateNote/NewEdit",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    3,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateNote: true,
			expect:               1,
		},
		{
			name: "CanUpdateNote/NewEdit/NoEditRemaining",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    4,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateNote: true,
			expect:               0,
		},
		{
			name: "CanUpdateNote/RecentEdit",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    3,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 23, 30, 0, 0, time.UTC)),
			},
			expect: 2,
		},
		{
			// You are still allowed to continue edit a recent note, if you just reached your maximum edit count.
			name: "CanUpdateNote/RecentEdit/EditsExhausted",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    5,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 23, 30, 0, 0, time.UTC)),
			},
			expect: 0,
		},
		{
			name: "CanUpdateNote/ReadOnly",
			data: &models.CanUpdateNoteRequest{
				AuthorID: "author-id-1",
				ReadOnly: true,
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote: true,
			countNoteResponse:   3,
			expect:              2,
		},
		{
			name: "CanUpdateNote/ReadOnly/EditsExhausted",
			data: &models.CanUpdateNoteRequest{
				AuthorID: "author-id-1",
				ReadOnly: true,
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote: true,
			countNoteResponse:   5,
			expect:              0,
		},

		// Local error cases.
		{
			name: "CanUpdateNote/EditsExhausted",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    5,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectErr: services.ErrNoteEditsExhausted,
		},
		{
			name:      "CanUpdateNote/InvalidRequest",
			data:      &models.CanUpdateNoteRequest{},
			expectErr: services.ErrInvalidRequest,
		},

		// Dependency error cases.
		{
			name: "CreateNoteError",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    3,
			shouldCallLatestNote: true,
			latestNoteResponse: &entities.NoteEdit{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				CreatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateNote: true,
			createNoteErr:        FooErr,
			expectErr:            FooErr,
		},
		{
			name: "GetLatestNoteError",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote:  true,
			countNoteResponse:    3,
			shouldCallLatestNote: true,
			latestNoteErr:        FooErr,
			expectErr:            FooErr,
		},
		{
			name: "CountNotesError",
			data: &models.CanUpdateNoteRequest{
				AuthorID:         "author-id-1",
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
			},
			now: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			tier: config.TierInformation{
				Notes: config.NoteTierInformation{
					CountEditsOver: lo.ToPtr(24 * time.Hour),
					MaxEdits:       5,
				},
			},
			shouldCallCountNote: true,
			countNoteErr:        FooErr,
			expectErr:           FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			countNoteRepository := daomocks.NewMockCountNoteEditsByAuthorRepository(t)
			latestNoteRepository := daomocks.NewMockGetLatestNoteEditByAuthorRepository(t)
			createNoteRepository := daomocks.NewMockCreateNoteEditRepository(t)

			if tt.shouldCallCountNote {
				countNoteRepository.
					On("CountNoteEditsByAuthor", context.TODO(), tt.data.AuthorID, lo.ToPtr(tt.now.UTC().Add(-*tt.tier.Notes.CountEditsOver))).
					Return(tt.countNoteResponse, tt.countNoteErr)
			}

			if tt.shouldCallLatestNote {
				latestNoteRepository.
					On(
						"GetLatestNoteEditByAuthor",
						context.TODO(),
						tt.data.AuthorID,
						entities.Target(tt.data.Target),
						tt.data.PublicIdentifier,
					).
					Return(tt.latestNoteResponse, tt.latestNoteErr)
			}

			if tt.shouldCallCreateNote {
				createNoteRepository.
					On(
						"CreateNoteEdit",
						context.TODO(),
						tt.data.AuthorID,
						&dao.CreateNoteEditData{
							Target:           entities.Target(tt.data.Target),
							PublicIdentifier: tt.data.PublicIdentifier,
						},
					).
					Return(nil, tt.createNoteErr)
			}

			service := services.NewCanUpdateNoteService(
				countNoteRepository,
				createNoteRepository,
				latestNoteRepository,
			)

			remainingEdits, err := service.Exec(context.TODO(), tt.data, tt.tier, tt.now)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, remainingEdits)

			countNoteRepository.AssertExpectations(t)
			latestNoteRepository.AssertExpectations(t)
			createNoteRepository.AssertExpectations(t)
		})
	}
}
