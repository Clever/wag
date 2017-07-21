<a name="module_swagger-test"></a>

## swagger-test
swagger-test client library.


* [swagger-test](#module_swagger-test)
    * [SwaggerTest](#exp_module_swagger-test--SwaggerTest) ⏏
        * [new SwaggerTest(options)](#new_module_swagger-test--SwaggerTest_new)
        * [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
            * [.Exponential](#module_swagger-test--SwaggerTest.RetryPolicies.Exponential)
            * [.Single](#module_swagger-test--SwaggerTest.RetryPolicies.Single)
            * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
        * [.Errors](#module_swagger-test--SwaggerTest.Errors)
            * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
            * [.NotFound](#module_swagger-test--SwaggerTest.Errors.NotFound) ⇐ <code>Error</code>
            * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="exp_module_swagger-test--SwaggerTest"></a>

### SwaggerTest ⏏
swagger-test client

**Kind**: Exported class  
<a name="new_module_swagger-test--SwaggerTest_new"></a>

#### new SwaggerTest(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| [options.address] | <code>string</code> |  | URL where the server is located. Must provide this or the discovery argument |
| [options.discovery] | <code>bool</code> |  | Use clever-discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. |
| [options.retryPolicy] | <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code> | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;swagger-test-wagclient&quot;)</code> | The Kayvee  logger to use in the client. |

<a name="module_swagger-test--SwaggerTest.RetryPolicies"></a>

#### SwaggerTest.RetryPolicies
Retry policies available to use.

**Kind**: static property of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  

* [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
    * [.Exponential](#module_swagger-test--SwaggerTest.RetryPolicies.Exponential)
    * [.Single](#module_swagger-test--SwaggerTest.RetryPolicies.Single)
    * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)

<a name="module_swagger-test--SwaggerTest.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code>  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code>  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code>  
<a name="module_swagger-test--SwaggerTest.Errors"></a>

#### SwaggerTest.Errors
Errors returned by methods.

**Kind**: static property of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  

* [.Errors](#module_swagger-test--SwaggerTest.Errors)
    * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
    * [.NotFound](#module_swagger-test--SwaggerTest.Errors.NotFound) ⇐ <code>Error</code>
    * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_swagger-test--SwaggerTest.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of <code>[Errors](#module_swagger-test--SwaggerTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.NotFound"></a>

##### Errors.NotFound ⇐ <code>Error</code>
NotFound

**Kind**: static class of <code>[Errors](#module_swagger-test--SwaggerTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of <code>[Errors](#module_swagger-test--SwaggerTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

