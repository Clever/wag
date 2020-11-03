package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the app-service service.
type Controller interface {

	// HealthCheck handles GET requests to /_health
	// Checks if the service is healthy
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error

	// GetAdmins handles GET requests to /v1/admins
	//
	// 200: []models.Admin
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdmins(ctx context.Context, i *models.GetAdminsInput) ([]models.Admin, error)

	// DeleteAdmin handles DELETE requests to /v1/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteAdmin(ctx context.Context, adminID string) error

	// GetAdminByID handles GET requests to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdminByID(ctx context.Context, adminID string) (*models.Admin, error)

	// UpdateAdmin handles PATCH requests to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateAdmin(ctx context.Context, i *models.UpdateAdminInput) (*models.Admin, error)

	// CreateAdmin handles PUT requests to /v1/admins/{adminID}
	//
	// 200: *models.Admin
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateAdmin(ctx context.Context, i *models.CreateAdminInput) (*models.Admin, error)

	// GetAppsForAdminDeprecated handles GET requests to /v1/admins/{adminID}/apps
	//
	// 200: []models.App
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppsForAdminDeprecated(ctx context.Context, adminID string) ([]models.App, error)

	// VerifyCode handles POST requests to /v1/admins/{adminID}/confirmation_code
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyCode(ctx context.Context, i *models.VerifyCodeInput) error

	// CreateVerificationCode handles PUT requests to /v1/admins/{adminID}/confirmation_code
	//
	// 200: *models.VerificationCodeResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateVerificationCode(ctx context.Context, i *models.CreateVerificationCodeInput) (*models.VerificationCodeResponse, error)

	// VerifyAdminEmail handles POST requests to /v1/admins/{adminID}/verify_email
	// set the verified email of an admin
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyAdminEmail(ctx context.Context, i *models.VerifyAdminEmailInput) error

	// GetAllAnalyticsApps handles GET requests to /v1/analytics/apps
	//
	// 200: *models.AnalyticsApps
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllAnalyticsApps(ctx context.Context) (*models.AnalyticsApps, error)

	// GetAnalyticsAppByShortname handles GET requests to /v1/analytics/apps/{shortname}
	//
	// 200: *models.AnalyticsApp
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAnalyticsAppByShortname(ctx context.Context, shortname string) (*models.AnalyticsApp, error)

	// GetAllTrackableApps handles GET requests to /v1/analytics/trackable_apps
	//
	// 200: *models.TrackableApps
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllTrackableApps(ctx context.Context) (*models.TrackableApps, error)

	// GetAnalyticsUsageUrls handles GET requests to /v1/analytics/usageUrls
	//
	// 200: *models.UsageUrls
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAnalyticsUsageUrls(ctx context.Context) (*models.UsageUrls, error)

	// GetAllUsageUrls handles GET requests to /v1/appUniverse/usageUrls
	//
	// 200: *models.UsageUrls
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllUsageUrls(ctx context.Context) (*models.UsageUrls, error)

	// GetApps handles GET requests to /v1/apps
	// The server takes in the intersection of input parameters
	// 200: []models.App
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetApps(ctx context.Context, i *models.GetAppsInput) ([]models.App, error)

	// DeleteApp handles DELETE requests to /v1/apps/{appID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteApp(ctx context.Context, appID string) error

	// GetAppByID handles GET requests to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppByID(ctx context.Context, appID string) (*models.App, error)

	// UpdateApp handles PATCH requests to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateApp(ctx context.Context, i *models.UpdateAppInput) (*models.App, error)

	// CreateApp handles PUT requests to /v1/apps/{appID}
	//
	// 200: *models.App
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateApp(ctx context.Context, i *models.CreateAppInput) (*models.App, error)

	// GetAdminsForApp handles GET requests to /v1/apps/{appID}/admins
	//
	// 200: []models.AppAdminResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAdminsForApp(ctx context.Context, appID string) ([]models.AppAdminResponse, error)

	// UnlinkAppAdmin handles DELETE requests to /v1/apps/{appID}/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UnlinkAppAdmin(ctx context.Context, i *models.UnlinkAppAdminInput) error

	// LinkAppAdmin handles PUT requests to /v1/apps/{appID}/admins/{adminID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	LinkAppAdmin(ctx context.Context, i *models.LinkAppAdminInput) error

	// GetGuideConfig handles GET requests to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
	//
	// 200: *models.GuideConfig
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetGuideConfig(ctx context.Context, i *models.GetGuideConfigInput) (*models.GuideConfig, error)

	// SetGuideConfig handles PUT requests to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
	//
	// 200: *models.GuideConfig
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetGuideConfig(ctx context.Context, i *models.SetGuideConfigInput) (*models.GuideConfig, error)

	// GetPermissionsForAdmin handles GET requests to /v1/apps/{appID}/admins/{adminID}/permissions
	//
	// 200: *models.PermissionList
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPermissionsForAdmin(ctx context.Context, i *models.GetPermissionsForAdminInput) (*models.PermissionList, error)

	// VerifyAppAdmin handles POST requests to /v1/apps/{appID}/admins/{adminID}/verify
	//
	// 200: nil
	// 400: *models.BadRequest
	// 403: *models.Forbidden
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	VerifyAppAdmin(ctx context.Context, i *models.VerifyAppAdminInput) error

	// GenerateNewBusinessToken handles POST requests to /v1/apps/{appID}/business_token
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GenerateNewBusinessToken(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GetCertifications handles GET requests to /v1/apps/{appID}/certifications/{schoolYearStart}
	//
	// 200: *models.Certifications
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetCertifications(ctx context.Context, i *models.GetCertificationsInput) (*models.Certifications, error)

	// SetCertifications handles POST requests to /v1/apps/{appID}/certifications/{schoolYearStart}
	//
	// 200: *models.Certifications
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetCertifications(ctx context.Context, i *models.SetCertificationsInput) (*models.Certifications, error)

	// GetSetupStep handles GET requests to /v1/apps/{appID}/customStep
	//
	// 200: *models.SetupStep
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSetupStep(ctx context.Context, appID string) (*models.SetupStep, error)

	// CreateSetupStep handles PATCH requests to /v1/apps/{appID}/customStep
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateSetupStep(ctx context.Context, i *models.CreateSetupStepInput) error

	// GetDataRules handles GET requests to /v1/apps/{appID}/data_rules
	//
	// 200: []models.DataRule
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetDataRules(ctx context.Context, appID string) ([]models.DataRule, error)

	// SetDataRules handles PUT requests to /v1/apps/{appID}/data_rules
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetDataRules(ctx context.Context, i *models.SetDataRulesInput) error

	// GetManagers handles GET requests to /v1/apps/{appID}/managers
	//
	// 200: *models.Managers
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetManagers(ctx context.Context, appID string) (*models.Managers, error)

	// GetOnboarding handles GET requests to /v1/apps/{appID}/onboarding
	//
	// 200: *models.Onboarding
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetOnboarding(ctx context.Context, appID string) (*models.Onboarding, error)

	// UpdateOnboarding handles PATCH requests to /v1/apps/{appID}/onboarding
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateOnboarding(ctx context.Context, i *models.UpdateOnboardingInput) error

	// InitializeOnboarding handles PUT requests to /v1/apps/{appID}/onboarding
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	InitializeOnboarding(ctx context.Context, appID string) error

	// DeletePlatform handles DELETE requests to /v1/apps/{appID}/platform/{clientID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeletePlatform(ctx context.Context, i *models.DeletePlatformInput) error

	// UpdatePlatform handles PATCH requests to /v1/apps/{appID}/platform/{clientID}
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdatePlatform(ctx context.Context, i *models.UpdatePlatformInput) (*models.Platform, error)

	// GetPlatformsByAppID handles GET requests to /v1/apps/{appID}/platforms
	//
	// 200: []models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatformsByAppID(ctx context.Context, appID string) ([]models.Platform, error)

	// CreatePlatform handles PUT requests to /v1/apps/{appID}/platforms
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreatePlatform(ctx context.Context, i *models.CreatePlatformInput) (*models.Platform, error)

	// DeleteAppSchema handles DELETE requests to /v1/apps/{appID}/schema
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteAppSchema(ctx context.Context, i *models.DeleteAppSchemaInput) error

	// GetAppSchema handles GET requests to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppSchema(ctx context.Context, appID string) (*models.AppSchema, error)

	// CreateAppSchema handles POST requests to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateAppSchema(ctx context.Context, i *models.CreateAppSchemaInput) (*models.AppSchema, error)

	// SetAppSchema handles PUT requests to /v1/apps/{appID}/schema
	//
	// 200: *models.AppSchema
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetAppSchema(ctx context.Context, i *models.SetAppSchemaInput) (*models.AppSchema, error)

	// GetSecrets handles GET requests to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSecrets(ctx context.Context, appID string) (*models.SecretConfig, error)

	// RevokeOldClientSecret handles PATCH requests to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	RevokeOldClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GenerateNewClientSecret handles POST requests to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GenerateNewClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// ResetClientSecret handles PUT requests to /v1/apps/{appID}/secrets
	//
	// 200: *models.SecretConfig
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	ResetClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error)

	// GetRecommendedSharing handles GET requests to /v1/apps/{appID}/sharing
	//
	// 200: *models.SharingRecommendations
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetRecommendedSharing(ctx context.Context, appID string) (*models.SharingRecommendations, error)

	// SetRecommendedSharing handles PUT requests to /v1/apps/{appID}/sharing
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SetRecommendedSharing(ctx context.Context, i *models.SetRecommendedSharingInput) error

	// UpdateAppIcon handles POST requests to /v1/apps/{appID}/update_icon
	//
	// 200: *models.Image
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 422: *models.UnprocessableEntity
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateAppIcon(ctx context.Context, i *models.UpdateAppIconInput) (*models.Image, error)

	// GetAllCategories handles GET requests to /v1/categories
	//
	// 200: *models.Categories
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllCategories(ctx context.Context) (*models.Categories, error)

	// GetKnownHosts handles GET requests to /v1/knownhosts
	//
	// 200: []models.KnownHost
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetKnownHosts(ctx context.Context) ([]models.KnownHost, error)

	// GetAllLibraryResources handles GET requests to /v1/libraryResources
	//
	// 200: *models.LibraryResources
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAllLibraryResources(ctx context.Context, i *models.GetAllLibraryResourcesInput) (*models.LibraryResources, error)

	// SearchLibraryResource handles GET requests to /v1/libraryResources/search
	//
	// 200: *models.LibraryResources
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	SearchLibraryResource(ctx context.Context, i *models.SearchLibraryResourceInput) (*models.LibraryResources, error)

	// GetLibraryResourceByShortname handles GET requests to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetLibraryResourceByShortname(ctx context.Context, i *models.GetLibraryResourceByShortnameInput) (*models.LibraryResource, error)

	// UpdateLibraryResourceByShortname handles PATCH requests to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	UpdateLibraryResourceByShortname(ctx context.Context, i *models.UpdateLibraryResourceByShortnameInput) (*models.LibraryResource, error)

	// CreateLibraryResource handles POST requests to /v1/libraryResources/{shortname}
	//
	// 200: *models.LibraryResource
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateLibraryResource(ctx context.Context, i *models.CreateLibraryResourceInput) (*models.LibraryResource, error)

	// DeleteLibraryResourceLink handles DELETE requests to /v1/libraryResources/{shortname}/link
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteLibraryResourceLink(ctx context.Context, shortname string) error

	// GetValidPermissions handles GET requests to /v1/permissions
	//
	// 200: *models.GetValidPermissionsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetValidPermissions(ctx context.Context) (*models.GetValidPermissionsResponse, error)

	// GetPlatforms handles GET requests to /v1/platforms
	// The server takes in the intersection of input parameters
	// 200: []models.Platform
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatforms(ctx context.Context, i *models.GetPlatformsInput) ([]models.Platform, error)

	// GetPlatformByClientID handles GET requests to /v1/platforms/{clientID}
	//
	// 200: *models.Platform
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetPlatformByClientID(ctx context.Context, clientID string) (*models.Platform, error)

	// GetAppsForAdmin handles GET requests to /v2/admins/{adminID}/apps
	//
	// 200: []models.AppForAdminResponse
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAppsForAdmin(ctx context.Context, adminID string) ([]models.AppForAdminResponse, error)

	// OverrideConfig handles POST requests to /v2/apps/{srcAppID}/override-config/{destAppID}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	OverrideConfig(ctx context.Context, i *models.OverrideConfigInput) error
}
