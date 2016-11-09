<a name="module_swagger-test"></a>

## swagger-test
swagger-test client library.


* [swagger-test](#module_swagger-test)
    * [SwaggerTest](#exp_module_swagger-test--SwaggerTest) ⏏
        * [new SwaggerTest(options)](#new_module_swagger-test--SwaggerTest_new)
        * _instance_
            * [.getBooks(params, [cb])](#module_swagger-test--SwaggerTest+getBooks) ⇒ <code>Promise</code>
            * [.createBook(newBook, [cb])](#module_swagger-test--SwaggerTest+createBook) ⇒ <code>Promise</code>
            * [.getBookByID(params, [cb])](#module_swagger-test--SwaggerTest+getBookByID) ⇒ <code>Promise</code>
            * [.getBookByID2(id, [cb])](#module_swagger-test--SwaggerTest+getBookByID2) ⇒ <code>Promise</code>
            * [.healthCheck([cb])](#module_swagger-test--SwaggerTest+healthCheck) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_swagger-test--SwaggerTest.RetryPolicies)
                * [.Default](#module_swagger-test--SwaggerTest.RetryPolicies.Default)
                    * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.backoffs) ⇒ <code>Array.&lt;number&gt;</code>
                    * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.Default.retry) ⇒ <code>boolean</code>
                * [.None](#module_swagger-test--SwaggerTest.RetryPolicies.None)
                    * [.backoffs()](#module_swagger-test--SwaggerTest.RetryPolicies.None.backoffs)
                    * [.retry()](#module_swagger-test--SwaggerTest.RetryPolicies.None.retry)
            * [.Errors](#module_swagger-test--SwaggerTest.Errors)
                * [.BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>
                * [.Unathorized](#module_swagger-test--SwaggerTest.Errors.Unathorized) ⇐ <code>Error</code>
                * [.Error](#module_swagger-test--SwaggerTest.Errors.Error) ⇐ <code>Error</code>

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

<a name="module_swagger-test--SwaggerTest+getBooks"></a>

#### swaggerTest.getBooks(params, [cb]) ⇒ <code>Promise</code>
Returns a list of books

**Kind**: instance method of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  
**Fulfill**: <code>Object[]</code>  
**Reject**: <code>[BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest)</code>  
**Reject**: <code>[InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type | Default |
| --- | --- | --- |
| params | <code>Object</code> |  | 
| [params.authors] | <code>Array.&lt;string&gt;</code> |  | 
| [params.available] | <code>boolean</code> | <code>true</code> | 
| [params.state] | <code>string</code> | <code>&quot;finished&quot;</code> | 
| [params.published] | <code>string</code> |  | 
| [params.snakeCase] | <code>string</code> |  | 
| [params.completed] | <code>string</code> |  | 
| [params.maxPages] | <code>number</code> | <code>500.5</code> | 
| [params.minPages] | <code>number</code> | <code>5</code> | 
| [params.pagesToTime] | <code>number</code> |  | 
| [cb] | <code>function</code> |  | 

<a name="module_swagger-test--SwaggerTest+createBook"></a>

#### swaggerTest.createBook(newBook, [cb]) ⇒ <code>Promise</code>
Creates a book

**Kind**: instance method of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  
**Fulfill**: <code>Object</code>  
**Reject**: <code>[BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest)</code>  
**Reject**: <code>[InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type |
| --- | --- |
| newBook |  | 
| [cb] | <code>function</code> | 

<a name="module_swagger-test--SwaggerTest+getBookByID"></a>

#### swaggerTest.getBookByID(params, [cb]) ⇒ <code>Promise</code>
Returns a book

**Kind**: instance method of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  
**Fulfill**: <code>Object</code>  
**Reject**: <code>[BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest)</code>  
**Reject**: <code>[Unathorized](#module_swagger-test--SwaggerTest.Errors.Unathorized)</code>  
**Reject**: <code>[Error](#module_swagger-test--SwaggerTest.Errors.Error)</code>  
**Reject**: <code>[InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type |
| --- | --- |
| params | <code>Object</code> | 
| params.bookID | <code>number</code> | 
| [params.authorID] | <code>string</code> | 
| [params.authorization] | <code>string</code> | 
| [params.randomBytes] | <code>string</code> | 
| [cb] | <code>function</code> | 

<a name="module_swagger-test--SwaggerTest+getBookByID2"></a>

#### swaggerTest.getBookByID2(id, [cb]) ⇒ <code>Promise</code>
Retrieve a book

**Kind**: instance method of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  
**Fulfill**: <code>Object</code>  
**Reject**: <code>[BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest)</code>  
**Reject**: <code>[Error](#module_swagger-test--SwaggerTest.Errors.Error)</code>  
**Reject**: <code>[InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type |
| --- | --- |
| id | <code>string</code> | 
| [cb] | <code>function</code> | 

<a name="module_swagger-test--SwaggerTest+healthCheck"></a>

#### swaggerTest.healthCheck([cb]) ⇒ <code>Promise</code>
**Kind**: instance method of <code>[SwaggerTest](#exp_module_swagger-test--SwaggerTest)</code>  
**Fulfill**: <code>undefined</code>  
**Reject**: <code>[BadRequest](#module_swagger-test--SwaggerTest.Errors.BadRequest)</code>  
**Reject**: <code>[InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type |
| --- | --- |
| [cb] | <code>function</code> | 

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
    * [.InternalError](#module_swagger-test--SwaggerTest.Errors.InternalError) ⇐ <code>Error</code>
    * [.Unathorized](#module_swagger-test--SwaggerTest.Errors.Unathorized) ⇐ <code>Error</code>
    * [.Error](#module_swagger-test--SwaggerTest.Errors.Error) ⇐ <code>Error</code>

<a name="module_swagger-test--SwaggerTest.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

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

<a name="module_swagger-test--SwaggerTest.Errors.Unathorized"></a>

##### Errors.Unathorized ⇐ <code>Error</code>
Unathorized

**Kind**: static class of <code>[Errors](#module_swagger-test--SwaggerTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_swagger-test--SwaggerTest.Errors.Error"></a>

##### Errors.Error ⇐ <code>Error</code>
Error

**Kind**: static class of <code>[Errors](#module_swagger-test--SwaggerTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| code | <code>number</code> | 
| message | <code>string</code> | 

