<a name="module_blog"></a>

## blog
blog client library.


* [blog](#module_blog)
    * [Blog](#exp_module_blog--Blog) ⏏
        * [new Blog(options)](#new_module_blog--Blog_new)
        * _instance_
            * [.close()](#module_blog--Blog+close)
            * [.postGradeFileForStudent(params, [options], [cb])](#module_blog--Blog+postGradeFileForStudent) ⇒ <code>Promise</code>
            * [.getSectionsForStudent(studentID, [options], [cb])](#module_blog--Blog+getSectionsForStudent) ⇒ <code>Promise</code>
            * [.postSectionsForStudent(params, [options], [cb])](#module_blog--Blog+postSectionsForStudent) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_blog--Blog.RetryPolicies)
                * [.Exponential](#module_blog--Blog.RetryPolicies.Exponential)
                * [.Single](#module_blog--Blog.RetryPolicies.Single)
                * [.None](#module_blog--Blog.RetryPolicies.None)
            * [.Errors](#module_blog--Blog.Errors)
                * [.BadRequest](#module_blog--Blog.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_blog--Blog.Errors.InternalError) ⇐ <code>Error</code>
            * [.DefaultCircuitOptions](#module_blog--Blog.DefaultCircuitOptions)

<a name="exp_module_blog--Blog"></a>

### Blog ⏏
blog client

**Kind**: Exported class  
<a name="new_module_blog--Blog_new"></a>

#### new Blog(options)
Create a new client object.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| options | <code>Object</code> |  | Options for constructing a client object. |
| [options.address] | <code>string</code> |  | URL where the server is located. Must provide this or the discovery argument |
| [options.discovery] | <code>bool</code> |  | Use clever-discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. |
| [options.keepalive] | <code>bool</code> |  | Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies) | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |
| [options.logger] | <code>module:kayvee.Logger</code> | <code>logger.New(&quot;blog-wagclient&quot;)</code> | The Kayvee logger to use in the client. |
| [options.circuit] | <code>Object</code> |  | Options for constructing the client's circuit breaker. |
| [options.circuit.forceClosed] | <code>bool</code> |  | When set to true the circuit will always be closed. Default: true. |
| [options.circuit.maxConcurrentRequests] | <code>number</code> |  | the maximum number of concurrent requests the client can make at the same time. Default: 100. |
| [options.circuit.requestVolumeThreshold] | <code>number</code> |  | The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. |
| [options.circuit.sleepWindow] | <code>number</code> |  | how long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. |
| [options.circuit.errorPercentThreshold] | <code>number</code> |  | the threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. |
| [options.asynclocalstore] | <code>object</code> |  | a request scoped async store |

<a name="module_blog--Blog+close"></a>

#### blog.close()
Releases handles used in client

**Kind**: instance method of [<code>Blog</code>](#exp_module_blog--Blog)  
<a name="module_blog--Blog+postGradeFileForStudent"></a>

#### blog.postGradeFileForStudent(params, [options], [cb]) ⇒ <code>Promise</code>
Posts the grade file for the specified student

**Kind**: instance method of [<code>Blog</code>](#exp_module_blog--Blog)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_blog--Blog.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_blog--Blog.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.studentID | <code>string</code> |  |
| [params.file] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.baggage] | <code>[ &#x27;Map&#x27; ].&lt;string, (string\|number)&gt;</code> | A request-specific baggage to be propagated |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies) | A request specific retryPolicy |
| [options.headers] | <code>[ &#x27;Object&#x27; ].&lt;string, string&gt;</code> | Additional headers to send with the request |
| [cb] | <code>function</code> |  |

<a name="module_blog--Blog+getSectionsForStudent"></a>

#### blog.getSectionsForStudent(studentID, [options], [cb]) ⇒ <code>Promise</code>
Gets the sections for the specified student

**Kind**: instance method of [<code>Blog</code>](#exp_module_blog--Blog)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_blog--Blog.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_blog--Blog.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| studentID | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.baggage] | <code>[ &#x27;Map&#x27; ].&lt;string, (string\|number)&gt;</code> | A request-specific baggage to be propagated |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies) | A request specific retryPolicy |
| [options.headers] | <code>[ &#x27;Object&#x27; ].&lt;string, string&gt;</code> | Additional headers to send with the request |
| [cb] | <code>function</code> |  |

<a name="module_blog--Blog+postSectionsForStudent"></a>

#### blog.postSectionsForStudent(params, [options], [cb]) ⇒ <code>Promise</code>
Posts the sections for the specified student

**Kind**: instance method of [<code>Blog</code>](#exp_module_blog--Blog)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_blog--Blog.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_blog--Blog.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.studentID | <code>string</code> |  |
| params.sections | <code>string</code> |  |
| params.userType | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.baggage] | <code>[ &#x27;Map&#x27; ].&lt;string, (string\|number)&gt;</code> | A request-specific baggage to be propagated |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies) | A request specific retryPolicy |
| [options.headers] | <code>[ &#x27;Object&#x27; ].&lt;string, string&gt;</code> | Additional headers to send with the request |
| [cb] | <code>function</code> |  |

<a name="module_blog--Blog.RetryPolicies"></a>

#### Blog.RetryPolicies
Retry policies available to use.

**Kind**: static property of [<code>Blog</code>](#exp_module_blog--Blog)  

* [.RetryPolicies](#module_blog--Blog.RetryPolicies)
    * [.Exponential](#module_blog--Blog.RetryPolicies.Exponential)
    * [.Single](#module_blog--Blog.RetryPolicies.Single)
    * [.None](#module_blog--Blog.RetryPolicies.None)

<a name="module_blog--Blog.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies)  
<a name="module_blog--Blog.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies)  
<a name="module_blog--Blog.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of [<code>RetryPolicies</code>](#module_blog--Blog.RetryPolicies)  
<a name="module_blog--Blog.Errors"></a>

#### Blog.Errors
Errors returned by methods.

**Kind**: static property of [<code>Blog</code>](#exp_module_blog--Blog)  

* [.Errors](#module_blog--Blog.Errors)
    * [.BadRequest](#module_blog--Blog.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_blog--Blog.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_blog--Blog.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of [<code>Errors</code>](#module_blog--Blog.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_blog--Blog.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of [<code>Errors</code>](#module_blog--Blog.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_blog--Blog.DefaultCircuitOptions"></a>

#### Blog.DefaultCircuitOptions
Default circuit breaker options.

**Kind**: static constant of [<code>Blog</code>](#exp_module_blog--Blog)  
