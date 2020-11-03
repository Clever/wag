<a name="module_app-service"></a>

## app-service
app-service client library.


* [app-service](#module_app-service)
    * [AppService](#exp_module_app-service--AppService) ⏏
        * [new AppService(options)](#new_module_app-service--AppService_new)
        * _instance_
            * [.healthCheck([options], [cb])](#module_app-service--AppService+healthCheck) ⇒ <code>Promise</code>
            * [.getAdmins(params, [options], [cb])](#module_app-service--AppService+getAdmins) ⇒ <code>Promise</code>
            * [.deleteAdmin(adminID, [options], [cb])](#module_app-service--AppService+deleteAdmin) ⇒ <code>Promise</code>
            * [.getAdminByID(adminID, [options], [cb])](#module_app-service--AppService+getAdminByID) ⇒ <code>Promise</code>
            * [.updateAdmin(params, [options], [cb])](#module_app-service--AppService+updateAdmin) ⇒ <code>Promise</code>
            * [.createAdmin(params, [options], [cb])](#module_app-service--AppService+createAdmin) ⇒ <code>Promise</code>
            * [.verifyCode(params, [options], [cb])](#module_app-service--AppService+verifyCode) ⇒ <code>Promise</code>
            * [.createVerificationCode(params, [options], [cb])](#module_app-service--AppService+createVerificationCode) ⇒ <code>Promise</code>
            * [.verifyAdminEmail(params, [options], [cb])](#module_app-service--AppService+verifyAdminEmail) ⇒ <code>Promise</code>
            * [.getAllAnalyticsApps([options], [cb])](#module_app-service--AppService+getAllAnalyticsApps) ⇒ <code>Promise</code>
            * [.getAnalyticsAppByShortname(shortname, [options], [cb])](#module_app-service--AppService+getAnalyticsAppByShortname) ⇒ <code>Promise</code>
            * [.getAllTrackableApps([options], [cb])](#module_app-service--AppService+getAllTrackableApps) ⇒ <code>Promise</code>
            * [.getAnalyticsUsageUrls([options], [cb])](#module_app-service--AppService+getAnalyticsUsageUrls) ⇒ <code>Promise</code>
            * [.getAllUsageUrls([options], [cb])](#module_app-service--AppService+getAllUsageUrls) ⇒ <code>Promise</code>
            * [.getApps(params, [options], [cb])](#module_app-service--AppService+getApps) ⇒ <code>Promise</code>
            * [.deleteApp(appID, [options], [cb])](#module_app-service--AppService+deleteApp) ⇒ <code>Promise</code>
            * [.getAppByID(appID, [options], [cb])](#module_app-service--AppService+getAppByID) ⇒ <code>Promise</code>
            * [.updateApp(params, [options], [cb])](#module_app-service--AppService+updateApp) ⇒ <code>Promise</code>
            * [.createApp(params, [options], [cb])](#module_app-service--AppService+createApp) ⇒ <code>Promise</code>
            * [.getAdminsForApp(appID, [options], [cb])](#module_app-service--AppService+getAdminsForApp) ⇒ <code>Promise</code>
            * [.unlinkAppAdmin(params, [options], [cb])](#module_app-service--AppService+unlinkAppAdmin) ⇒ <code>Promise</code>
            * [.linkAppAdmin(params, [options], [cb])](#module_app-service--AppService+linkAppAdmin) ⇒ <code>Promise</code>
            * [.getGuideConfig(params, [options], [cb])](#module_app-service--AppService+getGuideConfig) ⇒ <code>Promise</code>
            * [.setGuideConfig(params, [options], [cb])](#module_app-service--AppService+setGuideConfig) ⇒ <code>Promise</code>
            * [.getPermissionsForAdmin(params, [options], [cb])](#module_app-service--AppService+getPermissionsForAdmin) ⇒ <code>Promise</code>
            * [.verifyAppAdmin(params, [options], [cb])](#module_app-service--AppService+verifyAppAdmin) ⇒ <code>Promise</code>
            * [.generateNewBusinessToken(appID, [options], [cb])](#module_app-service--AppService+generateNewBusinessToken) ⇒ <code>Promise</code>
            * [.getCertifications(params, [options], [cb])](#module_app-service--AppService+getCertifications) ⇒ <code>Promise</code>
            * [.setCertifications(params, [options], [cb])](#module_app-service--AppService+setCertifications) ⇒ <code>Promise</code>
            * [.getSetupStep(appID, [options], [cb])](#module_app-service--AppService+getSetupStep) ⇒ <code>Promise</code>
            * [.createSetupStep(params, [options], [cb])](#module_app-service--AppService+createSetupStep) ⇒ <code>Promise</code>
            * [.getDataRules(appID, [options], [cb])](#module_app-service--AppService+getDataRules) ⇒ <code>Promise</code>
            * [.setDataRules(params, [options], [cb])](#module_app-service--AppService+setDataRules) ⇒ <code>Promise</code>
            * [.getManagers(appID, [options], [cb])](#module_app-service--AppService+getManagers) ⇒ <code>Promise</code>
            * [.getOnboarding(appID, [options], [cb])](#module_app-service--AppService+getOnboarding) ⇒ <code>Promise</code>
            * [.updateOnboarding(params, [options], [cb])](#module_app-service--AppService+updateOnboarding) ⇒ <code>Promise</code>
            * [.initializeOnboarding(appID, [options], [cb])](#module_app-service--AppService+initializeOnboarding) ⇒ <code>Promise</code>
            * [.deletePlatform(params, [options], [cb])](#module_app-service--AppService+deletePlatform) ⇒ <code>Promise</code>
            * [.updatePlatform(params, [options], [cb])](#module_app-service--AppService+updatePlatform) ⇒ <code>Promise</code>
            * [.getPlatformsByAppID(appID, [options], [cb])](#module_app-service--AppService+getPlatformsByAppID) ⇒ <code>Promise</code>
            * [.createPlatform(params, [options], [cb])](#module_app-service--AppService+createPlatform) ⇒ <code>Promise</code>
            * [.deleteAppSchema(params, [options], [cb])](#module_app-service--AppService+deleteAppSchema) ⇒ <code>Promise</code>
            * [.getAppSchema(appID, [options], [cb])](#module_app-service--AppService+getAppSchema) ⇒ <code>Promise</code>
            * [.createAppSchema(params, [options], [cb])](#module_app-service--AppService+createAppSchema) ⇒ <code>Promise</code>
            * [.setAppSchema(params, [options], [cb])](#module_app-service--AppService+setAppSchema) ⇒ <code>Promise</code>
            * [.getSecrets(appID, [options], [cb])](#module_app-service--AppService+getSecrets) ⇒ <code>Promise</code>
            * [.revokeOldClientSecret(appID, [options], [cb])](#module_app-service--AppService+revokeOldClientSecret) ⇒ <code>Promise</code>
            * [.generateNewClientSecret(appID, [options], [cb])](#module_app-service--AppService+generateNewClientSecret) ⇒ <code>Promise</code>
            * [.resetClientSecret(appID, [options], [cb])](#module_app-service--AppService+resetClientSecret) ⇒ <code>Promise</code>
            * [.getRecommendedSharing(appID, [options], [cb])](#module_app-service--AppService+getRecommendedSharing) ⇒ <code>Promise</code>
            * [.setRecommendedSharing(params, [options], [cb])](#module_app-service--AppService+setRecommendedSharing) ⇒ <code>Promise</code>
            * [.updateAppIcon(params, [options], [cb])](#module_app-service--AppService+updateAppIcon) ⇒ <code>Promise</code>
            * [.getAllCategories([options], [cb])](#module_app-service--AppService+getAllCategories) ⇒ <code>Promise</code>
            * [.getKnownHosts([options], [cb])](#module_app-service--AppService+getKnownHosts) ⇒ <code>Promise</code>
            * [.getAllLibraryResources(params, [options], [cb])](#module_app-service--AppService+getAllLibraryResources) ⇒ <code>Promise</code>
            * [.searchLibraryResource(params, [options], [cb])](#module_app-service--AppService+searchLibraryResource) ⇒ <code>Promise</code>
            * [.getLibraryResourceByShortname(params, [options], [cb])](#module_app-service--AppService+getLibraryResourceByShortname) ⇒ <code>Promise</code>
            * [.updateLibraryResourceByShortname(params, [options], [cb])](#module_app-service--AppService+updateLibraryResourceByShortname) ⇒ <code>Promise</code>
            * [.createLibraryResource(params, [options], [cb])](#module_app-service--AppService+createLibraryResource) ⇒ <code>Promise</code>
            * [.deleteLibraryResourceLink(shortname, [options], [cb])](#module_app-service--AppService+deleteLibraryResourceLink) ⇒ <code>Promise</code>
            * [.getValidPermissions([options], [cb])](#module_app-service--AppService+getValidPermissions) ⇒ <code>Promise</code>
            * [.getPlatforms(params, [options], [cb])](#module_app-service--AppService+getPlatforms) ⇒ <code>Promise</code>
            * [.getPlatformByClientID(clientID, [options], [cb])](#module_app-service--AppService+getPlatformByClientID) ⇒ <code>Promise</code>
            * [.getAppsForAdmin(adminID, [options], [cb])](#module_app-service--AppService+getAppsForAdmin) ⇒ <code>Promise</code>
            * [.overrideConfig(params, [options], [cb])](#module_app-service--AppService+overrideConfig) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_app-service--AppService.RetryPolicies)
                * [.Exponential](#module_app-service--AppService.RetryPolicies.Exponential)
                * [.Single](#module_app-service--AppService.RetryPolicies.Single)
                * [.None](#module_app-service--AppService.RetryPolicies.None)
            * [.Errors](#module_app-service--AppService.Errors)
                * [.BadRequest](#module_app-service--AppService.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_app-service--AppService.Errors.InternalError) ⇐ <code>Error</code>
                * [.NotFound](#module_app-service--AppService.Errors.NotFound) ⇐ <code>Error</code>
                * [.Forbidden](#module_app-service--AppService.Errors.Forbidden) ⇐ <code>Error</code>
                * [.UnprocessableEntity](#module_app-service--AppService.Errors.UnprocessableEntity) ⇐ <code>Error</code>
            * [.DefaultCircuitOptions](#module_app-service--AppService.DefaultCircuitOptions)

<a name="exp_module_app-service--AppService"></a>

### AppService ⏏
app-service client

**Kind**: Exported class  
<a name="new_module_app-service--AppService_new"></a>

#### new AppService(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| [options.address] | <code>string</code> |  | URL where the server is located. Must provide this or the discovery argument |
| [options.discovery] | <code>bool</code> |  | Use clever-discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. |
| [options.keepalive] | <code>bool</code> |  | Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;app-service-wagclient&quot;)</code> | The Kayvee logger to use in the client. |
| [options.circuit] | <code>Object</code> |  | Options for constructing the client's circuit breaker. |
| [options.circuit.forceClosed] | <code>bool</code> |  | When set to true the circuit will always be closed. Default: true. |
| [options.circuit.maxConcurrentRequests] | <code>number</code> |  | the maximum number of concurrent requests the client can make at the same time. Default: 100. |
| [options.circuit.requestVolumeThreshold] | <code>number</code> |  | The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. |
| [options.circuit.sleepWindow] | <code>number</code> |  | how long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. |
| [options.circuit.errorPercentThreshold] | <code>number</code> |  | the threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. |

<a name="module_app-service--AppService+healthCheck"></a>

#### appService.healthCheck([options], [cb]) ⇒ <code>Promise</code>
Checks if the service is healthy

**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAdmins"></a>

#### appService.getAdmins(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.email] | <code>string</code> |  |
| [params.password] | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+deleteAdmin"></a>

#### appService.deleteAdmin(adminID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAdminByID"></a>

#### appService.getAdminByID(adminID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updateAdmin"></a>

#### appService.updateAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.adminID | <code>string</code> |  |
| params.admin |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createAdmin"></a>

#### appService.createAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.createAdmin |  |  |
| params.adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+verifyCode"></a>

#### appService.verifyCode(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.code | <code>string</code> |  |
| [params.invalidate] | <code>boolean</code> |  |
| params.adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createVerificationCode"></a>

#### appService.createVerificationCode(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.duration | <code>number</code> |  |
| params.adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+verifyAdminEmail"></a>

#### appService.verifyAdminEmail(params, [options], [cb]) ⇒ <code>Promise</code>
set the verified email of an admin

**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.adminID | <code>string</code> |  |
| params.request |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAllAnalyticsApps"></a>

#### appService.getAllAnalyticsApps([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAnalyticsAppByShortname"></a>

#### appService.getAnalyticsAppByShortname(shortname, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| shortname | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAllTrackableApps"></a>

#### appService.getAllTrackableApps([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAnalyticsUsageUrls"></a>

#### appService.getAnalyticsUsageUrls([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAllUsageUrls"></a>

#### appService.getAllUsageUrls([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getApps"></a>

#### appService.getApps(params, [options], [cb]) ⇒ <code>Promise</code>
The server takes in the intersection of input parameters

**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.ids] | <code>Array.&lt;string&gt;</code> |  |
| [params.clientId] | <code>string</code> |  |
| [params.clientSecret] | <code>string</code> |  |
| [params.shortname] | <code>string</code> |  |
| [params.businessToken] | <code>string</code> |  |
| [params.tags] | <code>Array.&lt;string&gt;</code> |  |
| [params.skipTags] | <code>Array.&lt;string&gt;</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+deleteApp"></a>

#### appService.deleteApp(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAppByID"></a>

#### appService.getAppByID(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updateApp"></a>

#### appService.updateApp(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.withSchemaPropagation] | <code>boolean</code> | If scopes change, then the app schema will be updated. This flag will propagate app schema updates to all connection schemas as well |
| params.app |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createApp"></a>

#### appService.createApp(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.app] |  |  |
| params.appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAdminsForApp"></a>

#### appService.getAdminsForApp(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+unlinkAppAdmin"></a>

#### appService.unlinkAppAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>Forbidden</code>](#module_app-service--AppService.Errors.Forbidden)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+linkAppAdmin"></a>

#### appService.linkAppAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>Forbidden</code>](#module_app-service--AppService.Errors.Forbidden)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.adminID | <code>string</code> |  |
| params.permissions |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getGuideConfig"></a>

#### appService.getGuideConfig(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>Forbidden</code>](#module_app-service--AppService.Errors.Forbidden)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.adminID | <code>string</code> |  |
| params.guideID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+setGuideConfig"></a>

#### appService.setGuideConfig(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>Forbidden</code>](#module_app-service--AppService.Errors.Forbidden)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.adminID | <code>string</code> |  |
| params.guideID | <code>string</code> |  |
| params.guideConfig |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getPermissionsForAdmin"></a>

#### appService.getPermissionsForAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.adminID | <code>string</code> |  |
| params.appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+verifyAppAdmin"></a>

#### appService.verifyAppAdmin(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>Forbidden</code>](#module_app-service--AppService.Errors.Forbidden)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.adminID | <code>string</code> |  |
| params.verified | <code>boolean</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+generateNewBusinessToken"></a>

#### appService.generateNewBusinessToken(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getCertifications"></a>

#### appService.getCertifications(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.schoolYearStart | <code>number</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+setCertifications"></a>

#### appService.setCertifications(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.schoolYearStart | <code>number</code> |  |
| params.certifications |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getSetupStep"></a>

#### appService.getSetupStep(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createSetupStep"></a>

#### appService.createSetupStep(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.setupStep] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getDataRules"></a>

#### appService.getDataRules(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+setDataRules"></a>

#### appService.setDataRules(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.rules] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getManagers"></a>

#### appService.getManagers(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getOnboarding"></a>

#### appService.getOnboarding(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updateOnboarding"></a>

#### appService.updateOnboarding(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.update |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+initializeOnboarding"></a>

#### appService.initializeOnboarding(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+deletePlatform"></a>

#### appService.deletePlatform(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.clientID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updatePlatform"></a>

#### appService.updatePlatform(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.clientID | <code>string</code> |  |
| params.request |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getPlatformsByAppID"></a>

#### appService.getPlatformsByAppID(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createPlatform"></a>

#### appService.createPlatform(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.request |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+deleteAppSchema"></a>

#### appService.deleteAppSchema(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.deleteDataRules] | <code>boolean</code> | Delete field setting-style data warnings when app schema is deleted |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAppSchema"></a>

#### appService.getAppSchema(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createAppSchema"></a>

#### appService.createAppSchema(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.skipPropagation] | <code>boolean</code> | Skip propagation to connection schemas |
| [params.updateDataRules] | <code>boolean</code> | Update data warnings when app schema changes |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+setAppSchema"></a>

#### appService.setAppSchema(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.skipPropagation] | <code>boolean</code> | Skip propagation to connection schemas |
| [params.updateDataRules] | <code>boolean</code> | Update data warnings when app schema changes |
| [params.appSchema] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getSecrets"></a>

#### appService.getSecrets(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+revokeOldClientSecret"></a>

#### appService.revokeOldClientSecret(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+generateNewClientSecret"></a>

#### appService.generateNewClientSecret(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+resetClientSecret"></a>

#### appService.resetClientSecret(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getRecommendedSharing"></a>

#### appService.getRecommendedSharing(appID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| appID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+setRecommendedSharing"></a>

#### appService.setRecommendedSharing(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| [params.recommendations] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updateAppIcon"></a>

#### appService.updateAppIcon(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>UnprocessableEntity</code>](#module_app-service--AppService.Errors.UnprocessableEntity)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.appID | <code>string</code> |  |
| params.app |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAllCategories"></a>

#### appService.getAllCategories([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getKnownHosts"></a>

#### appService.getKnownHosts([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAllLibraryResources"></a>

#### appService.getAllLibraryResources(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.category] | <code>string</code> |  |
| [params.includeDevApps] | <code>boolean</code> |  |
| [params.includeLinks] | <code>boolean</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+searchLibraryResource"></a>

#### appService.searchLibraryResource(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.searchTerm | <code>string</code> |  |
| [params.showInLibraryOnly] | <code>boolean</code> |  |
| [params.includeLinks] | <code>boolean</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getLibraryResourceByShortname"></a>

#### appService.getLibraryResourceByShortname(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.shortname | <code>string</code> |  |
| [params.includeDevApps] | <code>boolean</code> |  |
| [params.includeLinks] | <code>boolean</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+updateLibraryResourceByShortname"></a>

#### appService.updateLibraryResourceByShortname(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.shortname | <code>string</code> |  |
| params.libraryResource |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+createLibraryResource"></a>

#### appService.createLibraryResource(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.shortname | <code>string</code> |  |
| params.libraryResource |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+deleteLibraryResourceLink"></a>

#### appService.deleteLibraryResourceLink(shortname, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| shortname | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getValidPermissions"></a>

#### appService.getValidPermissions([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getPlatforms"></a>

#### appService.getPlatforms(params, [options], [cb]) ⇒ <code>Promise</code>
The server takes in the intersection of input parameters

**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.appIds] | <code>Array.&lt;string&gt;</code> |  |
| [params.name] | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getPlatformByClientID"></a>

#### appService.getPlatformByClientID(clientID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| clientID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+getAppsForAdmin"></a>

#### appService.getAppsForAdmin(adminID, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| adminID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService+overrideConfig"></a>

#### appService.overrideConfig(params, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>AppService</code>](#exp_module_app-service--AppService)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_app-service--AppService.Errors.BadRequest)  
**Reject**: [<code>NotFound</code>](#module_app-service--AppService.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_app-service--AppService.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.srcAppID | <code>string</code> |  |
| params.destAppID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | [<code>Span</code>](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html) | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_app-service--AppService.RetryPolicies"></a>

#### AppService.RetryPolicies
Retry policies available to use.

**Kind**: static property of [<code>AppService</code>](#exp_module_app-service--AppService)  

* [.RetryPolicies](#module_app-service--AppService.RetryPolicies)
    * [.Exponential](#module_app-service--AppService.RetryPolicies.Exponential)
    * [.Single](#module_app-service--AppService.RetryPolicies.Single)
    * [.None](#module_app-service--AppService.RetryPolicies.None)

<a name="module_app-service--AppService.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies)  
<a name="module_app-service--AppService.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies)  
<a name="module_app-service--AppService.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_app-service--AppService.RetryPolicies)  
<a name="module_app-service--AppService.Errors"></a>

#### AppService.Errors
Errors returned by methods.

**Kind**: static property of [<code>AppService</code>](#exp_module_app-service--AppService)  

* [.Errors](#module_app-service--AppService.Errors)
    * [.BadRequest](#module_app-service--AppService.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_app-service--AppService.Errors.InternalError) ⇐ <code>Error</code>
    * [.NotFound](#module_app-service--AppService.Errors.NotFound) ⇐ <code>Error</code>
    * [.Forbidden](#module_app-service--AppService.Errors.Forbidden) ⇐ <code>Error</code>
    * [.UnprocessableEntity](#module_app-service--AppService.Errors.UnprocessableEntity) ⇐ <code>Error</code>

<a name="module_app-service--AppService.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of [<code>Errors</code>](#module_app-service--AppService.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code |  | 
| message | <code>string</code> | 

<a name="module_app-service--AppService.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of [<code>Errors</code>](#module_app-service--AppService.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_app-service--AppService.Errors.NotFound"></a>

##### Errors.NotFound ⇐ <code>Error</code>
NotFound

**Kind**: static class of [<code>Errors</code>](#module_app-service--AppService.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code |  | 
| message | <code>string</code> | 

<a name="module_app-service--AppService.Errors.Forbidden"></a>

##### Errors.Forbidden ⇐ <code>Error</code>
Forbidden

**Kind**: static class of [<code>Errors</code>](#module_app-service--AppService.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_app-service--AppService.Errors.UnprocessableEntity"></a>

##### Errors.UnprocessableEntity ⇐ <code>Error</code>
UnprocessableEntity

**Kind**: static class of [<code>Errors</code>](#module_app-service--AppService.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_app-service--AppService.DefaultCircuitOptions"></a>

#### AppService.DefaultCircuitOptions
Default circuit breaker options.

**Kind**: static constant of [<code>AppService</code>](#exp_module_app-service--AppService)  
