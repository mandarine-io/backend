package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/converter"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type svc struct {
	profileRepo repo.MasterProfileRepository
	serviceRepo repo.MasterServiceRepository
	logger      zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(
	profileRepo repo.MasterProfileRepository,
	serviceRepo repo.MasterServiceRepository,
	opts ...Option,
) domain.MasterServiceService {
	s := &svc{
		profileRepo: profileRepo,
		serviceRepo: serviceRepo,
		logger:      zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *svc) CreateMasterService(
	ctx context.Context,
	userID uuid.UUID,
	input v0.CreateMasterServiceInput,
) (v0.MasterServiceOutput, error) {
	log.Info().Msgf("create master service for master profile %s", userID)

	// Check if master profile exists
	masterProfile, err := s.profileRepo.FindMasterProfileByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to check if master profile exists")
		return v0.MasterServiceOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileDisabled).Msg("master profile is disabled")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileDisabled
	}

	// Save master service
	masterServiceEntity := converter.MapCreateMasterServiceInputToEntity(input)
	masterServiceEntity.MasterProfileID = masterProfile.UserID

	masterServiceEntity, err = s.serviceRepo.CreateMasterService(ctx, masterServiceEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to save master service to db")
		return v0.MasterServiceOutput{}, err
	}

	return converter.MapEntityToMasterServiceOutput(masterServiceEntity), nil
}

func (s *svc) UpdateMasterService(
	ctx context.Context,
	userID uuid.UUID,
	masterServiceID uuid.UUID,
	input v0.UpdateMasterServiceInput,
) (v0.MasterServiceOutput, error) {
	log.Info().Msgf("update master service %s for master profile %s", masterServiceID.String(), userID.String())

	// Check if master profile exists
	masterProfile, err := s.profileRepo.FindMasterProfileByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to check if master profile exists")
		return v0.MasterServiceOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileDisabled).Msg("master profile is disabled")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileDisabled
	}

	// Find master service
	serviceEntity, err := s.serviceRepo.FindMasterServiceByID(ctx, userID, masterServiceID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get master service")
		return v0.MasterServiceOutput{}, err
	}
	if serviceEntity == nil {
		log.Error().Stack().Err(domain.ErrMasterServiceNotExist).Msg("master service not found")
		return v0.MasterServiceOutput{}, domain.ErrMasterServiceNotExist
	}

	// Update master service
	serviceEntity = converter.MapUpdateMasterServiceInputToEntity(serviceEntity, input)
	serviceEntity, err = s.serviceRepo.UpdateMasterService(ctx, serviceEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update master service")
		return v0.MasterServiceOutput{}, err
	}

	return converter.MapEntityToMasterServiceOutput(serviceEntity), nil
}

func (s *svc) DeleteMasterService(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID) error {
	log.Info().Msgf("update master service %s for master profile %s", masterServiceID.String(), userID.String())

	// Check if master profile exists
	masterProfile, err := s.profileRepo.FindMasterProfileByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to check if master profile exists")
		return err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileDisabled).Msg("master profile is disabled")
		return domain.ErrMasterProfileDisabled
	}

	// Delete master service
	err = s.serviceRepo.DeleteMasterServiceByID(ctx, masterProfile.UserID, masterServiceID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update master service")
		return err
	}

	return nil
}

func (s *svc) FindAllMasterServices(
	ctx context.Context,
	input v0.FindMasterServicesInput,
) (v0.MasterServicesOutput, error) {
	log.Info().Msg("find all master services")

	// Generate scopes
	scopes, err := s.generateScopes(input)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to generate DB scopes")
		return v0.MasterServicesOutput{}, err
	}

	// Find all master services
	var (
		data  []*entity.MasterService
		count int64
	)

	group, _ := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			data, err = s.serviceRepo.FindMasterServices(ctx, scopes...)
			return err
		},
	)
	group.Go(
		func() error {
			count, err = s.serviceRepo.CountMasterServices(ctx, scopes...)
			return err
		},
	)

	if err := group.Wait(); err != nil {
		return v0.MasterServicesOutput{}, err
	}

	return converter.MapEntitiesToMasterServicesOutput(data, int(count)), nil
}

func (s *svc) FindAllMasterServicesByUsername(
	ctx context.Context,
	username string,
	input v0.FindMasterServicesInput,
) (v0.MasterServicesOutput, error) {
	log.Info().Msg("find all master services by username")

	// Find master profile
	masterProfile, err := s.profileRepo.FindMasterProfileByUsername(ctx, username)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master profile")
		return v0.MasterServicesOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServicesOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile is disabled")
		return v0.MasterServicesOutput{}, domain.ErrMasterProfileNotExist
	}

	// Generate scopes
	scopes, err := s.generateScopes(input)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to generate DB scopes")
		return v0.MasterServicesOutput{}, err
	}

	// Find all master services
	var (
		data  []*entity.MasterService
		count int64
	)

	group, _ := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			data, err = s.serviceRepo.FindMasterServicesByMasterProfileID(ctx, masterProfile.UserID, scopes...)
			return err
		},
	)
	group.Go(
		func() error {
			count, err = s.serviceRepo.CountMasterServicesByMasterProfileID(ctx, masterProfile.UserID, scopes...)
			return err
		},
	)

	if err := group.Wait(); err != nil {
		return v0.MasterServicesOutput{}, err
	}

	return converter.MapEntitiesToMasterServicesOutput(data, int(count)), nil
}

