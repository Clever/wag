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

type SwaggerTestOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = SwaggerTest.Models

declare class SwaggerTest {
	/**
	* Create a new client object.
	* @param options - Options for constructing a client object.
	*/
  constructor(options: SwaggerTestOptions);

  
}

declare namespace SwaggerTest {
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

    
    class BadRequest {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class NotFound {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class InternalError {
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
  }
}

export = SwaggerTest;
