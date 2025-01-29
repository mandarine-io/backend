package profile

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/converter"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"strings"
)

type svc struct {
	repo   repo.MasterProfileRepository
	logger zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(repo repo.MasterProfileRepository, opts ...Option) domain.MasterProfileService {
	s := &svc{repo: repo}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *svc) CreateMasterProfile(
	ctx context.Context,
	id uuid.UUID,
	input v0.CreateMasterProfileInput,
) (v0.MasterProfileOutput, error) {
	s.logger.Info().Msgf("create master profile for user %s", id.String())

	// Check if user exists
	exists, err := s.repo.ExistsMasterProfileByUserID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check if master profile exists")
		return v0.MasterProfileOutput{}, err
	}
	if exists {
		s.logger.Error().Stack().Err(domain.ErrDuplicateMasterProfile).Msg("master profile already exists")
		return v0.MasterProfileOutput{}, domain.ErrDuplicateMasterProfile
	}

	// Create master profile
	profileEntity := converter.MapCreateMasterProfileInputToEntity(id, input)
	profileEntity, err = s.repo.CreateMasterProfile(ctx, profileEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to create master profile")

		if errors.Is(err, repo.ErrDuplicateMasterProfile) {
			return v0.MasterProfileOutput{}, domain.ErrDuplicateMasterProfile
		}
		if errors.Is(err, repo.ErrUserForMasterProfileNotExist) {
			return v0.MasterProfileOutput{}, domain.ErrUserNotFound
		}
		return v0.MasterProfileOutput{}, err
	}

	output := converter.MapEntityToMasterProfileOutput(profileEntity)
	return output, nil
}

func (s *svc) FindMasterProfiles(ctx context.Context, input v0.FindMasterProfilesInput) (
	v0.MasterProfilesOutput,
	error,
) {
	s.logger.Info().Msg("find master profiles")

	scopes, err := s.generateScopes(input)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to generate DB scopes")
		return v0.MasterProfilesOutput{}, err
	}

	masterProfiles, count, err := s.findAndCountMasterProfiles(ctx, scopes)
	if err != nil {
		return v0.MasterProfilesOutput{}, err
	}

	output := converter.MapEntitiesToMasterProfilesOutput(masterProfiles, int(count))
	return output, nil
}

func (s *svc) findAndCountMasterProfiles(ctx context.Context, scopes []repo.Scope) (
	[]*entity.MasterProfile,
	int64,
	error,
) {
	var (
		masterProfiles []*entity.MasterProfile
		count          int64
	)
	executor, _ := errgroup.WithContext(ctx)
	executor.Go(
		func() error {
			var err error
			masterProfiles, err = s.repo.FindMasterProfiles(ctx, scopes...)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("failed to find master profiles")
			}
			return err
		},
	)
	executor.Go(
		func() error {
			var err error
			count, err = s.repo.CountMasterProfiles(ctx, scopes...)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("failed to count master profiles")
			}
			return err
		},
	)

	err := executor.Wait()
	return masterProfiles, count, err
}

func (s *svc) GetOwnMasterProfile(ctx context.Context, id uuid.UUID) (v0.MasterProfileOutput, error) {
	s.logger.Info().Msgf("get own master profile for user %s", id.String())

	profileEntity, err := s.repo.FindMasterProfileByUserID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get own master profile")
		return v0.MasterProfileOutput{}, err
	}
	if profileEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not found")
		return v0.MasterProfileOutput{}, domain.ErrMasterProfileNotExist
	}

	if !profileEntity.IsEnabled {
		return v0.MasterProfileOutput{}, domain.ErrMasterProfileDisabled
	}

	output := converter.MapEntityToMasterProfileOutput(profileEntity)
	return output, nil
}

func (s *svc) GetMasterProfileByUsername(ctx context.Context, username string) (v0.MasterProfileOutput, error) {
	s.logger.Info().Msgf("get master profile by username %s", username)

	profileEntity, err := s.repo.FindEnabledMasterProfileByUsername(ctx, username)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get master profile")
		return v0.MasterProfileOutput{}, err
	}
	if profileEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrMasterProfileNotFound).Msg("master profile not found")
		return v0.MasterProfileOutput{}, domain.ErrMasterProfileNotFound
	}

	output := converter.MapEntityToMasterProfileOutput(profileEntity)
	output.IsEnabled = nil

	return output, nil
}

func (s *svc) UpdateMasterProfile(
	ctx context.Context,
	id uuid.UUID,
	input v0.UpdateMasterProfileInput,
) (v0.MasterProfileOutput, error) {
	s.logger.Info().Msgf("update master profile for user %s", id.String())

	profileEntity, err := s.repo.FindMasterProfileByUserID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get master profile")
		return v0.MasterProfileOutput{}, err
	}
	if profileEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not found")
		return v0.MasterProfileOutput{}, domain.ErrMasterProfileNotExist
	}

	profileEntity = converter.MapUpdateMasterProfileInputToEntity(profileEntity, input)
	profileEntity, err = s.repo.UpdateMasterProfile(ctx, profileEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update master profile")
		return v0.MasterProfileOutput{}, err
	}

	output := converter.MapEntityToMasterProfileOutput(profileEntity)
	return output, nil
}

func (s *svc) generateScopes(input v0.FindMasterProfilesInput) ([]repo.Scope, error) {
	var (
		scopes     []repo.Scope
		filter     = input.FindMasterProfilesFilterInput
		pagination = input.PaginationInput
		sort       = input.SortInput
	)

	// Generate filters
	if filter != nil {
		if filter.DisplayName != nil {
			scopes = append(scopes, s.repo.WithDisplayNameFilter(*filter.DisplayName))
		}

		if filter.Job != nil {
			scopes = append(scopes, s.repo.WithJobFilter(*filter.Job))
		}

		if filter.Radius != nil && filter.Lng != nil && filter.Lat != nil {
			scopes = append(scopes, s.repo.WithPointFilter(*filter.Lat, *filter.Lng, *filter.Radius))
		}
	}

	// Generate pagination
	if pagination == nil {
		pagination = &v0.PaginationInput{
			Page: 1, PageSize: 10,
		}
	}
	scopes = append(scopes, s.repo.WithPagination(pagination.Page, pagination.PageSize))

	// Generate sort
	if sort != nil {
		asc := sort.Order == "" || strings.ToLower(sort.Order) == "asc"

		switch sort.Field {
		case "display_name", "job", "address":
			scopes = append(scopes, s.repo.WithColumnSort(sort.Field, asc))
		case "point":
			if filter == nil || filter.Lng == nil || filter.Lat == nil {
				return []repo.Scope{}, domain.ErrMissingPointForSorting
			}
			scopes = append(scopes, s.repo.WithPointSort(*filter.Lat, *filter.Lng, asc))
		default:
			return []repo.Scope{}, domain.ErrUnavailableSortField
		}
	}

	return scopes, nil
}
