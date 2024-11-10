package profile

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	commonmapper "github.com/mandarine-io/Backend/internal/domain/mapper"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/master/profile/mapper"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var (
	availableSortFields = []string{"display_name", "job", "address", "point"}
)

type svc struct {
	repo repo.MasterProfileRepository
}

func NewService(repo repo.MasterProfileRepository) service.MasterProfileService {
	return &svc{repo: repo}
}

func (s *svc) CreateMasterProfile(ctx context.Context, id uuid.UUID, input dto.CreateMasterProfileInput) (dto.OwnMasterProfileOutput, error) {
	log.Info().Msgf("create master profile for user %s", id.String())

	// Check if user exists
	exists, err := s.repo.ExistsMasterProfileByUserId(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to check if master profile exists")
		return dto.OwnMasterProfileOutput{}, err
	}
	if exists {
		log.Error().Stack().Err(service.ErrDuplicateMasterProfile).Msg("master profile already exists")
		return dto.OwnMasterProfileOutput{}, service.ErrDuplicateMasterProfile
	}

	// Create master profile
	entity := mapper.MapCreateMasterProfileInputToEntity(id, input)
	entity, err = s.repo.CreateMasterProfile(ctx, entity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to create master profile")

		if errors.Is(err, repo.ErrDuplicateMasterProfile) {
			return dto.OwnMasterProfileOutput{}, service.ErrDuplicateMasterProfile
		}
		if errors.Is(err, repo.ErrUserForMasterProfileNotExist) {
			return dto.OwnMasterProfileOutput{}, service.ErrUserNotFound
		}
		return dto.OwnMasterProfileOutput{}, err
	}

	output := mapper.MapEntityToOwnMasterProfileOutput(entity)
	return output, nil
}

func (s *svc) FindMasterProfiles(ctx context.Context, input dto.FindMasterProfilesInput) (dto.MasterProfilesOutput, error) {
	log.Info().Msg("find master profiles")

	// Generate filters, pagination and sorts
	filterMap := mapper.MapFindMasterProfilesToFilterMap(input.FindMasterProfilesFilterInput)
	dbPagination := commonmapper.MapDomainPaginationToPagination(input.PaginationInput)
	dbSorts, err := commonmapper.MapDomainSortsToSortsWithAvailables([]*dto.SortInput{input.SortInput}, availableSortFields)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to map domain sorts to db sorts")
		return dto.MasterProfilesOutput{}, service.ErrUnavailableSortField
	}

	var (
		masterProfiles []*model.MasterProfileEntity
		count          int64
	)
	executor, _ := errgroup.WithContext(ctx)
	executor.Go(func() error {
		var err error
		masterProfiles, err = s.repo.FindMasterProfiles(ctx, filterMap, dbPagination, dbSorts)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to find master profiles")
		}
		return err
	})
	executor.Go(func() error {
		var err error
		count, err = s.repo.CountMasterProfiles(ctx, filterMap)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to count master profiles")
		}
		return err
	})

	err = executor.Wait()
	if err != nil {
		return dto.MasterProfilesOutput{}, err
	}

	output := mapper.MapEntitiesToMasterProfilesOutput(masterProfiles, int(count))
	return output, nil
}

func (s *svc) GetOwnMasterProfile(ctx context.Context, id uuid.UUID) (dto.OwnMasterProfileOutput, error) {
	log.Info().Msgf("get own master profile for user %s", id.String())

	entity, err := s.repo.FindMasterProfileByUserId(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get own master profile")
		return dto.OwnMasterProfileOutput{}, err
	}
	if entity == nil {
		log.Error().Stack().Err(service.ErrMasterProfileNotExist).Msg("master profile not found")
		return dto.OwnMasterProfileOutput{}, service.ErrMasterProfileNotExist
	}

	output := mapper.MapEntityToOwnMasterProfileOutput(entity)
	return output, nil
}

func (s *svc) GetMasterProfileByUsername(ctx context.Context, username string) (dto.MasterProfileOutput, error) {
	log.Info().Msgf("get master profile by username %s", username)

	entity, err := s.repo.FindEnabledMasterProfileByUsername(ctx, username)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get master profile")
		return dto.MasterProfileOutput{}, err
	}
	if entity == nil {
		log.Error().Stack().Err(service.ErrMasterProfileNotFound).Msg("master profile not found")
		return dto.MasterProfileOutput{}, service.ErrMasterProfileNotFound
	}

	output := mapper.MapEntityToMasterProfileOutput(entity)
	return output, nil
}

func (s *svc) UpdateMasterProfile(ctx context.Context, id uuid.UUID, input dto.UpdateMasterProfileInput) (dto.OwnMasterProfileOutput, error) {
	log.Info().Msgf("update master profile for user %s", id.String())

	entity, err := s.repo.FindMasterProfileByUserId(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get master profile")
		return dto.OwnMasterProfileOutput{}, err
	}
	if entity == nil {
		log.Error().Stack().Err(service.ErrMasterProfileNotExist).Msg("master profile not found")
		return dto.OwnMasterProfileOutput{}, service.ErrMasterProfileNotExist
	}

	entity = mapper.MapUpdateMasterProfileInputToEntity(entity, input)
	entity, err = s.repo.UpdateMasterProfile(ctx, entity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update master profile")
		return dto.OwnMasterProfileOutput{}, err
	}

	output := mapper.MapEntityToOwnMasterProfileOutput(entity)
	return output, nil
}
