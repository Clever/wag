package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = fmt.Sprintf
var _ = url.QueryEscape
var _ = strconv.FormatInt
var _ = strings.Replace
var _ = validate.Maximum
var _ = strfmt.NewFormats

// HealthCheckInput holds the input parameters for a healthCheck operation.
type HealthCheckInput struct {
}

// Validate returns an error if any of the HealthCheckInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i HealthCheckInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i HealthCheckInput) Path() (string, error) {
	path := "/_health"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAdminsInput holds the input parameters for a getAdmins operation.
type GetAdminsInput struct {
	Email    *string
	Password *string
}

// Validate returns an error if any of the GetAdminsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAdminsInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetAdminsInput) Path() (string, error) {
	path := "/v1/admins"
	urlVals := url.Values{}

	if i.Email != nil {
		urlVals.Add("email", *i.Email)
	}

	if i.Password != nil {
		urlVals.Add("password", *i.Password)
	}

	return path + "?" + urlVals.Encode(), nil
}

// DeleteAdminInput holds the input parameters for a deleteAdmin operation.
type DeleteAdminInput struct {
	AdminID string
}

// ValidateDeleteAdminInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateDeleteAdminInput(adminID string) error {

	return nil
}

// DeleteAdminInputPath returns the URI path for the input.
func DeleteAdminInputPath(adminID string) (string, error) {
	path := "/v1/admins/{adminID}"
	urlVals := url.Values{}

	pathadminID := adminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAdminByIDInput holds the input parameters for a getAdminByID operation.
type GetAdminByIDInput struct {
	AdminID string
}

// ValidateGetAdminByIDInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAdminByIDInput(adminID string) error {

	return nil
}