func (s *svc) FindAllOwnMasterServices(
	ctx context.Context,
	userID uuid.UUID,
	input v0.FindMasterServicesInput,
) (v0.MasterServicesOutput, error) {
	log.Info().Msg("find all own master services")

	// Find master profile
	masterProfile, err := s.profileRepo.FindMasterProfileByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master profile")
		return v0.MasterServicesOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServicesOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileDisabled).Msg("master profile is disabled")
		return v0.MasterServicesOutput{}, domain.ErrMasterProfileDisabled
	}

	// Generate scopes
	scopes, err := s.generateScopes(input)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to generate DB scopes")
		return v0.MasterServicesOutput{}, err
	}

	// Find all master services
	var (
		data  []*entity.MasterService
		count int64
	)

	group, _ := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			data, err = s.serviceRepo.FindMasterServicesByMasterProfileID(ctx, masterProfile.UserID, scopes...)
			return err
		},
	)
	group.Go(
		func() error {
			count, err = s.serviceRepo.CountMasterServicesByMasterProfileID(ctx, masterProfile.UserID, scopes...)
			return err
		},
	)

	if err := group.Wait(); err != nil {
		return v0.MasterServicesOutput{}, err
	}

	return converter.MapEntitiesToMasterServicesOutput(data, int(count)), nil
}

func (s *svc) GetMasterServiceByID(ctx context.Context, username string, id uuid.UUID) (
	v0.MasterServiceOutput,
	error,
) {
	log.Info().Msg("get master service by id")

	// Find master profile
	masterProfile, err := s.profileRepo.FindMasterProfileByUsername(ctx, username)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master profile")
		return v0.MasterServiceOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile is disabled")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileNotExist
	}

	// Find master service
	masterService, err := s.serviceRepo.FindMasterServiceByID(ctx, masterProfile.UserID, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master service")
		return v0.MasterServiceOutput{}, err
	}
	if masterService == nil {
		log.Error().Stack().Err(domain.ErrMasterServiceNotExist).Msg("master service not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterServiceNotExist
	}

	return converter.MapEntityToMasterServiceOutput(masterService), nil
}

func (s *svc) GetOwnMasterServiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (
	v0.MasterServiceOutput,
	error,
) {
	log.Info().Msg("get own master service")

	// Find master profile
	masterProfile, err := s.profileRepo.FindMasterProfileByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master profile")
		return v0.MasterServiceOutput{}, err
	}
	if masterProfile == nil {
		log.Error().Stack().Err(domain.ErrMasterProfileNotExist).Msg("master profile not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileNotExist
	}

	// Check if master profile is disabled
	if !masterProfile.IsEnabled {
		log.Error().Stack().Err(domain.ErrMasterProfileDisabled).Msg("master profile is disabled")
		return v0.MasterServiceOutput{}, domain.ErrMasterProfileDisabled
	}

	// Find master service
	masterService, err := s.serviceRepo.FindMasterServiceByID(ctx, masterProfile.UserID, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to find master service")
		return v0.MasterServiceOutput{}, err
	}
	if masterService == nil {
		log.Error().Stack().Err(domain.ErrMasterServiceNotExist).Msg("master service not exists")
		return v0.MasterServiceOutput{}, domain.ErrMasterServiceNotExist
	}

	return converter.MapEntityToMasterServiceOutput(masterService), nil
}

func (s *svc) generateScopes(input v0.FindMasterServicesInput) ([]repo.Scope, error) {
	var (
		scopes     []repo.Scope
		filter     = input.FindMasterServicesFilterInput
		pagination = input.PaginationInput
		sort       = input.SortInput
	)

	// Generate filter
	if filter.Name != nil {
		scopes = append(scopes, s.serviceRepo.WithNameFilter(*filter.Name))
	}
	if filter.MinPrice != nil {
		scopes = append(scopes, s.serviceRepo.WithMinPriceFilter(*filter.MinPrice))
	}
	if filter.MaxPrice != nil {
		scopes = append(scopes, s.serviceRepo.WithMaxPriceFilter(*filter.MaxPrice))
	}
	if filter.MinInterval != nil {
		minInterval, err := time.ParseDuration(*filter.MinInterval)
		if err == nil {
			scopes = append(scopes, s.serviceRepo.WithMinIntervalFilter(minInterval))
		}
	}
	if filter.MaxInterval != nil {
		maxInterval, err := time.ParseDuration(*filter.MaxInterval)
		if err == nil {
			scopes = append(scopes, s.serviceRepo.WithMaxIntervalFilter(maxInterval))
		}
	}

	// Generate pagination
	if pagination == nil {
		pagination = &v0.PaginationInput{
			Page: 1, PageSize: 10,
		}
	}
	scopes = append(scopes, s.serviceRepo.WithPagination(pagination.Page, pagination.PageSize))

	// Generate sort
	if sort != nil {
		asc := sort.Order == "" || strings.ToLower(sort.Order) == "asc"

		switch sort.Field {
		case "name", "min_price", "max_price", "min_interval", "max_interval":
			scopes = append(scopes, s.serviceRepo.WithSort(sort.Field, asc))
		default:
			return []repo.Scope{}, domain.ErrUnavailableSortField
		}
	}

	return scopes, nil
}
