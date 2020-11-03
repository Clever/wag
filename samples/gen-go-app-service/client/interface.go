package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the app-service service.
type Client interface {

	// HealthCheck makes a GET request to /_health
	// Checks if the service is healthy
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error

	// GetAdmins makes a GET request to /v1/admins
	//
	// 200: []models.Admin
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdmins(ctx context.Context, i *models.GetAdminsInput) ([]models.Admin, error)

	// DeleteAdmin makes a DELETE request to /v1/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteAdmin(ctx context.Context, adminID string) error

	// GetAdminByID makes a GET request to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdminByID(ctx context.Context, adminID string) (*models.Admin, error)

	// UpdateAdmin makes a PATCH request to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateAdmin(ctx context.Context, i *models.UpdateAdminInput) (*models.Admin, error)

	// CreateAdmin makes a PUT request to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateAdmin(ctx context.Context, i *models.CreateAdminInput) (*models.Admin, error)

	// VerifyCode makes a POST request to /v1/admins/{adminID}/confirmation_code
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyCode(ctx context.Context, i *models.VerifyCodeInput) error

	// CreateVerificationCode makes a PUT request to /v1/admins/{adminID}/confirmation_code
	//
	// 200: *models.VerificationCodeResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateVerificationCode(ctx context.Context, i *models.CreateVerificationCodeInput) (*models.VerificationCodeResponse, error)

	// VerifyAdminEmail makes a POST request to /v1/admins/{adminID}/verify_email
	// set the verified email of an admin
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyAdminEmail(ctx context.Context, i *models.VerifyAdminEmailInput) error

	// GetAllAnalyticsApps makes a GET request to /v1/analytics/apps
	//
	// 200: *models.AnalyticsApps
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllAnalyticsApps(ctx context.Context) (*models.AnalyticsApps, error)

	// GetAnalyticsAppByShortname makes a GET request to /v1/analytics/apps/{shortname}
	//
	// 200: *models.AnalyticsApp
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAnalyticsAppByShortname(ctx context.Context, shortname string) (*models.AnalyticsApp, error)

	// GetAllTrackableApps makes a GET request to /v1/analytics/trackable_apps
	//
	// 200: *models.TrackableApps
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllTrackableApps(ctx context.Context) (*models.TrackableApps, error)

	// GetAnalyticsUsageUrls makes a GET request to /v1/analytics/usageUrls
	//
	// 200: *models.UsageUrls
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAnalyticsUsageUrls(ctx context.Context) (*models.UsageUrls, error)

	// GetAllUsageUrls makes a GET request to /v1/appUniverse/usageUrls
	//
	// 200: *models.UsageUrls
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllUsageUrls(ctx context.Context) (*models.UsageUrls, error)

	// GetApps makes a GET request to /v1/apps
	// The server takes in the intersection of input parameters
	// 200: []models.App
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetApps(ctx context.Context, i *models.GetAppsInput) ([]models.App, error)

	// DeleteApp makes a DELETE request to /v1/apps/{appID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteApp(ctx context.Context, appID string) error

	// GetAppByID makes a GET request to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppByID(ctx context.Context, appID string) (*models.App, error)

	// UpdateApp makes a PATCH request to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateApp(ctx context.Context, i *models.UpdateAppInput) (*models.App, error)

	// CreateApp makes a PUT request to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateApp(ctx context.Context, i *models.CreateAppInput) (*models.App, error)

	// GetAdminsForApp makes a GET request to /v1/apps/{appID}/admins
	//
	// 200: []models.AppAdminResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdminsForApp(ctx context.Context, appID string) ([]models.AppAdminResponse, error)

	// UnlinkAppAdmin makes a DELETE request to /v1/apps/{appID}/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UnlinkAppAdmin(ctx context.Context, i *models.UnlinkAppAdminInput) error

	// LinkAppAdmin makes a PUT request to /v1/apps/{appID}/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	LinkAppAdmin(ctx context.Context, i *models.LinkAppAdminInput) error

	// GetGuideConfig makes a GET request to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
	//
	// 200: *models.GuideConfig
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetGuideConfig(ctx context.Context, i *models.GetGuideConfigInput) (*models.GuideConfig, error)

	// SetGuideConfig makes a PUT request to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
	//
	// 200: *models.GuideConfig
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetGuideConfig(ctx context.Context, i *models.SetGuideConfigInput) (*models.GuideConfig, error)

	// GetPermissionsForAdmin makes a GET request to /v1/apps/{appID}/admins/{adminID}/permissions
	//
	// 200: *models.PermissionList
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPermissionsForAdmin(ctx context.Context, i *models.GetPermissionsForAdminInput) (*models.PermissionList, error)

	// VerifyAppAdmin makes a POST request to /v1/apps/{appID}/admins/{adminID}/verify
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyAppAdmin(ctx context.Context, i *models.VerifyAppAdminInput) error

	// GenerateNewBusinessToken makes a POST request to /v1/apps/{appID}/business_token
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GenerateNewBusinessToken(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GetCertifications makes a GET request to /v1/apps/{appID}/certifications/{schoolYearStart}
	//
	// 200: *models.Certifications
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetCertifications(ctx context.Context, i *models.GetCertificationsInput) (*models.Certifications, error)

	// SetCertifications makes a POST request to /v1/apps/{appID}/certifications/{schoolYearStart}
	//
	// 200: *models.Certifications
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetCertifications(ctx context.Context, i *models.SetCertificationsInput) (*models.Certifications, error)

	// GetSetupStep makes a GET request to /v1/apps/{appID}/customStep
	//
	// 200: *models.SetupStep
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSetupStep(ctx context.Context, appID string) (*models.SetupStep, error)

	// CreateSetupStep makes a PATCH request to /v1/apps/{appID}/customStep
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateSetupStep(ctx context.Context, i *models.CreateSetupStepInput) error

	// GetDataRules makes a GET request to /v1/apps/{appID}/data_rules
	//
	// 200: []models.DataRule
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetDataRules(ctx context.Context, appID string) ([]models.DataRule, error)

	// SetDataRules makes a PUT request to /v1/apps/{appID}/data_rules
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetDataRules(ctx context.Context, i *models.SetDataRulesInput) error

	// GetManagers makes a GET request to /v1/apps/{appID}/managers
	//
	// 200: *models.Managers
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetManagers(ctx context.Context, appID string) (*models.Managers, error)

	// GetOnboarding makes a GET request to /v1/apps/{appID}/onboarding
	//
	// 200: *models.Onboarding
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetOnboarding(ctx context.Context, appID string) (*models.Onboarding, error)

	// UpdateOnboarding makes a PATCH request to /v1/apps/{appID}/onboarding
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateOnboarding(ctx context.Context, i *models.UpdateOnboardingInput) error

	// InitializeOnboarding makes a PUT request to /v1/apps/{appID}/onboarding
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	InitializeOnboarding(ctx context.Context, appID string) error

	// DeletePlatform makes a DELETE request to /v1/apps/{appID}/platform/{clientID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeletePlatform(ctx context.Context, i *models.DeletePlatformInput) error

	// UpdatePlatform makes a PATCH request to /v1/apps/{appID}/platform/{clientID}
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdatePlatform(ctx context.Context, i *models.UpdatePlatformInput) (*models.Platform, error)

	// GetPlatformsByAppID makes a GET request to /v1/apps/{appID}/platforms
	//
	// 200: []models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatformsByAppID(ctx context.Context, appID string) ([]models.Platform, error)

	// CreatePlatform makes a PUT request to /v1/apps/{appID}/platforms
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreatePlatform(ctx context.Context, i *models.CreatePlatformInput) (*models.Platform, error)

	// DeleteAppSchema makes a DELETE request to /v1/apps/{appID}/schema
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteAppSchema(ctx context.Context, i *models.DeleteAppSchemaInput) error

	// GetAppSchema makes a GET request to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppSchema(ctx context.Context, appID string) (*models.AppSchema, error)

	// CreateAppSchema makes a POST request to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateAppSchema(ctx context.Context, i *models.CreateAppSchemaInput) (*models.AppSchema, error)

	// SetAppSchema makes a PUT request to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetAppSchema(ctx context.Context, i *models.SetAppSchemaInput) (*models.AppSchema, error)

	// GetSecrets makes a GET request to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSecrets(ctx context.Context, appID string) (*models.SecretConfig, error)

	// RevokeOldClientSecret makes a PATCH request to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	RevokeOldClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GenerateNewClientSecret makes a POST request to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GenerateNewClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// ResetClientSecret makes a PUT request to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	ResetClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GetRecommendedSharing makes a GET request to /v1/apps/{appID}/sharing
	//
	// 200: *models.SharingRecommendations
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetRecommendedSharing(ctx context.Context, appID string) (*models.SharingRecommendations, error)

	// SetRecommendedSharing makes a PUT request to /v1/apps/{appID}/sharing
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetRecommendedSharing(ctx context.Context, i *models.SetRecommendedSharingInput) error

	// UpdateAppIcon makes a POST request to /v1/apps/{appID}/update_icon
	//
	// 200: *models.Image
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 422: *models.UnprocessableEntity
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateAppIcon(ctx context.Context, i *models.UpdateAppIconInput) (*models.Image, error)

	// GetAllCategories makes a GET request to /v1/categories
	//
	// 200: *models.Categories
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllCategories(ctx context.Context) (*models.Categories, error)

	// GetKnownHosts makes a GET request to /v1/knownhosts
	//
	// 200: []models.KnownHost
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetKnownHosts(ctx context.Context) ([]models.KnownHost, error)

	// GetAllLibraryResources makes a GET request to /v1/libraryResources
	//
	// 200: *models.LibraryResources
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllLibraryResources(ctx context.Context, i *models.GetAllLibraryResourcesInput) (*models.LibraryResources, error)

	// SearchLibraryResource makes a GET request to /v1/libraryResources/search
	//
	// 200: *models.LibraryResources
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SearchLibraryResource(ctx context.Context, i *models.SearchLibraryResourceInput) (*models.LibraryResources, error)

	// GetLibraryResourceByShortname makes a GET request to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetLibraryResourceByShortname(ctx context.Context, i *models.GetLibraryResourceByShortnameInput) (*models.LibraryResource, error)

	// UpdateLibraryResourceByShortname makes a PATCH request to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateLibraryResourceByShortname(ctx context.Context, i *models.UpdateLibraryResourceByShortnameInput) (*models.LibraryResource, error)

	// CreateLibraryResource makes a POST request to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateLibraryResource(ctx context.Context, i *models.CreateLibraryResourceInput) (*models.LibraryResource, error)

	// DeleteLibraryResourceLink makes a DELETE request to /v1/libraryResources/{shortname}/link
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteLibraryResourceLink(ctx context.Context, shortname string) error

	// GetValidPermissions makes a GET request to /v1/permissions
	//
	// 200: *models.GetValidPermissionsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetValidPermissions(ctx context.Context) (*models.GetValidPermissionsResponse, error)

	// GetPlatforms makes a GET request to /v1/platforms
	// The server takes in the intersection of input parameters
	// 200: []models.Platform
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatforms(ctx context.Context, i *models.GetPlatformsInput) ([]models.Platform, error)

	// GetPlatformByClientID makes a GET request to /v1/platforms/{clientID}
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatformByClientID(ctx context.Context, clientID string) (*models.Platform, error)

	// GetAppsForAdmin makes a GET request to /v2/admins/{adminID}/apps
	//
	// 200: []models.AppForAdminResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppsForAdmin(ctx context.Context, adminID string) ([]models.AppForAdminResponse, error)

	// OverrideConfig makes a POST request to /v2/apps/{srcAppID}/override-config/{destAppID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	OverrideConfig(ctx context.Context, i *models.OverrideConfigInput) error
}
