package domain

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/pkg/model/health"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/third_party/oauth"
	"golang.org/x/text/language"
	"net/http"
)

var (
	// Account error

	ErrUserNotFound         = v0.NewI18nError("user not found", "errors.user_not_found")
	ErrDuplicateUsername    = v0.NewI18nError("username already in use", "errors.duplicate_username")
	ErrDuplicateEmail       = v0.NewI18nError("email already in use", "errors.duplicate_email")
	ErrPasswordIsSet        = v0.NewI18nError("password is already set", "errors.password_is_set")
	ErrIncorrectOldPassword = v0.NewI18nError("incorrect old password", "errors.incorrect_old_password")
	ErrUserNotDeleted       = v0.NewI18nError("user not deleted", "errors.user_not_deleted")
	ErrUserAlreadyDeleted   = v0.NewI18nError("user already deleted", "errors.user_already_deleted")
	ErrSendEmail            = v0.NewI18nError("failed to send email", "errors.failed_to_send_email")

	// Auth error

	ErrDuplicateUser       = v0.NewI18nError("duplicate user", "errors.duplicate_user")
	ErrBadCredentials      = v0.NewI18nError("bad credentials", "errors.bad_credentials")
	ErrUserIsBlocked       = v0.NewI18nError("user is blocked", "errors.user_is_blocked")
	ErrUserInfoNotReceived = v0.NewI18nError("user info not received", "errors.userinfo_not_received")
	ErrInvalidProvider     = v0.NewI18nError("invalid provider", "errors.invalid_provider")

	// Geocoding error

	ErrGeocodeProvidersUnavailable = v0.NewI18nError(
		"geocode providers unavailable",
		"errors.geocode_providers_unavailable",
	)

	// Master profile error

	ErrDuplicateMasterProfile = v0.NewI18nError("duplicate master profile", "errors.duplicate_master_profile")
	ErrMasterProfileNotExist  = v0.NewI18nError("master profile not exist", "errors.master_profile_not_exist")
	ErrMasterProfileDisabled  = v0.NewI18nError("master profile disabled", "errors.master_profile_disabled")
	ErrMasterProfileNotFound  = v0.NewI18nError("master profile not found", "errors.master_profile_not_found")
	ErrUnavailableSortField   = v0.NewI18nError("unavailable sort field", "errors.unavailable_sort_field")
	ErrMissingPointForSorting = v0.NewI18nError("missing point for sorting", "errors.missing_point_for_sorting")

	// Master service error

	ErrMasterServiceCreation     = v0.NewI18nError("master service creation", "errors.master_service_creation")
	ErrMasterServiceModification = v0.NewI18nError(
		"master service modification",
		"errors.master_service_modification",
	)
	ErrMasterServiceDeletion = v0.NewI18nError("master service deletion", "errors.master_service_deletion")
	ErrMasterServiceNotExist = v0.NewI18nError("master service not exist", "errors.master_service_not_exist")

	// Resource error

	ErrResourceNotUploaded = v0.NewI18nError("resource not uploaded", "errors.resource_not_uploaded")
)

