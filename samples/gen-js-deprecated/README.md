<a name="module_swagger-test"></a>

## swagger-test
swagger-test client library.


* [swagger-test](#module_swagger-test)
    * [SwaggerTest](#exp_module_swagger-test--SwaggerTest) ⏏
        * [new SwaggerTest(options)](#new_module_swagger-test--SwaggerTest_new)
        * [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
            * [.Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)
                * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.backoffs) ⇒ <code>Array.&lt;number&gt;</code>
                * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.retry) ⇒ <code>boolean</code>
            * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
                * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.None.backoffs)
                * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.None.retry)
        * [.Errors](#module_swagger-test--SwaggerTest.Errors)
            * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
            * [.NotFound](#module_swagger-test--SwaggerTest.Errors.NotFound) ⇐ <code>Error</code>
            * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="exp_module_swagger-test--SwaggerTest"></a>

### SwaggerTest ⏏
The main client object to instantiate.

**Kind**: Exported class  
<a name="new_module_swagger-test--SwaggerTest_new"></a>

#### new SwaggerTest(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| options.address | <code>string</code> |  | URL where the server is located. If not specified, the address will be discovered via @clever/discovery. |
| options.timeout | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. |
| [options.retryPolicy] | <code>Object</code> | <code>RetryPolicies.Default</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| options.retryPolicy.backoffs | <code>function</code> |  |  |
| options.retryPolicy.retry | <code>function</code> |  |  |

<a name="module_swagger-test--SwaggerTest.RetryPolicies"></a>

#### SwaggerTest.RetryPolicies
Retry policies available to use.

**Kind**: static property of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  

* [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
    * [.Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)
        * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.backoffs) ⇒ <code>Array.&lt;number&gt;</code>
        * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.retry) ⇒ <code>boolean</code>
    * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
        * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.None.backoffs)
        * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.None.retry)

<a name="module_swagger-test--SwaggerTest.RetryPolicies.Default"></a>

##### RetryPolicies.Default
The default retry policy will retry five times with an exponential backoff.

**Kind**: static constant of <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code>  

* [.Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)
    * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.backoffs) ⇒ <code>Array.&lt;number&gt;</code>
    * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.retry) ⇒ <code>boolean</code>

<a name="module_swagger-test--SwaggerTest.RetryPolicies.Default.backoffs"></a>

###### Default.backoffs() ⇒ <code>Array.&lt;number&gt;</code>
backoffs returns an array of five backoffs: 100ms, 200ms, 400ms, 800ms, and
1.6s. It adds a random 5% jitter to each backoff.

**Kind**: static method of <code>[Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)</code>  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.Default.retry"></a>

###### Default.retry() ⇒ <code>boolean</code>
retry will not retry a request if the HTTP client returns an error, if the
is a POST or PATCH, or if the status code is less than 500. It will retry
all other requests.

**Kind**: static method of <code>[Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)</code>  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of <code>[RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)</code>  

* [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
    * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.None.backoffs)
    * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.None.retry)

<a name="module_swagger-test--SwaggerTest.RetryPolicies.None.backoffs"></a>

###### None.backoffs()
returns an empty array

**Kind**: static method of <code>[None](#module_swagger-test--SwaggerTest.RetryPolicies.None)</code>  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.None.retry"></a>

###### None.retry()
returns false

**Kind**: static method of <code>[None](#module_swagger-test--SwaggerTest.RetryPolicies.None)</code>  
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

