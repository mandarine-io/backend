package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	httpdto "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

// Account error
var (
	ErrInvalidOrExpiredOtp  = httpdto.NewI18nError("invalid or expired otp", "errors.invalid_or_expired_otp")
	ErrUserNotFound         = httpdto.NewI18nError("user not found", "errors.user_not_found")
	ErrDuplicateUsername    = httpdto.NewI18nError("username already in use", "errors.duplicate_username")
	ErrDuplicateEmail       = httpdto.NewI18nError("email already in use", "errors.duplicate_email")
	ErrPasswordIsSet        = httpdto.NewI18nError("password is already set", "errors.password_is_set")
	ErrIncorrectOldPassword = httpdto.NewI18nError("incorrect old password", "errors.incorrect_old_password")
	ErrUserNotDeleted       = httpdto.NewI18nError("user not deleted", "errors.user_not_deleted")
	ErrUserAlreadyDeleted   = httpdto.NewI18nError("user already deleted", "errors.user_already_deleted")
	ErrSendEmail            = httpdto.NewI18nError("failed to send email", "errors.failed_to_send_email")
)

// Auth error
var (
	ErrDuplicateUser       = httpdto.NewI18nError("duplicate user", "errors.duplicate_user")
	ErrBadCredentials      = httpdto.NewI18nError("bad credentials", "errors.bad_credentials")
	ErrUserIsBlocked       = httpdto.NewI18nError("user is blocked", "errors.user_is_blocked")
	ErrInvalidJwtToken     = httpdto.NewI18nError("invalid JWT token", "errors.invalid_jwt_token")
	ErrUserInfoNotReceived = httpdto.NewI18nError("user info not received", "errors.userinfo_not_received")
	ErrInvalidProvider     = httpdto.NewI18nError("invalid provider", "errors.invalid_provider")
)

// Geocoding error
var (
	ErrGeocodeProvidersUnavailable = httpdto.NewI18nError("geocode providers unavailable", "errors.geocode_providers_unavailable")
)

// Master profile error
var (
	ErrDuplicateMasterProfile = httpdto.NewI18nError("duplicate master profile", "errors.duplicate_master_profile")
	ErrMasterProfileNotExist  = httpdto.NewI18nError("master profile not exist", "errors.master_profile_not_exist")
	ErrMasterProfileNotFound  = httpdto.NewI18nError("master profile not found", "errors.master_profile_not_found")
	ErrUnavailableSortField   = httpdto.NewI18nError("unavailable sort field", "errors.unavailable_sort_field")
)

// Resource error
var (
	ErrResourceNotUploaded = httpdto.NewI18nError("resource not uploaded", "errors.resource_not_uploaded")
)

type AccountService interface {
	GetAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error)
	UpdateUsername(ctx context.Context, id uuid.UUID, input dto.UpdateUsernameInput) (dto.AccountOutput, error)
	UpdateEmail(ctx context.Context, id uuid.UUID, input dto.UpdateEmailInput, localizer *i18n.Localizer) (dto.AccountOutput, error)
	VerifyEmail(ctx context.Context, id uuid.UUID, req dto.VerifyEmailInput) error
	SetPassword(ctx context.Context, id uuid.UUID, input dto.SetPasswordInput) error
	UpdatePassword(ctx context.Context, id uuid.UUID, input dto.UpdatePasswordInput) error
	RestoreAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterInput, localizer *i18n.Localizer) error
	RegisterConfirm(ctx context.Context, input dto.RegisterConfirmInput) error
	Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error)
	RefreshTokens(ctx context.Context, refreshToken string) (dto.JwtTokensOutput, error)
	Logout(ctx context.Context, jti string) error
	RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput, localizer *i18n.Localizer) error
	VerifyRecoveryCode(ctx context.Context, input dto.VerifyRecoveryCodeInput) error
	ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error
	GetConsentPageUrl(_ context.Context, provider string, redirectUrl string) (dto.GetConsentPageUrlOutput, error)
	FetchUserInfo(ctx context.Context, provider string, input dto.FetchUserInfoInput) (oauth.UserInfo, error)
	RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (dto.JwtTokensOutput, error)
}

type GeocodingService interface {
	Geocode(ctx context.Context, input dto.GeocodingInput, lang language.Tag) (dto.GeocodingOutput, error)
	ReverseGeocode(ctx context.Context, input dto.ReverseGeocodingInput, lang language.Tag) (dto.ReverseGeocodingOutput, error)
}

type HealthService interface {
	Health() []dto.HealthOutput
}

type MasterProfileService interface {
	CreateMasterProfile(ctx context.Context, id uuid.UUID, input dto.CreateMasterProfileInput) (dto.OwnMasterProfileOutput, error)
	FindMasterProfiles(ctx context.Context, input dto.FindMasterProfilesInput) (dto.MasterProfilesOutput, error)
	GetOwnMasterProfile(ctx context.Context, id uuid.UUID) (dto.OwnMasterProfileOutput, error)
	GetMasterProfileByUsername(ctx context.Context, username string) (dto.MasterProfileOutput, error)
	UpdateMasterProfile(ctx context.Context, id uuid.UUID, input dto.UpdateMasterProfileInput) (dto.OwnMasterProfileOutput, error)
}

type ResourceService interface {
	UploadResource(ctx context.Context, input *dto.UploadResourceInput) (dto.UploadResourceOutput, error)
	UploadResources(ctx context.Context, input *dto.UploadResourcesInput) (dto.UploadResourcesOutput, error)
	DownloadResource(ctx context.Context, objectID string) (*s3.FileData, error)
}

type WebsocketService interface {
	RegisterClient(userId uuid.UUID, r *http.Request, w http.ResponseWriter) error
}
