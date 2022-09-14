<a name="module_nil-test"></a>

## nil-test
nil-test client library.


* [nil-test](#module_nil-test)
    * [NilTest](#exp_module_nil-test--NilTest) ⏏
        * [new NilTest(options)](#new_module_nil-test--NilTest_new)
        * _instance_
            * [.close()](#module_nil-test--NilTest+close)
            * [.nilCheck(params, [options], [cb])](#module_nil-test--NilTest+nilCheck) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_nil-test--NilTest.RetryPolicies)
                * [.Exponential](#module_nil-test--NilTest.RetryPolicies.Exponential)
                * [.Single](#module_nil-test--NilTest.RetryPolicies.Single)
                * [.None](#module_nil-test--NilTest.RetryPolicies.None)
            * [.Errors](#module_nil-test--NilTest.Errors)
                * [.BadRequest](#module_nil-test--NilTest.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_nil-test--NilTest.Errors.InternalError) ⇐ <code>Error</code>
            * [.DefaultCircuitOptions](#module_nil-test--NilTest.DefaultCircuitOptions)

<a name="exp_module_nil-test--NilTest"></a>

### NilTest ⏏
nil-test client

**Kind**: Exported class  
<a name="new_module_nil-test--NilTest_new"></a>

#### new NilTest(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| [options.address] | <code>string</code> |  | URL where the server is located. Must provide this or the discovery argument |
| [options.discovery] | <code>bool</code> |  | Use clever-discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. |
| [options.keepalive] | <code>bool</code> |  | Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_nil-test--NilTest.RetryPolicies) | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;nil-test-wagclient&quot;)</code> | The Kayvee logger to use in the client. |
| [options.circuit] | <code>Object</code> |  | Options for constructing the client's circuit breaker. |
| [options.circuit.forceClosed] | <code>bool</code> |  | When set to true the circuit will always be closed. Default: true. |
| [options.circuit.maxConcurrentRequests] | <code>number</code> |  | the maximum number of concurrent requests the client can make at the same time. Default: 100. |
| [options.circuit.requestVolumeThreshold] | <code>number</code> |  | The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. |
| [options.circuit.sleepWindow] | <code>number</code> |  | how long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. |
| [options.circuit.errorPercentThreshold] | <code>number</code> |  | the threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. |

<a name="module_nil-test--NilTest+close"></a>

#### nilTest.close()
Releases handles used in client

**Kind**: instance method of [<code>NilTest</code>](#exp_module_nil-test--NilTest)  
<a name="module_nil-test--NilTest+nilCheck"></a>

#### nilTest.nilCheck(params, [options], [cb]) ⇒ <code>Promise</code>
Nil check tests

**Kind**: instance method of [<code>NilTest</code>](#exp_module_nil-test--NilTest)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_nil-test--NilTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_nil-test--NilTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.id | <code>string</code> |  |
| [params.query] | <code>string</code> |  |
| [params.header] | <code>string</code> |  |
| [params.array] | <code>Array.&lt;string&gt;</code> |  |
| [params.body] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_nil-test--NilTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_nil-test--NilTest.RetryPolicies"></a>

#### NilTest.RetryPolicies
Retry policies available to use.

**Kind**: static property of [<code>NilTest</code>](#exp_module_nil-test--NilTest)  

* [.RetryPolicies](#module_nil-test--NilTest.RetryPolicies)
    * [.Exponential](#module_nil-test--NilTest.RetryPolicies.Exponential)
    * [.Single](#module_nil-test--NilTest.RetryPolicies.Single)
    * [.None](#module_nil-test--NilTest.RetryPolicies.None)

<a name="module_nil-test--NilTest.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_nil-test--NilTest.RetryPolicies)  
<a name="module_nil-test--NilTest.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_nil-test--NilTest.RetryPolicies)  
<a name="module_nil-test--NilTest.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_nil-test--NilTest.RetryPolicies)  
<a name="module_nil-test--NilTest.Errors"></a>

#### NilTest.Errors
Errors returned by methods.

**Kind**: static property of [<code>NilTest</code>](#exp_module_nil-test--NilTest)  

* [.Errors](#module_nil-test--NilTest.Errors)
    * [.BadRequest](#module_nil-test--NilTest.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_nil-test--NilTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_nil-test--NilTest.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of [<code>Errors</code>](#module_nil-test--NilTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_nil-test--NilTest.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of [<code>Errors</code>](#module_nil-test--NilTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_nil-test--NilTest.DefaultCircuitOptions"></a>

#### NilTest.DefaultCircuitOptions
Default circuit breaker options.

**Kind**: static constant of [<code>NilTest</code>](#exp_module_nil-test--NilTest)  