// GetAdminByIDInputPath returns the URI path for the input.
func GetAdminByIDInputPath(adminID string) (string, error) {
	path := "/v1/admins/{adminID}"
	urlVals := url.Values{}

	pathadminID := adminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UpdateAdminInput holds the input parameters for a updateAdmin operation.
type UpdateAdminInput struct {
	AdminID string
	Admin   *PatchAdminRequest
}

// Validate returns an error if any of the UpdateAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdateAdminInput) Validate() error {

	if i.Admin != nil {
		if err := i.Admin.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdateAdminInput) Path() (string, error) {
	path := "/v1/admins/{adminID}"
	urlVals := url.Values{}

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateAdminInput holds the input parameters for a createAdmin operation.
type CreateAdminInput struct {
	CreateAdmin *CreateAdminRequest
	AdminID     string
}

// Validate returns an error if any of the CreateAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateAdminInput) Validate() error {

	if i.CreateAdmin != nil {
		if err := i.CreateAdmin.Validate(nil); err != nil {
			return err
		}
	}

	if err := validate.FormatOf("adminID", "path", "mongo-id", i.AdminID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i CreateAdminInput) Path() (string, error) {
	path := "/v1/admins/{adminID}"
	urlVals := url.Values{}

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAppsForAdminDeprecatedInput holds the input parameters for a getAppsForAdminDeprecated operation.
type GetAppsForAdminDeprecatedInput struct {
	AdminID string
}

// ValidateGetAppsForAdminDeprecatedInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAppsForAdminDeprecatedInput(adminID string) error {

	return nil
}

// GetAppsForAdminDeprecatedInputPath returns the URI path for the input.
func GetAppsForAdminDeprecatedInputPath(adminID string) (string, error) {
	path := "/v1/admins/{adminID}/apps"
	urlVals := url.Values{}

	pathadminID := adminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// VerifyCodeInput holds the input parameters for a verifyCode operation.
type VerifyCodeInput struct {
	Code       string
	Invalidate *bool
	AdminID    string
}

// Validate returns an error if any of the VerifyCodeInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i VerifyCodeInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i VerifyCodeInput) Path() (string, error) {
	path := "/v1/admins/{adminID}/confirmation_code"
	urlVals := url.Values{}

	urlVals.Add("code", i.Code)

	if i.Invalidate != nil {
		urlVals.Add("invalidate", strconv.FormatBool(*i.Invalidate))
	}

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateVerificationCodeInput holds the input parameters for a createVerificationCode operation.
type CreateVerificationCodeInput struct {
	Duration int32
	AdminID  string
}

// Validate returns an error if any of the CreateVerificationCodeInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateVerificationCodeInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i CreateVerificationCodeInput) Path() (string, error) {
	path := "/v1/admins/{adminID}/confirmation_code"
	urlVals := url.Values{}

	urlVals.Add("duration", strconv.FormatInt(int64(i.Duration), 10))

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// VerifyAdminEmailInput holds the input parameters for a verifyAdminEmail operation.
type VerifyAdminEmailInput struct {
	AdminID string
	Request *VerifyAdminEmailRequest
}

// Validate returns an error if any of the VerifyAdminEmailInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i VerifyAdminEmailInput) Validate() error {

	if i.Request != nil {
		if err := i.Request.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i VerifyAdminEmailInput) Path() (string, error) {
	path := "/v1/admins/{adminID}/verify_email"
	urlVals := url.Values{}

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAllAnalyticsAppsInput holds the input parameters for a getAllAnalyticsApps operation.
type GetAllAnalyticsAppsInput struct {
}

// Validate returns an error if any of the GetAllAnalyticsAppsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAllAnalyticsAppsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetAllAnalyticsAppsInput) Path() (string, error) {
	path := "/v1/analytics/apps"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAnalyticsAppByShortnameInput holds the input parameters for a getAnalyticsAppByShortname operation.
type GetAnalyticsAppByShortnameInput struct {
	Shortname string
}

// ValidateGetAnalyticsAppByShortnameInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAnalyticsAppByShortnameInput(shortname string) error {

	return nil
}

// GetAnalyticsAppByShortnameInputPath returns the URI path for the input.
func GetAnalyticsAppByShortnameInputPath(shortname string) (string, error) {
	path := "/v1/analytics/apps/{shortname}"
	urlVals := url.Values{}

	pathshortname := shortname
	if pathshortname == "" {
		err := fmt.Errorf("shortname cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{shortname}", pathshortname, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAllTrackableAppsInput holds the input parameters for a getAllTrackableApps operation.
type GetAllTrackableAppsInput struct {
}

// Validate returns an error if any of the GetAllTrackableAppsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAllTrackableAppsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetAllTrackableAppsInput) Path() (string, error) {
	path := "/v1/analytics/trackable_apps"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAnalyticsUsageUrlsInput holds the input parameters for a getAnalyticsUsageUrls operation.
type GetAnalyticsUsageUrlsInput struct {
}

// Validate returns an error if any of the GetAnalyticsUsageUrlsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAnalyticsUsageUrlsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetAnalyticsUsageUrlsInput) Path() (string, error) {
	path := "/v1/analytics/usageUrls"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAllUsageUrlsInput holds the input parameters for a getAllUsageUrls operation.
type GetAllUsageUrlsInput struct {
}

// Validate returns an error if any of the GetAllUsageUrlsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAllUsageUrlsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetAllUsageUrlsInput) Path() (string, error) {
	path := "/v1/appUniverse/usageUrls"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAppsInput holds the input parameters for a getApps operation.
type GetAppsInput struct {
	Ids           []string
	ClientId      *string
	ClientSecret  *string
	Shortname     *string
	BusinessToken *string
	Tags          []string
	SkipTags      []string
}

// Validate returns an error if any of the GetAppsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAppsInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetAppsInput) Path() (string, error) {
	path := "/v1/apps"
	urlVals := url.Values{}

	for _, v := range i.Ids {
		urlVals.Add("ids", v)
	}

	if i.ClientId != nil {
		urlVals.Add("clientId", *i.ClientId)
	}

	if i.ClientSecret != nil {
		urlVals.Add("clientSecret", *i.ClientSecret)
	}

	if i.Shortname != nil {
		urlVals.Add("shortname", *i.Shortname)
	}

	if i.BusinessToken != nil {
		urlVals.Add("businessToken", *i.BusinessToken)
	}

	for _, v := range i.Tags {
		urlVals.Add("tags", v)
	}

	for _, v := range i.SkipTags {
		urlVals.Add("skipTags", v)
	}

	return path + "?" + urlVals.Encode(), nil
}

// DeleteAppInput holds the input parameters for a deleteApp operation.
type DeleteAppInput struct {
	AppID string
}

// ValidateDeleteAppInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateDeleteAppInput(appID string) error {

	if err := validate.FormatOf("appID", "path", "mongo-id", appID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// DeleteAppInputPath returns the URI path for the input.
func DeleteAppInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAppByIDInput holds the input parameters for a getAppByID operation.
type GetAppByIDInput struct {
	AppID string
}

// ValidateGetAppByIDInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAppByIDInput(appID string) error {

	if err := validate.FormatOf("appID", "path", "mongo-id", appID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// GetAppByIDInputPath returns the URI path for the input.
func GetAppByIDInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UpdateAppInput holds the input parameters for a updateApp operation.
type UpdateAppInput struct {
	AppID                 string
	WithSchemaPropagation *bool
	App                   *PatchAppRequest
}

// Validate returns an error if any of the UpdateAppInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdateAppInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	if i.App != nil {
		if err := i.App.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdateAppInput) Path() (string, error) {
	path := "/v1/apps/{appID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	if i.WithSchemaPropagation != nil {
		urlVals.Add("withSchemaPropagation", strconv.FormatBool(*i.WithSchemaPropagation))
	}

	return path + "?" + urlVals.Encode(), nil
}

// CreateAppInput holds the input parameters for a createApp operation.
type CreateAppInput struct {
	App   *App
	AppID string
}

// Validate returns an error if any of the CreateAppInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateAppInput) Validate() error {

	if i.App != nil {
		if err := i.App.Validate(nil); err != nil {
			return err
		}
	}

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i CreateAppInput) Path() (string, error) {
	path := "/v1/apps/{appID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAdminsForAppInput holds the input parameters for a getAdminsForApp operation.
type GetAdminsForAppInput struct {
	AppID string
}

// ValidateGetAdminsForAppInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAdminsForAppInput(appID string) error {

	return nil
}

// GetAdminsForAppInputPath returns the URI path for the input.
func GetAdminsForAppInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/admins"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UnlinkAppAdminInput holds the input parameters for a unlinkAppAdmin operation.
type UnlinkAppAdminInput struct {
	AppID   string
	AdminID string
}

// Validate returns an error if any of the UnlinkAppAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UnlinkAppAdminInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i UnlinkAppAdminInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// LinkAppAdminInput holds the input parameters for a linkAppAdmin operation.
type LinkAppAdminInput struct {
	AppID       string
	AdminID     string
	Permissions *PermissionList
}

// Validate returns an error if any of the LinkAppAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i LinkAppAdminInput) Validate() error {

	if i.Permissions != nil {
		if err := i.Permissions.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i LinkAppAdminInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetGuideConfigInput holds the input parameters for a getGuideConfig operation.
type GetGuideConfigInput struct {
	AppID   string
	AdminID string
	GuideID string
}

// Validate returns an error if any of the GetGuideConfigInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetGuideConfigInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	if err := validate.FormatOf("adminID", "path", "mongo-id", i.AdminID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i GetGuideConfigInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}/guides/{guideID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	pathguideID := i.GuideID
	if pathguideID == "" {
		err := fmt.Errorf("guideID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{guideID}", pathguideID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// SetGuideConfigInput holds the input parameters for a setGuideConfig operation.
type SetGuideConfigInput struct {
	AppID       string
	AdminID     string
	GuideID     string
	GuideConfig *GuideConfig
}

// Validate returns an error if any of the SetGuideConfigInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SetGuideConfigInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	if err := validate.FormatOf("adminID", "path", "mongo-id", i.AdminID, strfmt.Default); err != nil {
		return err
	}

	if i.GuideConfig != nil {
		if err := i.GuideConfig.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i SetGuideConfigInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}/guides/{guideID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	pathguideID := i.GuideID
	if pathguideID == "" {
		err := fmt.Errorf("guideID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{guideID}", pathguideID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetPermissionsForAdminInput holds the input parameters for a getPermissionsForAdmin operation.
type GetPermissionsForAdminInput struct {
	AdminID string
	AppID   string
}

// Validate returns an error if any of the GetPermissionsForAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetPermissionsForAdminInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetPermissionsForAdminInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}/permissions"
	urlVals := url.Values{}

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// VerifyAppAdminInput holds the input parameters for a verifyAppAdmin operation.
type VerifyAppAdminInput struct {
	AppID    string
	AdminID  string
	Verified bool
}

// Validate returns an error if any of the VerifyAppAdminInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i VerifyAppAdminInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i VerifyAppAdminInput) Path() (string, error) {
	path := "/v1/apps/{appID}/admins/{adminID}/verify"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathadminID := i.AdminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	urlVals.Add("verified", strconv.FormatBool(i.Verified))

	return path + "?" + urlVals.Encode(), nil
}

// GenerateNewBusinessTokenInput holds the input parameters for a generateNewBusinessToken operation.
type GenerateNewBusinessTokenInput struct {
	AppID string
}

// ValidateGenerateNewBusinessTokenInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGenerateNewBusinessTokenInput(appID string) error {

	return nil
}

// GenerateNewBusinessTokenInputPath returns the URI path for the input.
func GenerateNewBusinessTokenInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/business_token"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetCertificationsInput holds the input parameters for a getCertifications operation.
type GetCertificationsInput struct {
	AppID           string
	SchoolYearStart int32
}

// Validate returns an error if any of the GetCertificationsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetCertificationsInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetCertificationsInput) Path() (string, error) {
	path := "/v1/apps/{appID}/certifications/{schoolYearStart}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathschoolYearStart := strconv.FormatInt(int64(i.SchoolYearStart), 10)
	if pathschoolYearStart == "" {
		err := fmt.Errorf("schoolYearStart cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{schoolYearStart}", pathschoolYearStart, -1)

	return path + "?" + urlVals.Encode(), nil
}

// SetCertificationsInput holds the input parameters for a setCertifications operation.
type SetCertificationsInput struct {
	AppID           string
	SchoolYearStart int32
	Certifications  *SetCertificationsRequest
}

// Validate returns an error if any of the SetCertificationsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SetCertificationsInput) Validate() error {

	if i.Certifications != nil {
		if err := i.Certifications.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i SetCertificationsInput) Path() (string, error) {
	path := "/v1/apps/{appID}/certifications/{schoolYearStart}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathschoolYearStart := strconv.FormatInt(int64(i.SchoolYearStart), 10)
	if pathschoolYearStart == "" {
		err := fmt.Errorf("schoolYearStart cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{schoolYearStart}", pathschoolYearStart, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetSetupStepInput holds the input parameters for a getSetupStep operation.
type GetSetupStepInput struct {
	AppID string
}

// ValidateGetSetupStepInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetSetupStepInput(appID string) error {

	return nil
}

// GetSetupStepInputPath returns the URI path for the input.
func GetSetupStepInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/customStep"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateSetupStepInput holds the input parameters for a createSetupStep operation.
type CreateSetupStepInput struct {
	AppID     string
	SetupStep *SetupStep
}

// Validate returns an error if any of the CreateSetupStepInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateSetupStepInput) Validate() error {

	if i.SetupStep != nil {
		if err := i.SetupStep.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i CreateSetupStepInput) Path() (string, error) {
	path := "/v1/apps/{appID}/customStep"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetDataRulesInput holds the input parameters for a getDataRules operation.
type GetDataRulesInput struct {
	AppID string
}

// ValidateGetDataRulesInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetDataRulesInput(appID string) error {

	return nil
}

// GetDataRulesInputPath returns the URI path for the input.
func GetDataRulesInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/data_rules"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// SetDataRulesInput holds the input parameters for a setDataRules operation.
type SetDataRulesInput struct {
	AppID string
	Rules *SetDataRulesRequest
}

// Validate returns an error if any of the SetDataRulesInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SetDataRulesInput) Validate() error {

	if i.Rules != nil {
		if err := i.Rules.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i SetDataRulesInput) Path() (string, error) {
	path := "/v1/apps/{appID}/data_rules"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetManagersInput holds the input parameters for a getManagers operation.
type GetManagersInput struct {
	AppID string
}

// ValidateGetManagersInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetManagersInput(appID string) error {

	return nil
}

// GetManagersInputPath returns the URI path for the input.
func GetManagersInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/managers"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetOnboardingInput holds the input parameters for a getOnboarding operation.
type GetOnboardingInput struct {
	AppID string
}

// ValidateGetOnboardingInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetOnboardingInput(appID string) error {

	return nil
}

// GetOnboardingInputPath returns the URI path for the input.
func GetOnboardingInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/onboarding"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UpdateOnboardingInput holds the input parameters for a updateOnboarding operation.
type UpdateOnboardingInput struct {
	AppID  string
	Update *UpdateOnboardingRequest
}

// Validate returns an error if any of the UpdateOnboardingInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdateOnboardingInput) Validate() error {

	if i.Update != nil {
		if err := i.Update.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdateOnboardingInput) Path() (string, error) {
	path := "/v1/apps/{appID}/onboarding"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// InitializeOnboardingInput holds the input parameters for a initializeOnboarding operation.
type InitializeOnboardingInput struct {
	AppID string
}

// ValidateInitializeOnboardingInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateInitializeOnboardingInput(appID string) error {

	return nil
}

// InitializeOnboardingInputPath returns the URI path for the input.
func InitializeOnboardingInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/onboarding"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// DeletePlatformInput holds the input parameters for a deletePlatform operation.
type DeletePlatformInput struct {
	AppID    string
	ClientID string
}

// Validate returns an error if any of the DeletePlatformInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i DeletePlatformInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i DeletePlatformInput) Path() (string, error) {
	path := "/v1/apps/{appID}/platform/{clientID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathclientID := i.ClientID
	if pathclientID == "" {
		err := fmt.Errorf("clientID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{clientID}", pathclientID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UpdatePlatformInput holds the input parameters for a updatePlatform operation.
type UpdatePlatformInput struct {
	AppID    string
	ClientID string
	Request  *PatchPlatformRequest
}

// Validate returns an error if any of the UpdatePlatformInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdatePlatformInput) Validate() error {

	if i.Request != nil {
		if err := i.Request.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdatePlatformInput) Path() (string, error) {
	path := "/v1/apps/{appID}/platform/{clientID}"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	pathclientID := i.ClientID
	if pathclientID == "" {
		err := fmt.Errorf("clientID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{clientID}", pathclientID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetPlatformsByAppIDInput holds the input parameters for a getPlatformsByAppID operation.
type GetPlatformsByAppIDInput struct {
	AppID string
}

// ValidateGetPlatformsByAppIDInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetPlatformsByAppIDInput(appID string) error {

	return nil
}

// GetPlatformsByAppIDInputPath returns the URI path for the input.
func GetPlatformsByAppIDInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/platforms"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreatePlatformInput holds the input parameters for a createPlatform operation.
type CreatePlatformInput struct {
	AppID   string
	Request *CreatePlatformRequest
}

// Validate returns an error if any of the CreatePlatformInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreatePlatformInput) Validate() error {

	if i.Request != nil {
		if err := i.Request.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i CreatePlatformInput) Path() (string, error) {
	path := "/v1/apps/{appID}/platforms"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// DeleteAppSchemaInput holds the input parameters for a deleteAppSchema operation.
type DeleteAppSchemaInput struct {
	AppID           string
	DeleteDataRules *bool
}

// Validate returns an error if any of the DeleteAppSchemaInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i DeleteAppSchemaInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i DeleteAppSchemaInput) Path() (string, error) {
	path := "/v1/apps/{appID}/schema"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	if i.DeleteDataRules != nil {
		urlVals.Add("deleteDataRules", strconv.FormatBool(*i.DeleteDataRules))
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetAppSchemaInput holds the input parameters for a getAppSchema operation.
type GetAppSchemaInput struct {
	AppID string
}

// ValidateGetAppSchemaInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAppSchemaInput(appID string) error {

	if err := validate.FormatOf("appID", "path", "mongo-id", appID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// GetAppSchemaInputPath returns the URI path for the input.
func GetAppSchemaInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/schema"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateAppSchemaInput holds the input parameters for a createAppSchema operation.
type CreateAppSchemaInput struct {
	AppID           string
	SkipPropagation *bool
	UpdateDataRules *bool
}

// Validate returns an error if any of the CreateAppSchemaInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateAppSchemaInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i CreateAppSchemaInput) Path() (string, error) {
	path := "/v1/apps/{appID}/schema"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	if i.SkipPropagation != nil {
		urlVals.Add("skipPropagation", strconv.FormatBool(*i.SkipPropagation))
	}

	if i.UpdateDataRules != nil {
		urlVals.Add("updateDataRules", strconv.FormatBool(*i.UpdateDataRules))
	}

	return path + "?" + urlVals.Encode(), nil
}

// SetAppSchemaInput holds the input parameters for a setAppSchema operation.
type SetAppSchemaInput struct {
	AppID           string
	SkipPropagation *bool
	UpdateDataRules *bool
	AppSchema       *AppSchema
}

// Validate returns an error if any of the SetAppSchemaInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SetAppSchemaInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	if i.AppSchema != nil {
		if err := i.AppSchema.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i SetAppSchemaInput) Path() (string, error) {
	path := "/v1/apps/{appID}/schema"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	if i.SkipPropagation != nil {
		urlVals.Add("skipPropagation", strconv.FormatBool(*i.SkipPropagation))
	}

	if i.UpdateDataRules != nil {
		urlVals.Add("updateDataRules", strconv.FormatBool(*i.UpdateDataRules))
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetSecretsInput holds the input parameters for a getSecrets operation.
type GetSecretsInput struct {
	AppID string
}

// ValidateGetSecretsInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetSecretsInput(appID string) error {

	return nil
}

// GetSecretsInputPath returns the URI path for the input.
func GetSecretsInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/secrets"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// RevokeOldClientSecretInput holds the input parameters for a revokeOldClientSecret operation.
type RevokeOldClientSecretInput struct {
	AppID string
}

// ValidateRevokeOldClientSecretInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateRevokeOldClientSecretInput(appID string) error {

	return nil
}

// RevokeOldClientSecretInputPath returns the URI path for the input.
func RevokeOldClientSecretInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/secrets"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GenerateNewClientSecretInput holds the input parameters for a generateNewClientSecret operation.
type GenerateNewClientSecretInput struct {
	AppID string
}

// ValidateGenerateNewClientSecretInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGenerateNewClientSecretInput(appID string) error {

	return nil
}

// GenerateNewClientSecretInputPath returns the URI path for the input.
func GenerateNewClientSecretInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/secrets"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// ResetClientSecretInput holds the input parameters for a resetClientSecret operation.
type ResetClientSecretInput struct {
	AppID string
}

// ValidateResetClientSecretInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateResetClientSecretInput(appID string) error {

	return nil
}

// ResetClientSecretInputPath returns the URI path for the input.
func ResetClientSecretInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/secrets"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetRecommendedSharingInput holds the input parameters for a getRecommendedSharing operation.
type GetRecommendedSharingInput struct {
	AppID string
}

// ValidateGetRecommendedSharingInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetRecommendedSharingInput(appID string) error {

	return nil
}

// GetRecommendedSharingInputPath returns the URI path for the input.
func GetRecommendedSharingInputPath(appID string) (string, error) {
	path := "/v1/apps/{appID}/sharing"
	urlVals := url.Values{}

	pathappID := appID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// SetRecommendedSharingInput holds the input parameters for a setRecommendedSharing operation.
type SetRecommendedSharingInput struct {
	AppID           string
	Recommendations *SharingRecommendations
}

// Validate returns an error if any of the SetRecommendedSharingInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SetRecommendedSharingInput) Validate() error {

	if i.Recommendations != nil {
		if err := i.Recommendations.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i SetRecommendedSharingInput) Path() (string, error) {
	path := "/v1/apps/{appID}/sharing"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// UpdateAppIconInput holds the input parameters for a updateAppIcon operation.
type UpdateAppIconInput struct {
	AppID string
	App   *UpdateAppIconRequest
}

// Validate returns an error if any of the UpdateAppIconInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdateAppIconInput) Validate() error {

	if err := validate.FormatOf("appID", "path", "mongo-id", i.AppID, strfmt.Default); err != nil {
		return err
	}

	if i.App != nil {
		if err := i.App.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdateAppIconInput) Path() (string, error) {
	path := "/v1/apps/{appID}/update_icon"
	urlVals := url.Values{}

	pathappID := i.AppID
	if pathappID == "" {
		err := fmt.Errorf("appID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{appID}", pathappID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAllCategoriesInput holds the input parameters for a getAllCategories operation.
type GetAllCategoriesInput struct {
}

// Validate returns an error if any of the GetAllCategoriesInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAllCategoriesInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetAllCategoriesInput) Path() (string, error) {
	path := "/v1/categories"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetKnownHostsInput holds the input parameters for a getKnownHosts operation.
type GetKnownHostsInput struct {
}

// Validate returns an error if any of the GetKnownHostsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetKnownHostsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetKnownHostsInput) Path() (string, error) {
	path := "/v1/knownhosts"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetAllLibraryResourcesInput holds the input parameters for a getAllLibraryResources operation.
type GetAllLibraryResourcesInput struct {
	Category       *string
	IncludeDevApps *bool
	IncludeLinks   *bool
}

// Validate returns an error if any of the GetAllLibraryResourcesInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAllLibraryResourcesInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetAllLibraryResourcesInput) Path() (string, error) {
	path := "/v1/libraryResources"
	urlVals := url.Values{}

	if i.Category != nil {
		urlVals.Add("category", *i.Category)
	}

	if i.IncludeDevApps != nil {
		urlVals.Add("includeDevApps", strconv.FormatBool(*i.IncludeDevApps))
	}

	if i.IncludeLinks != nil {
		urlVals.Add("includeLinks", strconv.FormatBool(*i.IncludeLinks))
	}

	return path + "?" + urlVals.Encode(), nil
}

// SearchLibraryResourceInput holds the input parameters for a searchLibraryResource operation.
type SearchLibraryResourceInput struct {
	SearchTerm        string
	ShowInLibraryOnly *bool
	IncludeLinks      *bool
}

// Validate returns an error if any of the SearchLibraryResourceInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i SearchLibraryResourceInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i SearchLibraryResourceInput) Path() (string, error) {
	path := "/v1/libraryResources/search"
	urlVals := url.Values{}

	urlVals.Add("searchTerm", i.SearchTerm)

	if i.ShowInLibraryOnly != nil {
		urlVals.Add("showInLibraryOnly", strconv.FormatBool(*i.ShowInLibraryOnly))
	}

	if i.IncludeLinks != nil {
		urlVals.Add("includeLinks", strconv.FormatBool(*i.IncludeLinks))
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetLibraryResourceByShortnameInput holds the input parameters for a getLibraryResourceByShortname operation.
type GetLibraryResourceByShortnameInput struct {
	Shortname      string
	IncludeDevApps *bool
	IncludeLinks   *bool
}

// Validate returns an error if any of the GetLibraryResourceByShortnameInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetLibraryResourceByShortnameInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetLibraryResourceByShortnameInput) Path() (string, error) {
	path := "/v1/libraryResources/{shortname}"
	urlVals := url.Values{}

	pathshortname := i.Shortname
	if pathshortname == "" {
		err := fmt.Errorf("shortname cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{shortname}", pathshortname, -1)

	if i.IncludeDevApps != nil {
		urlVals.Add("includeDevApps", strconv.FormatBool(*i.IncludeDevApps))
	}

	if i.IncludeLinks != nil {
		urlVals.Add("includeLinks", strconv.FormatBool(*i.IncludeLinks))
	}

	return path + "?" + urlVals.Encode(), nil
}

// UpdateLibraryResourceByShortnameInput holds the input parameters for a updateLibraryResourceByShortname operation.
type UpdateLibraryResourceByShortnameInput struct {
	Shortname       string
	LibraryResource *PatchLibraryResourceRequest
}

// Validate returns an error if any of the UpdateLibraryResourceByShortnameInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i UpdateLibraryResourceByShortnameInput) Validate() error {

	if i.LibraryResource != nil {
		if err := i.LibraryResource.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i UpdateLibraryResourceByShortnameInput) Path() (string, error) {
	path := "/v1/libraryResources/{shortname}"
	urlVals := url.Values{}

	pathshortname := i.Shortname
	if pathshortname == "" {
		err := fmt.Errorf("shortname cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{shortname}", pathshortname, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateLibraryResourceInput holds the input parameters for a createLibraryResource operation.
type CreateLibraryResourceInput struct {
	Shortname       string
	LibraryResource *CreateLibraryResourceRequest
}

// Validate returns an error if any of the CreateLibraryResourceInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateLibraryResourceInput) Validate() error {

	if i.LibraryResource != nil {
		if err := i.LibraryResource.Validate(nil); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i CreateLibraryResourceInput) Path() (string, error) {
	path := "/v1/libraryResources/{shortname}"
	urlVals := url.Values{}

	pathshortname := i.Shortname
	if pathshortname == "" {
		err := fmt.Errorf("shortname cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{shortname}", pathshortname, -1)

	return path + "?" + urlVals.Encode(), nil
}

// DeleteLibraryResourceLinkInput holds the input parameters for a deleteLibraryResourceLink operation.
type DeleteLibraryResourceLinkInput struct {
	Shortname string
}

// ValidateDeleteLibraryResourceLinkInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateDeleteLibraryResourceLinkInput(shortname string) error {

	return nil
}

// DeleteLibraryResourceLinkInputPath returns the URI path for the input.
func DeleteLibraryResourceLinkInputPath(shortname string) (string, error) {
	path := "/v1/libraryResources/{shortname}/link"
	urlVals := url.Values{}

	pathshortname := shortname
	if pathshortname == "" {
		err := fmt.Errorf("shortname cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{shortname}", pathshortname, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetValidPermissionsInput holds the input parameters for a getValidPermissions operation.
type GetValidPermissionsInput struct {
}

// Validate returns an error if any of the GetValidPermissionsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetValidPermissionsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetValidPermissionsInput) Path() (string, error) {
	path := "/v1/permissions"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetPlatformsInput holds the input parameters for a getPlatforms operation.
type GetPlatformsInput struct {
	AppIds []string
	Name   *string
}

// Validate returns an error if any of the GetPlatformsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetPlatformsInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetPlatformsInput) Path() (string, error) {
	path := "/v1/platforms"
	urlVals := url.Values{}

	for _, v := range i.AppIds {
		urlVals.Add("appIds", v)
	}

	if i.Name != nil {
		urlVals.Add("name", *i.Name)
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetPlatformByClientIDInput holds the input parameters for a getPlatformByClientID operation.
type GetPlatformByClientIDInput struct {
	ClientID string
}

// ValidateGetPlatformByClientIDInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetPlatformByClientIDInput(clientID string) error {

	return nil
}

// GetPlatformByClientIDInputPath returns the URI path for the input.
func GetPlatformByClientIDInputPath(clientID string) (string, error) {
	path := "/v1/platforms/{clientID}"
	urlVals := url.Values{}

	pathclientID := clientID
	if pathclientID == "" {
		err := fmt.Errorf("clientID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{clientID}", pathclientID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetAppsForAdminInput holds the input parameters for a getAppsForAdmin operation.
type GetAppsForAdminInput struct {
	AdminID string
}

// ValidateGetAppsForAdminInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetAppsForAdminInput(adminID string) error {

	return nil
}

// GetAppsForAdminInputPath returns the URI path for the input.
func GetAppsForAdminInputPath(adminID string) (string, error) {
	path := "/v2/admins/{adminID}/apps"
	urlVals := url.Values{}

	pathadminID := adminID
	if pathadminID == "" {
		err := fmt.Errorf("adminID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{adminID}", pathadminID, -1)

	return path + "?" + urlVals.Encode(), nil
}

// OverrideConfigInput holds the input parameters for a overrideConfig operation.
type OverrideConfigInput struct {
	SrcAppID  string
	DestAppID string
}

// Validate returns an error if any of the OverrideConfigInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i OverrideConfigInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i OverrideConfigInput) Path() (string, error) {
	path := "/v2/apps/{srcAppID}/override-config/{destAppID}"
	urlVals := url.Values{}

	pathsrcAppID := i.SrcAppID
	if pathsrcAppID == "" {
		err := fmt.Errorf("srcAppID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{srcAppID}", pathsrcAppID, -1)

	pathdestAppID := i.DestAppID
	if pathdestAppID == "" {
		err := fmt.Errorf("destAppID cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{destAppID}", pathdestAppID, -1)

	return path + "?" + urlVals.Encode(), nil
}
