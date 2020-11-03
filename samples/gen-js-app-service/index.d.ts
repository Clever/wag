import { Span, Tracer } from "opentracing";
import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
	/** The timeout to use for all client requests, in milliseconds. */
	timeout?: number;
	/** An OpenTracing span - For example from the parent request */
	span?: Span;
	/** The logic to determine which requests to retry, as well as how many times to retry. */
  retryPolicy?: RetryPolicy;
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
  forEachAsync(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CircuitOptions {
	/** When set to true the circuit will always be closed. Default: true. */
	forceClosed?: boolean;
	/** The maximum number of concurrent requests the client can make at the same time. Default: 100. */
	maxConcurrentRequests?: number;
	/** The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. */
	requestVolumeThreshold?: number;
	/** How long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. */
	sleepWindow?: number;
	/** The threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. */
  errorPercentThreshold?: number;
}

interface GenericOptions {
	/** The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. */
	timeout?: number;
	/** Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. */
	keepalive?: boolean;
	/** The logic to determine which requests to retry, as well as how many times to retry. */
	retryPolicy?: RetryPolicy;
	/** The Kayvee logger to use in the client. */
	logger?: Logger;
	/** The OpenTracing Tracer to use. Defaults to the OpenTracing globalTracer. */
	tracer?: Tracer;
	/** Options for constructing the client's circuit breaker. */
	circuit?: CircuitOptions;
	/** Overrides the default service name. This is necessary if the same client is used multiple times, but with different settings, such as with sso- clients. */
  serviceName?: string;
}

interface DiscoveryOptions {
	/** Use clever-discovery to locate the server. Must provide this or the address argument. */
	discovery: true;
  address?: undefined;
}

interface AddressOptions {
	discovery?: false;
	/** URL where the server is located. Must provide this or the discovery argument. */
  address: string;
}

type AppServiceOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = AppService.Models

declare class AppService {
	/**
	* Create a new client object.
	* @param options - Options for constructing a client object.
	*/
  constructor(options: AppServiceOptions);

  
  /**
	  Checks if the service is healthy
@throws BadRequest
@throws InternalError
	*/
	healthCheck(options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Gets app admins
@throws BadRequest
@throws InternalError
	*/
	getAdmins(params: models.GetAdminsParams, options?: RequestOptions, cb?: Callback<models.Admin[]>): Promise<models.Admin[]>
  
  /**
	  Delete an admin
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	deleteAdmin(adminID: string, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get admin by ID
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAdminByID(adminID: string, options?: RequestOptions, cb?: Callback<models.Admin>): Promise<models.Admin>
  
  /**
	  Update an admin
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	updateAdmin(params: models.UpdateAdminParams, options?: RequestOptions, cb?: Callback<models.Admin>): Promise<models.Admin>
  
  /**
	  Creates an app admin
@throws BadRequest
@throws InternalError
	*/
	createAdmin(params: models.CreateAdminParams, options?: RequestOptions, cb?: Callback<models.Admin>): Promise<models.Admin>
  
  /**
	  Verify and possibly remove the verification code
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	verifyCode(params: models.VerifyCodeParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Create verification code
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	createVerificationCode(params: models.CreateVerificationCodeParams, options?: RequestOptions, cb?: Callback<models.VerificationCodeResponse>): Promise<models.VerificationCodeResponse>
  
  /**
	  set the verified email of an admin
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	verifyAdminEmail(params: models.VerifyAdminEmailParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Returns all apps that are not found in district or library data (for p4a)
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAllAnalyticsApps(options?: RequestOptions, cb?: Callback<models.AnalyticsApps>): Promise<models.AnalyticsApps>
  
  /**
	  Returns an analytics app matching the shortname
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAnalyticsAppByShortname(shortname: string, options?: RequestOptions, cb?: Callback<models.AnalyticsApp>): Promise<models.AnalyticsApp>
  
  /**
	  Returns all apps that are used in analytics tracking
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAllTrackableApps(options?: RequestOptions, cb?: Callback<models.TrackableApps>): Promise<models.TrackableApps>
  
  /**
	  Returns all usage urls
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAnalyticsUsageUrls(options?: RequestOptions, cb?: Callback<models.UsageUrls>): Promise<models.UsageUrls>
  
  /**
	  Returns all usage urls
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAllUsageUrls(options?: RequestOptions, cb?: Callback<models.UsageUrls>): Promise<models.UsageUrls>
  
  /**
	  Gets applications filtered by the parameters
		
		The server takes in the intersection of input parameters
@throws BadRequest
@throws InternalError
	*/
	getApps(params: models.GetAppsParams, options?: RequestOptions, cb?: Callback<models.App[]>): Promise<models.App[]>
  
  /**
	  Delete an application
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	deleteApp(appID: string, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get application by ID
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAppByID(appID: string, options?: RequestOptions, cb?: Callback<models.App>): Promise<models.App>
  
  /**
	  Update an application
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	updateApp(params: models.UpdateAppParams, options?: RequestOptions, cb?: Callback<models.App>): Promise<models.App>
  
  /**
	  Creates an app
@throws BadRequest
@throws InternalError
	*/
	createApp(params: models.CreateAppParams, options?: RequestOptions, cb?: Callback<models.App>): Promise<models.App>
  
  /**
	  Admins for an app
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAdminsForApp(appID: string, options?: RequestOptions, cb?: Callback<models.AppAdminResponse[]>): Promise<models.AppAdminResponse[]>
  
  /**
	  Remove an admin from an app
@throws BadRequest
@throws Forbidden
@throws NotFound
@throws InternalError
	*/
	unlinkAppAdmin(params: models.UnlinkAppAdminParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Add an admin to an app
@throws BadRequest
@throws Forbidden
@throws NotFound
@throws InternalError
	*/
	linkAppAdmin(params: models.LinkAppAdminParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get a guide for an admin for an app
@throws BadRequest
@throws Forbidden
@throws NotFound
@throws InternalError
	*/
	getGuideConfig(params: models.GetGuideConfigParams, options?: RequestOptions, cb?: Callback<models.GuideConfig>): Promise<models.GuideConfig>
  
  /**
	  Add a guide for an admin for an app
@throws BadRequest
@throws Forbidden
@throws NotFound
@throws InternalError
	*/
	setGuideConfig(params: models.SetGuideConfigParams, options?: RequestOptions, cb?: Callback<models.GuideConfig>): Promise<models.GuideConfig>
  
  /**
	  Permissions for an app admin.
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getPermissionsForAdmin(params: models.GetPermissionsForAdminParams, options?: RequestOptions, cb?: Callback<models.PermissionList>): Promise<models.PermissionList>
  
  /**
	  Updates an admin's verification for an app
@throws BadRequest
@throws Forbidden
@throws NotFound
@throws InternalError
	*/
	verifyAppAdmin(params: models.VerifyAppAdminParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Generate a new business token and immediately revoke the previous one (if it exists).
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	generateNewBusinessToken(appID: string, options?: RequestOptions, cb?: Callback<models.SecretConfig>): Promise<models.SecretConfig>
  
  /**
	  Returns the given app's certifications for the given school year.
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getCertifications(params: models.GetCertificationsParams, options?: RequestOptions, cb?: Callback<models.Certifications>): Promise<models.Certifications>
  
  /**
	  Sets the given app's certifications for the given school year.
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	setCertifications(params: models.SetCertificationsParams, options?: RequestOptions, cb?: Callback<models.Certifications>): Promise<models.Certifications>
  
  /**
	  Returns an app-specific setup step
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getSetupStep(appID: string, options?: RequestOptions, cb?: Callback<models.SetupStep>): Promise<models.SetupStep>
  
  /**
	  Creates a new custom setup step
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	createSetupStep(params: models.CreateSetupStepParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get data rules
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getDataRules(appID: string, options?: RequestOptions, cb?: Callback<models.DataRule[]>): Promise<models.DataRule[]>
  
  /**
	  Set data rules
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	setDataRules(params: models.SetDataRulesParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get account and relationship managers
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getManagers(appID: string, options?: RequestOptions, cb?: Callback<models.Managers>): Promise<models.Managers>
  
  /**
	  Returns an app's onboarding progress
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getOnboarding(appID: string, options?: RequestOptions, cb?: Callback<models.Onboarding>): Promise<models.Onboarding>
  
  /**
	  Updates an app's onboarding progress
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	updateOnboarding(params: models.UpdateOnboardingParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Initializes an app's onboarding progress
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	initializeOnboarding(appID: string, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Delete a platform
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	deletePlatform(params: models.DeletePlatformParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Update a platform
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	updatePlatform(params: models.UpdatePlatformParams, options?: RequestOptions, cb?: Callback<models.Platform>): Promise<models.Platform>
  
  /**
	  Returns list of platforms for an application
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getPlatformsByAppID(appID: string, options?: RequestOptions, cb?: Callback<models.Platform[]>): Promise<models.Platform[]>
  
  /**
	  Creates a platform
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	createPlatform(params: models.CreatePlatformParams, options?: RequestOptions, cb?: Callback<models.Platform>): Promise<models.Platform>
  
  /**
	  Delete field settings specified by the app
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	deleteAppSchema(params: models.DeleteAppSchemaParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get field settings specified by the app
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAppSchema(appID: string, options?: RequestOptions, cb?: Callback<models.AppSchema>): Promise<models.AppSchema>
  
  /**
	  Creates the schema for an app. Will propagate app schema changes to all connections associated with this app
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	createAppSchema(params: models.CreateAppSchemaParams, options?: RequestOptions, cb?: Callback<models.AppSchema>): Promise<models.AppSchema>
  
  /**
	  Set field settings for an app. Will propagate app schema changes to all connections associated with this app
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	setAppSchema(params: models.SetAppSchemaParams, options?: RequestOptions, cb?: Callback<models.AppSchema>): Promise<models.AppSchema>
  
  /**
	  Get secret information
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getSecrets(appID: string, options?: RequestOptions, cb?: Callback<models.SecretConfig>): Promise<models.SecretConfig>
  
  /**
	  Revokes the old client secret
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	revokeOldClientSecret(appID: string, options?: RequestOptions, cb?: Callback<models.SecretConfig>): Promise<models.SecretConfig>
  
  /**
	  Generate a new client secret. This creates a new client secret, but retains the previous secret so that clients can migrate their code to the secret.
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	generateNewClientSecret(appID: string, options?: RequestOptions, cb?: Callback<models.SecretConfig>): Promise<models.SecretConfig>
  
  /**
	  Hard resets the client secret
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	resetClientSecret(appID: string, options?: RequestOptions, cb?: Callback<models.SecretConfig>): Promise<models.SecretConfig>
  
  /**
	  Get recommended sharing settings
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getRecommendedSharing(appID: string, options?: RequestOptions, cb?: Callback<models.SharingRecommendations>): Promise<models.SharingRecommendations>
  
  /**
	  Set recommended sharing settings
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	setRecommendedSharing(params: models.SetRecommendedSharingParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Update icon
@throws BadRequest
@throws NotFound
@throws UnprocessableEntity
@throws InternalError
	*/
	updateAppIcon(params: models.UpdateAppIconParams, options?: RequestOptions, cb?: Callback<models.Image>): Promise<models.Image>
  
  /**
	  All valid categories an app can belong to
@throws BadRequest
@throws InternalError
	*/
	getAllCategories(options?: RequestOptions, cb?: Callback<models.Categories>): Promise<models.Categories>
  
  /**
	  Get a list of known hosts for apps
@throws BadRequest
@throws InternalError
	*/
	getKnownHosts(options?: RequestOptions, cb?: Callback<models.KnownHost[]>): Promise<models.KnownHost[]>
  
  /**
	  Returns all library resources
@throws BadRequest
@throws InternalError
	*/
	getAllLibraryResources(params: models.GetAllLibraryResourcesParams, options?: RequestOptions, cb?: Callback<models.LibraryResources>): Promise<models.LibraryResources>
  
  /**
	  Returns a list of library resource matching a given query
@throws BadRequest
@throws InternalError
	*/
	searchLibraryResource(params: models.SearchLibraryResourceParams, options?: RequestOptions, cb?: Callback<models.LibraryResources>): Promise<models.LibraryResources>
  
  /**
	  Returns a library resource with a given shortname
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getLibraryResourceByShortname(params: models.GetLibraryResourceByShortnameParams, options?: RequestOptions, cb?: Callback<models.LibraryResource>): Promise<models.LibraryResource>
  
  /**
	  Updates the library resource with a given shortname
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	updateLibraryResourceByShortname(params: models.UpdateLibraryResourceByShortnameParams, options?: RequestOptions, cb?: Callback<models.LibraryResource>): Promise<models.LibraryResource>
  
  /**
	  Creates a library resource with the given shortname
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	createLibraryResource(params: models.CreateLibraryResourceParams, options?: RequestOptions, cb?: Callback<models.LibraryResource>): Promise<models.LibraryResource>
  
  /**
	  Deletes a library link with a given shortname
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	deleteLibraryResourceLink(shortname: string, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  /**
	  Get all valid permissions.
@throws BadRequest
@throws InternalError
	*/
	getValidPermissions(options?: RequestOptions, cb?: Callback<models.GetValidPermissionsResponse>): Promise<models.GetValidPermissionsResponse>
  
  /**
	  Gets platforms filtered by the parameters
		
		The server takes in the intersection of input parameters
@throws BadRequest
@throws InternalError
	*/
	getPlatforms(params: models.GetPlatformsParams, options?: RequestOptions, cb?: Callback<models.Platform[]>): Promise<models.Platform[]>
  
  /**
	  Get platform by client ID
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getPlatformByClientID(clientID: string, options?: RequestOptions, cb?: Callback<models.Platform>): Promise<models.Platform>
  
  /**
	  Apps an admin is associated with
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	getAppsForAdmin(adminID: string, options?: RequestOptions, cb?: Callback<models.AppForAdminResponse[]>): Promise<models.AppForAdminResponse[]>
  
  /**
	  Override one app's config with another's
@throws BadRequest
@throws NotFound
@throws InternalError
	*/
	overrideConfig(params: models.OverrideConfigParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
}

declare namespace AppService {
  const RetryPolicies: {
    Single: RetryPolicy;
    Exponential: RetryPolicy;
    None: RetryPolicy;
  }

  const DefaultCircuitOptions: CircuitOptions;

  namespace Errors {
    interface ErrorBody {
      message: string;
      [key: string]: any;
    }

    
    /** The request was bad... */
class BadRequest {
  /** The error code, if there is one */
code?: models.ErrorCode;
  message?: string;

  constructor(body: ErrorBody);
}
    
    class InternalError {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class NotFound {
  code?: models.ErrorCode;
  message?: string;

  constructor(body: ErrorBody);
}
    
    class Forbidden {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class UnprocessableEntity {
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
    type Admin = {
  email?: string;
  id?: string;
  name?: string;
  phone?: string;
  twoFactor?: TwoFactor;
  verifiedEmail?: string;
};
    
    /** App metadata for apps that don't exist on Clever (for p4a) */
type AnalyticsApp = {
  bannerIconURL?: string;
  categories?: ResourceCategories;
  certifications?: ResourceCertifications;
  cleverID?: string;
  description?: string;
  devices?: ResourceDevices;
  featured?: number;
  features?: ResourceFeatures;
  gradeLevels?: ResourceGradeLevels;
  iconURL?: string;
  id: string;
  insights?: InsightsMetadata;
  installCount?: number;
  libraryID?: string;
  libraryIntegration?: LibraryIntegrationType;
  marketingCollateralStatus?: MarketingCollateralStatus;
  name?: string;
  pricing?: ResourcePricing;
  privacyPolicyURL?: string;
  searchTags?: string[];
  shortname?: string;
  showInLibrary?: boolean;
  stripeUserID?: string;
  tagline?: string;
  tiers?: Tiers;
  tosURL?: string;
  universalDSASignatory?: boolean;
  url?: string;
};
    
    type AnalyticsApps = AnalyticsApp[];
    
    /** Represents an application returned by the app service. */
type App = {
  altNames?: string[];
  categories?: string[];
  clientId?: string;
  created?: string;
  customerSolutionsNotes?: string;
  defaultAccessTier?: string;
  description?: string;
  edsurgeUrl?: string;
  id?: string;
  image?: Image;
  inviteMessage?: string;
  iosRedirectUri?: string;
  name?: string;
  nonOwnerImpersonationDisabled?: boolean;
  redirectUris?: string[];
  requiredScopes?: string[];
  salesContact?: string;
  shortname?: string;
  supportContact?: string;
  supportContactName?: string;
  supportedUserTypes?: UserType[];
  tags?: string[];
  versions?: VersionSettings;
  websiteUrl?: string;
};
    
    /** Admin with the data specific to a particular app */
type AppAdminResponse = {
  admin?: Admin;
  dateAdded?: string;
  isOwner?: boolean;
  permissions?: PermissionList;
  verificationPending?: boolean;
};
    
    /** App with the data specific to a particular admin */
type AppForAdminResponse = {
  app?: App;
  dateAdded?: string;
  isOwner?: boolean;
  permissions?: PermissionList;
  verificationPending?: boolean;
};
    
    /** Tracks app onboarding data for library apps */
type AppOnboarding = {
  announced?: ("notScheduled" | "scheduled" | "announced");
};
    
    type AppSchema = {
  collectionSchemas?: CollectionSchema[];
  needsPropagation?: boolean;
};
    
    type Categories = string[];
    
    /** Integration certifications for an app, for a specific school year. */
type Certifications = {
  appID?: string;
  /** DEPRECATED - use certS3PathIntegration and certS3PathTradeshow instead. This field is mapped to certS3PathTradeshow. */
certS3Path?: string;
  certS3PathIntegration?: string;
  certS3PathTradeshow?: string;
  dateUpdated?: string;
  events?: boolean;
  generic?: boolean;
  goals?: boolean;
  instantLogin?: boolean;
  /** DEPRECATED - use nativeMobile instead. This field is mapped to nativeMobile for now. */
ios?: boolean;
  nativeMobile?: boolean;
  /** Start year for the school year. */
schoolYearStart?: number;
  secureSync?: boolean;
};
    
    type CollectionSchema = {
  extensionFieldsStatus?: ("none" | "required" | "optional");
  fieldSchemas?: FieldSchema[];
  name?: string;
  optInStatus?: ("in" | "out" | "clever_required");
};
    
    type CreateAdminParams = {
  createAdmin: CreateAdminRequest;
  adminID: string;
};
    
    type CreateAdminRequest = {
  admin?: Admin;
  password?: string;
};
    
    type CreateAppParams = {
  app?: App;
  appID: string;
};
    
    type CreateAppSchemaParams = {
  appID: string;
  /** Skip propagation to connection schemas */
skipPropagation?: boolean;
  /** Update data warnings when app schema changes */
updateDataRules?: boolean;
};
    
    type CreateLibraryResourceParams = {
  shortname: string;
  libraryResource: CreateLibraryResourceRequest;
};
    
    type CreateLibraryResourceRequest = {
  cleverID?: string;
};
    
    type CreatePlatformParams = {
  appID: string;
  request: CreatePlatformRequest;
};
    
    /** Used for the creation of a platform */
type CreatePlatformRequest = {
  disabled?: boolean;
  name?: string;
  redirectUris?: RedirectUri[];
};
    
    type CreateSetupStepParams = {
  appID: string;
  setupStep?: SetupStep;
};
    
    type CreateVerificationCodeParams = {
  duration: number;
  adminID: string;
};
    
    /** DataRule represents a rule on district configured by an app. The results of the rule are displayed on the data warnings page. Requirements for Rule Validity: - id is auto-generated by the app service when a rule is first created - simple rules apply validation on a document by document basis for a single field. They must have field set and can have field_required, max_length and regex set. - uniqueness rules validate that a field is unique across the entire data set. They only have a field set - custom rules exist for a College Board special case and don't have any of the other fields set. Severity of Rule Type: - warn indicates data may be not correct, but failing records are passed as is to app. - remove indicates app will not work for failing records and they will not be passed on to the app. */
type DataRule = {
  collection?: ("students" | "teachers" | "schools" | "sections" | "studentcontacts" | "schooladmins" | "courses" | "contacts" | "terms");
  description?: string;
  field?: string;
  fieldRequired?: boolean;
  groupByFields?: string[];
  id?: string;
  maxLength?: number;
  regex?: string;
  ruleType?: ("simple" | "uniqueness" | "custom" | "custom_dob_valid" | "custom_dob_range" | "dualenrollments" | "hasenrollments" | "uniqueenrollments" | "hasstudentcontacts" | "emptyschools" | "hascontacts" | "unlinkedadmins" | "teacher_staff_uniqueness" | "sectionhasschool");
  severity?: ("warn" | "remove");
};
    
    type DeleteAppSchemaParams = {
  appID: string;
  /** Delete field setting-style data warnings when app schema is deleted */
deleteDataRules?: boolean;
};
    
    type DeletePlatformParams = {
  appID: string;
  clientID: string;
};
    
    /** Error Codes */
type ErrorCode = ("Unknown" | "AdminAlreadyExists" | "AppNotFound" | "DuplicateRedirectURI" | "DuplicateShortname" | "EmailExistsAsDistrictUser" | "InvalidAccessTier" | "InvalidCategory" | "InvalidEmail" | "InvalidIconURL" | "InvalidID" | "InvalidLogoURL" | "InvalidPassword" | "InvalidRedirectURI" | "InvalidScope" | "InvalidVersionSettings" | "InvalidWebsiteURL" | "LastAdminPermission" | "OnboardingNotFound" | "ScopeNotAllowedForAccessTier" | "InvalidShortname" | "InvalidAppSchema");
    
    type FieldSchema = {
  accessType?: ("required" | "optional" | "available" | "clever_required" | "clever_generated");
  name?: string;
};
    
    type Forbidden = {
  message?: string;
};
    
    type GetAdminsParams = {
  email?: string;
  password?: string;
};
    
    type GetAllLibraryResourcesParams = {
  category?: string;
  includeDevApps?: boolean;
  includeLinks?: boolean;
};
    
    type GetAppsParams = {
  ids?: string[];
  clientId?: string;
  clientSecret?: string;
  shortname?: string;
  businessToken?: string;
  tags?: string[];
  skipTags?: string[];
};
    
    type GetCertificationsParams = {
  appID: string;
  schoolYearStart: number;
};
    
    type GetGuideConfigParams = {
  appID: string;
  adminID: string;
  guideID: string;
};
    
    type GetLibraryResourceByShortnameParams = {
  shortname: string;
  includeDevApps?: boolean;
  includeLinks?: boolean;
};
    
    type GetPermissionsForAdminParams = {
  adminID: string;
  appID: string;
};
    
    type GetPlatformsParams = {
  appIds?: string[];
  name?: string;
};
    
    /** Set of valid permissions. */
type GetValidPermissionsResponse = ValidPermission[];
    
    type GuideConfig = {
  guideData?: GuideData;
  stage?: string;
};
    
    type GuideData = {
  [key: string]: {
  
};
};
    
    type Image = {
  icon?: string;
  logo?: string;
};
    
    type InsightsMetadata = {
  enabled?: boolean;
  usageUrls?: string[];
};
    
    type KnownHost = {
  appId?: string;
  connectType?: ("ssh-rsa" | "ecdsa-sha2-nistp256");
  ipAddress?: string;
  key?: string;
  url?: string;
};
    
    type LibraryIntegrationType = ("link" | "partneredApplication");
    
    /** A resource to show in the Clever Library */
type LibraryResource = {
  appOnboarding?: AppOnboarding;
  bannerIconURL?: string;
  categories?: ResourceCategories;
  certifications?: ResourceCertifications;
  cleverID?: string;
  covidPromotion?: string;
  description?: string;
  devices?: ResourceDevices;
  featured?: number;
  features?: ResourceFeatures;
  gradeLevels?: ResourceGradeLevels;
  iconURL?: string;
  id?: string;
  insights?: InsightsMetadata;
  installCount?: number;
  libraryIntegration?: LibraryIntegrationType;
  marketingCollateralStatus?: MarketingCollateralStatus;
  name?: string;
  pricing?: ResourcePricing;
  privacyPolicyURL?: string;
  searchTags?: string[];
  shortname?: string;
  showInLibrary?: boolean;
  stripeUserID?: string;
  tagline?: string;
  teacherSetupGuide?: TeacherSetupGuide;
  teacherSetupGuideStatus?: TeacherSetupGuideStatus;
  tiers?: Tiers;
  tosURL?: string;
  universalDSASignatory?: boolean;
  url?: string;
};
    
    type LibraryResources = LibraryResource[];
    
    type LinkAppAdminParams = {
  appID: string;
  adminID: string;
  permissions: PermissionList;
};
    
    type Managers = {
  accountManager?: string;
  relationshipManager?: string;
};
    
    type MarketingCollateralStatus = ("notStarted" | "started" | "completed");
    
    /** A data structure used to track an app's progress through the self-serve onboarding flow */
type Onboarding = {
  events?: OnboardingEvent[];
  instantLogin?: OnboardingEvent[];
  ios?: OnboardingEvent[];
  secureSync?: OnboardingEvent[];
  secureSyncLite?: OnboardingEvent[];
};
    
    /** An event representing one of the following: an app starting its integration, an app submitting its integration for review, Clever approving an app's integration, or Clever rejecting an app's integration */
type OnboardingEvent = {
  eventType?: OnboardingEventType;
  survey?: OnboardingSurveyItem[];
  timestamp?: string;
};
    
    type OnboardingEventType = ("started" | "submitted" | "approved" | "rejected");
    
    type OnboardingIntegrationType = ("instantLogin" | "ios" | "secureSyncLite" | "secureSync" | "events");
    
    /** An app's response to a single onboarding survey question */
type OnboardingSurveyItem = {
  question?: string;
  response?: string;
};
    
    type OverrideConfigParams = {
  srcAppID: string;
  destAppID: string;
};
    
    type PatchAdminRequest = {
  email?: string;
  name?: string;
  password?: string;
  phone?: string;
  twoFactor?: TwoFactor;
};
    
    /** An update to an app */
type PatchAppRequest = {
  altNames?: string[];
  categories?: string[];
  customerSolutionsNotes?: string;
  defaultAccessTier?: string;
  description?: string;
  edsurgeUrl?: string;
  image?: Image;
  inviteMessage?: string;
  iosRedirectUri?: string;
  name?: string;
  nonOwnerImpersonationDisabled?: boolean;
  redirectUris?: string[];
  requiredScopes?: string[];
  salesContact?: string;
  shortname?: string;
  supportContact?: string;
  supportContactName?: string;
  supportedUserTypes?: UserType[];
  tags?: string[];
  versions?: VersionSettings;
  websiteUrl?: string;
};
    
    type PatchLibraryResourceRequest = {
  appOnboarding?: AppOnboarding;
  bannerIconURL?: string;
  categories?: ResourceCategories;
  certifications?: ResourceCertifications;
  cleverID?: string;
  covidPromotion?: string;
  description?: string;
  devices?: ResourceDevices;
  featured?: number;
  features?: ResourceFeatures;
  gradeLevels?: ResourceGradeLevels;
  iconURL?: string;
  insights?: InsightsMetadata;
  installCount?: number;
  libraryIntegration?: LibraryIntegrationType;
  marketingCollateralStatus?: MarketingCollateralStatus;
  name?: string;
  pricing?: ResourcePricing;
  privacyPolicyURL?: string;
  searchTags?: string[];
  shortname?: string;
  showInLibrary?: boolean;
  stripeUserID?: string;
  tagline?: string;
  teacherSetupGuide?: TeacherSetupGuide;
  teacherSetupGuideStatus?: TeacherSetupGuideStatus;
  tosURL?: string;
  universalDSASignatory?: boolean;
  url?: string;
};
    
    /** Used for an update to a platform */
type PatchPlatformRequest = {
  disabled?: boolean;
  redirectUris?: RedirectUri[];
};
    
    /** A permission for an app admin. */
type Permission = {
  action?: ("view" | "edit" | "none");
  resource?: ("admin" | "disconnect_district" | "district_requests" | "data_tools" | "app_settings" | "app_secrets" | "app_filters");
};
    
    type PermissionList = Permission[];
    
    /** Platforms are a way of distinguishing the different software clients of an app (e.g. iOS app, Desktop App, etc.) */
type Platform = {
  appId?: string;
  clientId?: string;
  disabled?: boolean;
  name?: ("ios" | "android" | "desktop");
  redirectUris?: RedirectUri[];
};
    
    /** Redirect URIs specify where codes and tokens should be sent in an OAuth flow. */
type RedirectUri = {
  /** Tags describe the purpose of the redirect URI, to distinguish it from other redirect URIs for the platform. */
tags?: string[];
  uri?: string;
};
    
    type RelatedApp = {
  cleverID?: string;
};
    
    type RelatedApps = {
  analytics?: string;
  connector?: string;
  districtPartner?: string;
  library?: string;
  saml?: string;
};
    
    type RelatedAppsV2 = {
  analytics?: string;
  districtPartners?: string[];
  library?: string[];
  saml?: string[];
  savedPasswords?: string[];
};
    
    type ResourceCategories = {
  assessment?: boolean;
  authoring?: boolean;
  classroom_management?: boolean;
  english?: boolean;
  math?: boolean;
  other?: boolean;
  presentation?: boolean;
  science?: boolean;
  social_studies?: boolean;
  technology?: boolean;
};
    
    type ResourceCertifications = {
  commonsensemedia_reviewed?: boolean;
  ikeepsafe_coppa?: boolean;
  ikeepsafe_ferpa?: boolean;
  student_privacy_pledge?: boolean;
};
    
    type ResourceDevices = {
  android?: boolean;
  browser?: boolean;
  iOS?: boolean;
  macOS?: boolean;
  windows?: boolean;
};
    
    type ResourceFeature = {
  description?: string;
  screenshotURL?: string;
};
    
    type ResourceFeatures = ResourceFeature[];
    
    type ResourceGradeLevels = {
  k2?: boolean;
  nineTwelve?: boolean;
  sixEight?: boolean;
  staff?: boolean;
  threeFive?: boolean;
};
    
    type ResourcePricing = ("free" | "freemium" | "free_trial" | "subscription");
    
    type SearchLibraryResourceParams = {
  searchTerm: string;
  showInLibraryOnly?: boolean;
  includeLinks?: boolean;
};
    
    type SecretConfig = {
  businessToken?: string;
  currentSecret?: string;
  oldSecret?: string;
  oldSecretExpires?: string;
};
    
    type SetAppSchemaParams = {
  appID: string;
  /** Skip propagation to connection schemas */
skipPropagation?: boolean;
  /** Update data warnings when app schema changes */
updateDataRules?: boolean;
  appSchema?: AppSchema;
};
    
    type SetCertificationsParams = {
  appID: string;
  schoolYearStart: number;
  certifications: SetCertificationsRequest;
};
    
    /** Integration certifications for an app, for a specific school year. */
type SetCertificationsRequest = {
  /** DEPRECATED - use certS3PathIntegration and certS3PathTradeshow instead. This field is mapped to certS3PathTradeshow. */
certS3Path?: string;
  certS3PathIntegration?: string;
  certS3PathTradeshow?: string;
  events?: boolean;
  generic?: boolean;
  goals?: boolean;
  instantLogin?: boolean;
  /** DEPRECATED - use nativeMobile instead. This field is mapped to nativeMobile for now. */
ios?: boolean;
  nativeMobile?: boolean;
  secureSync?: boolean;
};
    
    type SetDataRules = {
  rules?: DataRule[];
};
    
    type SetDataRulesParams = {
  appID: string;
  rules?: SetDataRulesRequest;
};
    
    type SetDataRulesRequest = DataRule[];
    
    type SetGuideConfigParams = {
  appID: string;
  adminID: string;
  guideID: string;
  guideConfig: GuideConfig;
};
    
    type SetRecommendedSharingParams = {
  appID: string;
  recommendations?: SharingRecommendations;
};
    
    type SetupStep = {
  app_id?: string;
  content?: string;
  creator?: string;
  edit_date?: string;
  id?: string;
  title?: string;
};
    
    type SharingConstraint = {
  collection?: ("sections" | "students" | "teachers" | "schooladmins");
  field?: string;
  operator?: ("Contains" | "EqualsAnyOf" | "NotContains" | "NotEqualsAnyOf" | "StartsWith");
  values?: string[];
};
    
    /** The sharing recommendations for an application. Should match the schema used in District Authorizations */
type SharingRecommendations = {
  rules?: SharingRule[];
  sharingType?: ("none" | "district" | "schools" | "sections" | "rules");
};
    
    type SharingRule = {
  collection?: ("sections" | "students" | "teachers" | "schooladmins");
  constraints?: SharingConstraint[];
};
    
    type TeacherSetupGuide = {
  steps?: TeacherSetupGuideStep[];
};
    
    type TeacherSetupGuideStatus = ("notStarted" | "started" | "completed");
    
    type TeacherSetupGuideStep = {
  content?: string;
  header?: string;
  screenshotURL?: string;
  type?: ("accountSyncing" | "howToUse");
};
    
    type Tier = {
  duration?: number;
  limit?: number;
  limitType?: string;
  name?: string;
  price?: number;
  studentLimit?: number;
  tierFeatures?: string[];
};
    
    type Tiers = Tier[];
    
    /** An app being tracked for P4A */
type TrackableApp = {
  _name?: string;
  id: string;
  isAppGroup?: boolean;
  relatedApps?: RelatedApps;
  relatedAppsV2?: RelatedAppsV2;
  trackingWhitelist?: TrackingWhitelist;
};
    
    type TrackableApps = {
  [key: string]: TrackableApp;
};
    
    type TrackingWhitelist = {
  enableTracking?: boolean;
  hideInUI?: boolean;
  urlMatchers?: string[];
};
    
    type TwoFactor = {
  authType?: string;
  authyId?: string;
  confirmed?: boolean;
  rememberMeKey?: string;
};
    
    type UnlinkAppAdminParams = {
  appID: string;
  adminID: string;
};
    
    type UnprocessableEntity = {
  message?: string;
};
    
    type UpdateAdminParams = {
  adminID: string;
  admin: PatchAdminRequest;
};
    
    type UpdateAppIconParams = {
  appID: string;
  app: UpdateAppIconRequest;
};
    
    /** The new S3 image to use for an app icon */
type UpdateAppIconRequest = {
  newIcon?: string;
};
    
    type UpdateAppParams = {
  appID: string;
  /**
	If scopes change, then the app schema will be updated. This flag will propagate app schema updates to all connection schemas as well

	*/
	withSchemaPropagation?: boolean;
  app: PatchAppRequest;
};
    
    type UpdateLibraryResourceByShortnameParams = {
  shortname: string;
  libraryResource: PatchLibraryResourceRequest;
};
    
    type UpdateOnboardingParams = {
  appID: string;
  update: UpdateOnboardingRequest;
};
    
    /** An update to an app's onboarding progress */
type UpdateOnboardingRequest = {
  eventType?: OnboardingEventType;
  integrationType?: OnboardingIntegrationType;
  survey?: OnboardingSurveyItem[];
};
    
    type UpdatePlatformParams = {
  appID: string;
  clientID: string;
  request: PatchPlatformRequest;
};
    
    type UsageUrls = string[];
    
    /** Valid user types */
type UserType = ("teachers" | "students" | "district_admins" | "school_admins" | "clever_users");
    
    type ValidPermission = {
  description?: string;
  value?: Permission;
};
    
    type VerificationCodeResponse = string;
    
    type VerifyAdminEmailParams = {
  adminID: string;
  request: VerifyAdminEmailRequest;
};
    
    type VerifyAdminEmailRequest = {
  email?: string;
};
    
    type VerifyAppAdminParams = {
  appID: string;
  adminID: string;
  verified: boolean;
};
    
    type VerifyCodeParams = {
  code: string;
  invalidate?: boolean;
  adminID: string;
};
    
    /** Valid versions */
type Version = ("v1.1" | "v1.2" | "v2.0" | "v2.1");
    
    /** Version settings for an application */
type VersionSettings = {
  allowed?: Version[];
  default?: Version;
};
    
  }
}

export = AppService;
