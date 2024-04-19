<a name="module_swagger-test"></a>

## swagger-test
swagger-test client library.


* [swagger-test](#module_swagger-test)
    * [SwaggerTest](#exp_module_swagger-test--SwaggerTest) ⏏
        * [new SwaggerTest(options)](#new_module_swagger-test--SwaggerTest_new)
        * _instance_
            * [.close()](#module_swagger-test--SwaggerTest+close)
            * [.getAuthors(params, [options], [cb])](#module_swagger-test--SwaggerTest+getAuthors) ⇒ <code>Promise</code>
            * [.getAuthorsIter(params, [options])](#module_swagger-test--SwaggerTest+getAuthorsIter) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
            * [.getAuthorsWithPut(params, [options], [cb])](#module_swagger-test--SwaggerTest+getAuthorsWithPut) ⇒ <code>Promise</code>
            * [.getAuthorsWithPutIter(params, [options])](#module_swagger-test--SwaggerTest+getAuthorsWithPutIter) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
            * [.getBooks(params, [options], [cb])](#module_swagger-test--SwaggerTest+getBooks) ⇒ <code>Promise</code>
            * [.getBooksIter(params, [options])](#module_swagger-test--SwaggerTest+getBooksIter) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
            * [.createBook(newBook, [options], [cb])](#module_swagger-test--SwaggerTest+createBook) ⇒ <code>Promise</code>
            * [.putBook(newBook, [options], [cb])](#module_swagger-test--SwaggerTest+putBook) ⇒ <code>Promise</code>
            * [.getBookByID(params, [options], [cb])](#module_swagger-test--SwaggerTest+getBookByID) ⇒ <code>Promise</code>
            * [.getBookByID2(id, [options], [cb])](#module_swagger-test--SwaggerTest+getBookByID2) ⇒ <code>Promise</code>
            * [.healthCheck([options], [cb])](#module_swagger-test--SwaggerTest+healthCheck) ⇒ <code>Promise</code>
            * [.lowercaseModelsTest(params, [options], [cb])](#module_swagger-test--SwaggerTest+lowercaseModelsTest) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
                * [.Exponential](#module_swagger-test--SwaggerTest.RetryPolicies.Exponential)
                * [.Single](#module_swagger-test--SwaggerTest.RetryPolicies.Single)
                * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
            * [.Errors](#module_swagger-test--SwaggerTest.Errors)
                * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>
                * [.Unathorized](#module_swagger-test--SwaggerTest.Errors.Unathorized) ⇐ <code>Error</code>
                * [.Error](#module_swagger-test--SwaggerTest.Errors.Error) ⇐ <code>Error</code>
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

<a name="module_swagger-test--SwaggerTest+close"></a>

#### swaggerTest.close()
Releases handles used in client

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
<a name="module_swagger-test--SwaggerTest+getAuthors"></a>

#### swaggerTest.getAuthors(params, [options], [cb]) ⇒ <code>Promise</code>
Gets authors

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.name] | <code>string</code> |  |
| [params.startingAfter] | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+getAuthorsIter"></a>

#### swaggerTest.getAuthorsIter(params, [options]) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
Gets authors

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Returns**: <code>Object</code> - iter<code>function</code> - iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array<code>function</code> - iter.toArray - returns a promise to the resources as an array<code>function</code> - iter.forEach - takes in a function, applies it to each resource<code>function</code> - iter.forEachAsync - takes in an async function, applies it to each resource  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.name] | <code>string</code> |  |
| [params.startingAfter] | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |

<a name="module_swagger-test--SwaggerTest+getAuthorsWithPut"></a>

#### swaggerTest.getAuthorsWithPut(params, [options], [cb]) ⇒ <code>Promise</code>
Gets authors, but needs to use the body so it's a PUT

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.name] | <code>string</code> |  |
| [params.startingAfter] | <code>string</code> |  |
| [params.favoriteBooks] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+getAuthorsWithPutIter"></a>

#### swaggerTest.getAuthorsWithPutIter(params, [options]) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
Gets authors, but needs to use the body so it's a PUT

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Returns**: <code>Object</code> - iter<code>function</code> - iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array<code>function</code> - iter.toArray - returns a promise to the resources as an array<code>function</code> - iter.forEach - takes in a function, applies it to each resource<code>function</code> - iter.forEachAsync - takes in an async function, applies it to each resource  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| [params.name] | <code>string</code> |  |
| [params.startingAfter] | <code>string</code> |  |
| [params.favoriteBooks] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |

<a name="module_swagger-test--SwaggerTest+getBooks"></a>

#### swaggerTest.getBooks(params, [options], [cb]) ⇒ <code>Promise</code>
Returns a list of books

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object[]</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| params | <code>Object</code> |  |  |
| [params.authors] | <code>[ &#x27;Array&#x27; ].&lt;string&gt;</code> |  | A list of authors. Must specify at least one and at most two |
| [params.available] | <code>boolean</code> | <code>true</code> |  |
| [params.state] | <code>string</code> | <code>&quot;finished&quot;</code> |  |
| [params.published] | <code>string</code> |  |  |
| [params.snakeCase] | <code>string</code> |  |  |
| [params.completed] | <code>string</code> |  |  |
| [params.maxPages] | <code>number</code> | <code>500.5</code> |  |
| [params.minPages] | <code>number</code> | <code>5</code> |  |
| [params.pagesToTime] | <code>number</code> |  |  |
| [params.authorization] | <code>string</code> |  |  |
| [params.startingAfter] | <code>number</code> |  |  |
| [options] | <code>object</code> |  |  |
| [options.timeout] | <code>number</code> |  | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) |  | A request specific retryPolicy |
| [cb] | <code>function</code> |  |  |

<a name="module_swagger-test--SwaggerTest+getBooksIter"></a>

#### swaggerTest.getBooksIter(params, [options]) ⇒ <code>Object</code> \| <code>function</code> \| <code>function</code> \| <code>function</code> \| <code>function</code>
Returns a list of books

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Returns**: <code>Object</code> - iter<code>function</code> - iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array<code>function</code> - iter.toArray - returns a promise to the resources as an array<code>function</code> - iter.forEach - takes in a function, applies it to each resource<code>function</code> - iter.forEachAsync - takes in an async function, applies it to each resource  

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| params | <code>Object</code> |  |  |
| [params.authors] | <code>[ &#x27;Array&#x27; ].&lt;string&gt;</code> |  | A list of authors. Must specify at least one and at most two |
| [params.available] | <code>boolean</code> | <code>true</code> |  |
| [params.state] | <code>string</code> | <code>&quot;finished&quot;</code> |  |
| [params.published] | <code>string</code> |  |  |
| [params.snakeCase] | <code>string</code> |  |  |
| [params.completed] | <code>string</code> |  |  |
| [params.maxPages] | <code>number</code> | <code>500.5</code> |  |
| [params.minPages] | <code>number</code> | <code>5</code> |  |
| [params.pagesToTime] | <code>number</code> |  |  |
| [params.authorization] | <code>string</code> |  |  |
| [params.startingAfter] | <code>number</code> |  |  |
| [options] | <code>object</code> |  |  |
| [options.timeout] | <code>number</code> |  | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) |  | A request specific retryPolicy |

<a name="module_swagger-test--SwaggerTest+createBook"></a>

#### swaggerTest.createBook(newBook, [options], [cb]) ⇒ <code>Promise</code>
Creates a book

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| newBook |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+putBook"></a>

#### swaggerTest.putBook(newBook, [options], [cb]) ⇒ <code>Promise</code>
Puts a book

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| newBook |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+getBookByID"></a>

#### swaggerTest.getBookByID(params, [options], [cb]) ⇒ <code>Promise</code>
Returns a book

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>Unathorized</code>](#module_swagger-test--SwaggerTest.Errors.Unathorized)  
**Reject**: [<code>Error</code>](#module_swagger-test--SwaggerTest.Errors.Error)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.bookID | <code>number</code> |  |
| [params.authorID] | <code>string</code> |  |
| [params.authorization] | <code>string</code> |  |
| [params.XDontRateLimitMeBro] | <code>string</code> |  |
| [params.randomBytes] | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+getBookByID2"></a>

#### swaggerTest.getBookByID2(id, [options], [cb]) ⇒ <code>Promise</code>
Retrieve a book

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>Object</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>Error</code>](#module_swagger-test--SwaggerTest.Errors.Error)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| id | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+healthCheck"></a>

#### swaggerTest.healthCheck([options], [cb]) ⇒ <code>Promise</code>
**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.retryPolicy] | [<code>RetryPolicies</code>](#module_swagger-test--SwaggerTest.RetryPolicies) | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_swagger-test--SwaggerTest+lowercaseModelsTest"></a>

#### swaggerTest.lowercaseModelsTest(params, [options], [cb]) ⇒ <code>Promise</code>
testing that we can use a lowercase name for a model

**Kind**: instance method of [<code>SwaggerTest</code>](#exp_module_swagger-test--SwaggerTest)  
**Fulfill**: <code>undefined</code>  
**Reject**: [<code>BadRequest</code>](#module_swagger-test--SwaggerTest.Errors.BadRequest)  
**Reject**: [<code>InternalError</code>](#module_swagger-test--SwaggerTest.Errors.InternalError)  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.lowercase |  |  |
| params.pathParam | <code>string</code> |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
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
    * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>
    * [.Unathorized](#module_swagger-test--SwaggerTest.Errors.Unathorized) ⇐ <code>Error</code>
    * [.Error](#module_swagger-test--SwaggerTest.Errors.Error) ⇐ <code>Error</code>

<a name="module_swagger-test--SwaggerTest.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

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
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.Unathorized"></a>

##### Errors.Unathorized ⇐ <code>Error</code>
Unathorized

**Kind**: static class of [<code>Errors</code>](#module_swagger-test--SwaggerTest.Errors)  
**Extends**: <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.Error"></a>

##### Errors.Error ⇐ <code>Error</code>
Error

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
