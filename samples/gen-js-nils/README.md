<a name="module_nil-test"></a>

## nil-test
nil-test client library.


* [nil-test](#module_nil-test)
    * [NilTest](#exp_module_nil-test--NilTest) ⏏
        * [new NilTest(options)](#new_module_nil-test--NilTest_new)
        * _instance_
            * [.nilCheck(params, [options], [cb])](#module_nil-test--NilTest+nilCheck) ⇒ <code>Promise</code>
        * _static_
            * [.RetryPolicies](#module_nil-test--NilTest.RetryPolicies)
                * [.Exponential](#module_nil-test--NilTest.RetryPolicies.Exponential)
                * [.Single](#module_nil-test--NilTest.RetryPolicies.Single)
                * [.None](#module_nil-test--NilTest.RetryPolicies.None)
            * [.Errors](#module_nil-test--NilTest.Errors)
                * [.BadRequest](#module_nil-test--NilTest.Errors.BadRequest) ⇐ <code>Error</code>
                * [.InternalError](#module_nil-test--NilTest.Errors.InternalError) ⇐ <code>Error</code>

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
| [options.discovery] | <code>bool</code> |  | Use @clever/discovery to locate the server. Must provide this or the address argument |
| [options.timeout] | <code>number</code> |  | The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. |
| [options.retryPolicy] | <code>[RetryPolicies](#module_nil-test--NilTest.RetryPolicies)</code> | <code>RetryPolicies.Single</code> | The logic to determine which requests to retry, as well as how many times to retry. |

<a name="module_nil-test--NilTest+nilCheck"></a>

#### nilTest.nilCheck(params, [options], [cb]) ⇒ <code>Promise</code>
Nil check tests

**Kind**: instance method of <code>[NilTest](#exp_module_nil-test--NilTest)</code>  
**Fulfill**: <code>undefined</code>  
**Reject**: <code>[BadRequest](#module_nil-test--NilTest.Errors.BadRequest)</code>  
**Reject**: <code>[InternalError](#module_nil-test--NilTest.Errors.InternalError)</code>  
**Reject**: <code>Error</code>  

| Param | Type | Description |
| --- | --- | --- |
| params | <code>Object</code> |  |
| params.id | <code>string</code> |  |
| [params.query] | <code>string</code> |  |
| [params.header] | <code>string</code> |  |
| [params.body] |  |  |
| [options] | <code>object</code> |  |
| [options.timeout] | <code>number</code> | A request specific timeout |
| [options.span] | <code>[Span](https://doc.esdoc.org/github.com/opentracing/opentracing-javascript/class/src/span.js~Span.html)</code> | An OpenTracing span - For example from the parent request |
| [options.retryPolicy] | <code>[RetryPolicies](#module_nil-test--NilTest.RetryPolicies)</code> | A request specific retryPolicy |
| [cb] | <code>function</code> |  |

<a name="module_nil-test--NilTest.RetryPolicies"></a>

#### NilTest.RetryPolicies
Retry policies available to use.

**Kind**: static property of <code>[NilTest](#exp_module_nil-test--NilTest)</code>  

* [.RetryPolicies](#module_nil-test--NilTest.RetryPolicies)
    * [.Exponential](#module_nil-test--NilTest.RetryPolicies.Exponential)
    * [.Single](#module_nil-test--NilTest.RetryPolicies.Single)
    * [.None](#module_nil-test--NilTest.RetryPolicies.None)

<a name="module_nil-test--NilTest.RetryPolicies.Exponential"></a>

##### RetryPolicies.Exponential
The exponential retry policy will retry five times with an exponential backoff.

**Kind**: static constant of <code>[RetryPolicies](#module_nil-test--NilTest.RetryPolicies)</code>  
<a name="module_nil-test--NilTest.RetryPolicies.Single"></a>

##### RetryPolicies.Single
Use this retry policy to retry a request once.

**Kind**: static constant of <code>[RetryPolicies](#module_nil-test--NilTest.RetryPolicies)</code>  
<a name="module_nil-test--NilTest.RetryPolicies.None"></a>

##### RetryPolicies.None
Use this retry policy to turn off retries.

**Kind**: static constant of <code>[RetryPolicies](#module_nil-test--NilTest.RetryPolicies)</code>  
<a name="module_nil-test--NilTest.Errors"></a>

#### NilTest.Errors
Errors returned by methods.

**Kind**: static property of <code>[NilTest](#exp_module_nil-test--NilTest)</code>  

* [.Errors](#module_nil-test--NilTest.Errors)
    * [.BadRequest](#module_nil-test--NilTest.Errors.BadRequest) ⇐ <code>Error</code>
    * [.InternalError](#module_nil-test--NilTest.Errors.InternalError) ⇐ <code>Error</code>

<a name="module_nil-test--NilTest.Errors.BadRequest"></a>

##### Errors.BadRequest ⇐ <code>Error</code>
BadRequest

**Kind**: static class of <code>[Errors](#module_nil-test--NilTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

<a name="module_nil-test--NilTest.Errors.InternalError"></a>

##### Errors.InternalError ⇐ <code>Error</code>
InternalError

**Kind**: static class of <code>[Errors](#module_nil-test--NilTest.Errors)</code>  
**Extends:** <code>Error</code>  
**Properties**

| Name | Type |
| --- | --- |
| message | <code>string</code> | 

