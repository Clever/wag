import { Span, Tracer } from "opentracing";
import { Logger } from "kayvee";

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  span?: Span;
  retryPolicy?: RetryPolicy;
}

type Callback<R> = (err: Error, result: R) => void;

interface IterResult<R> {
	map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
	toArray(cb?: Callback<R[]>): Promise<R[]>;
	forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CallOptions {
  timeout?: number;
  span?: Span;
  retryPolicy?: RetryPolicy;
}

interface CircuitOptions {
  forceClosed?: boolean;
  maxConcurrentRequests?: number;
  requestVolumeThreshold?: number;
  sleepWindow?: number;
  errorPercentThreshold?: number;
}

interface GenericOptions {
  timeout?: number;
  keepalive?: boolean;
  retryPolicy?: RetryPolicy;
	logger?: Logger;
	tracer?: Tracer;
	circuit?: CircuitOptions;
}

interface DiscoveryOptions {
  discovery: true;
  address?: undefined;
}

interface AddressOptions {
  discovery?: false;
  address: string;
}

type SwaggerTestOptions = (DiscoveryOptions | AddressOptions) & GenericOptions; 


declare class SwaggerTest {
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
		
		class BadRequest {
  message?: string;
}
		
		class NotFound {
  message?: string;
}
		
		class InternalError {
  message?: string;
}
		
  }
}

export = SwaggerTest;