type AccountService interface {
	GetAccount(ctx context.Context, id uuid.UUID) (v0.AccountOutput, error)
	UpdateUsername(ctx context.Context, id uuid.UUID, input v0.UpdateUsernameInput) (v0.AccountOutput, error)
	UpdateEmail(
		ctx context.Context,
		id uuid.UUID,
		input v0.UpdateEmailInput,
		localizer locale.Localizer,
	) (v0.AccountOutput, error)
	VerifyEmail(ctx context.Context, id uuid.UUID, input v0.VerifyEmailInput) error
	SetPassword(ctx context.Context, id uuid.UUID, input v0.SetPasswordInput) error
	UpdatePassword(ctx context.Context, id uuid.UUID, input v0.UpdatePasswordInput) error
	RestoreAccount(ctx context.Context, id uuid.UUID) (v0.AccountOutput, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type AuthService interface {
	Register(ctx context.Context, input v0.RegisterInput, localizer locale.Localizer) error
	RegisterConfirm(ctx context.Context, input v0.RegisterConfirmInput) error
	Login(ctx context.Context, input v0.LoginInput) (v0.JwtTokensOutput, error)
	RefreshTokens(ctx context.Context, input v0.RefreshTokensInput) (v0.JwtTokensOutput, error)
	Logout(ctx context.Context, jti string) error
	RecoveryPassword(ctx context.Context, input v0.RecoveryPasswordInput, localizer locale.Localizer) error
	VerifyRecoveryCode(ctx context.Context, input v0.VerifyRecoveryCodeInput) error
	ResetPassword(ctx context.Context, input v0.ResetPasswordInput) error
	GetConsentPageURL(_ context.Context, provider string, redirectURL string) (v0.GetConsentPageURLOutput, error)
	FetchUserInfo(ctx context.Context, provider string, input v0.FetchUserInfoInput) (oauth.UserInfo, error)
	RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (v0.JwtTokensOutput, error)
}

type GeocodingService interface {
	Geocode(ctx context.Context, input v0.GeocodingInput, lang language.Tag) (v0.GeocodingOutput, error)
	ReverseGeocode(
		ctx context.Context,
		input v0.ReverseGeocodingInput,
		lang language.Tag,
	) (v0.ReverseGeocodingOutput, error)
}

type HealthService interface {
	Health() []health.HealthOutput
}

type MasterProfileService interface {
	CreateMasterProfile(
		ctx context.Context,
		id uuid.UUID,
		input v0.CreateMasterProfileInput,
	) (v0.MasterProfileOutput, error)
	UpdateMasterProfile(
		ctx context.Context,
		id uuid.UUID,
		input v0.UpdateMasterProfileInput,
	) (v0.MasterProfileOutput, error)
	FindMasterProfiles(ctx context.Context, input v0.FindMasterProfilesInput) (v0.MasterProfilesOutput, error)
	GetOwnMasterProfile(ctx context.Context, id uuid.UUID) (v0.MasterProfileOutput, error)
	GetMasterProfileByUsername(ctx context.Context, username string) (v0.MasterProfileOutput, error)
}

type MasterServiceService interface {
	CreateMasterService(
		ctx context.Context,
		userID uuid.UUID,
		input v0.CreateMasterServiceInput,
	) (v0.MasterServiceOutput, error)
	UpdateMasterService(
		ctx context.Context,
		userID uuid.UUID,
		masterServiceID uuid.UUID,
		input v0.UpdateMasterServiceInput,
	) (v0.MasterServiceOutput, error)
	DeleteMasterService(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID) error
	FindAllMasterServices(ctx context.Context, input v0.FindMasterServicesInput) (v0.MasterServicesOutput, error)
	FindAllOwnMasterServices(
		ctx context.Context,
		userID uuid.UUID,
		input v0.FindMasterServicesInput,
	) (v0.MasterServicesOutput, error)
	FindAllMasterServicesByUsername(
		ctx context.Context,
		username string,
		input v0.FindMasterServicesInput,
	) (v0.MasterServicesOutput, error)
	GetMasterServiceByID(ctx context.Context, username string, id uuid.UUID) (v0.MasterServiceOutput, error)
	GetOwnMasterServiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (v0.MasterServiceOutput, error)
}

type ResourceService interface {
	UploadResource(ctx context.Context, input *v0.UploadResourceInput) (v0.UploadResourceOutput, error)
	UploadResources(ctx context.Context, input *v0.UploadResourcesInput) (v0.UploadResourcesOutput, error)
	DownloadResource(ctx context.Context, objectID string) (*s3.FileData, error)
}

type WebsocketService interface {
	RegisterClient(userID uuid.UUID, r *http.Request, w http.ResponseWriter) error
}
