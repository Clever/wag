<a name="module_swagger-test"></a>

## swagger-test
swagger-test client library.


* [swagger-test](#module_swagger-test)
    * [SwaggerTest](#exp_module_swagger-test--SwaggerTest) ⏏
        * [new SwaggerTest(options)](#new_module_swagger-test--SwaggerTest_new)
        * _instance_
            * [.close()](#module_swagger-test--SwaggerTest+close)
            * [.getBook(id, [options], [cb])](#module_swagger-test--SwaggerTest+getBook) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
                * [.Exponential](#module_swagger-test--SwaggerTest.RetryPolicies.Exponential)
                * [.Single](#module_swagger-test--SwaggerTest.RetryPolicies.Single)
                * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
            * [.Errors](#module_swagger-test--SwaggerTest.Errors)
                * [.ExtendedError](#module_swagger-test--SwaggerTest.Errors.ExtendedError) ⇐ <code>Error</code>
                * [.NotFound](#module_swagger-test--SwaggerTest.Errors.NotFound) ⇐ <code>Error</code>
                * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>
            * [.DefaultCircuitOptions](#module_swagger-test--SwaggerTest.DefaultCircuitOptions)

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
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. |
| [options.keepalive] | <code>bool</code> |  | Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;swagger-test-wagclient&quot;)</code> | The Kayvee logger to use in the client. |
| [options.circuit] | <code>Object</code> |  | Options for constructing the client's circuit breaker. |
| [options.circuit.forceClosed] | <code>bool</code> |  | When set to true the circuit will always be closed. Default: true. |
| [options.circuit.maxConcurrentRequests] | <code>number</code> |  | the maximum number of concurrent requests the client can make at the same time. Default: 100. |
| [options.circuit.requestVolumeThreshold] | <code>number</code> |  | The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. |
| [options.circuit.sleepWindow] | <code>number</code> |  | how long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. |
| [options.circuit.errorPercentThreshold] | <code>number</code> |  | the threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. |
| [options.asynclocalstore] | <code>object</code> |  | a request scoped async store |

<a name="module_swagger-test--SwaggerTest+close"></a>

#### swaggerTest.close()
Releases handles used in client

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
<a name="module_swagger-test--SwaggerTest+getBook"></a>

#### swaggerTest.getBook(id, [options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>ExtendedError</code>](#module_swagger-test--SwaggerTest.Errors.ExtendedError)  
**Reject**: [<code>NotFound</code>](#module_swagger-test--SwaggerTest.Errors.NotFound)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| id | <code>number</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.baggage] | <code>[ &#x27;Map&#x27; ].&lt;string, (string\|number)&gt;</code> | A request-specific baggage to be propagated |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest.RetryPolicies"></a>

#### SwaggerTest.RetryPolicies
Retry policies available to use.

**Kind**: static property of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  

* [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
    * [.Exponential](#module_swagger-test--SwaggerTest.RetryPolicies.Exponential)
    * [.Single](#module_swagger-test--SwaggerTest.RetryPolicies.Single)
    * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)

<a name="module_swagger-test--SwaggerTest.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies)  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies)  
<a name="module_swagger-test--SwaggerTest.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies)  
<a name="module_swagger-test--SwaggerTest.Errors"></a>

#### SwaggerTest.Errors
Errors returned by methods.

**Kind**: static property of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  

* [.Errors](#module_swagger-test--SwaggerTest.Errors)
    * [.ExtendedError](#module_swagger-test--SwaggerTest.Errors.ExtendedError) ⇐ <code>Error</code>
    * [.NotFound](#module_swagger-test--SwaggerTest.Errors.NotFound) ⇐ <code>Error</code>
    * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_swagger-test--SwaggerTest.Errors.ExtendedError"></a>

##### Errors.ExtendedError ⇐ <code>Error</code>
ExtendedError

**Kind**: static class of [<code>Errors</code>](#module_swagger-test--SwaggerTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code | <code>number</code> | 
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.NotFound"></a>

##### Errors.NotFound ⇐ <code>Error</code>
NotFound

**Kind**: static class of [<code>Errors</code>](#module_swagger-test--SwaggerTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of [<code>Errors</code>](#module_swagger-test--SwaggerTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code | <code>number</code> | 
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.DefaultCircuitOptions"></a>

#### SwaggerTest.DefaultCircuitOptions
Default circuit breaker options.

**Kind**: static constant of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
