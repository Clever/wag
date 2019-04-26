import { Span } from "opentracing";

interface RetryPolicy {
  backoffs(): number[],
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean,
}

interface RetryPolicies {
  Single: RetryPolicy,
  Exponential: RetryPolicy,
  None: RetryPolicy,
}

interface CallOptions {
  timeout?: number,
  span?: Span,
  retryPolicy?: RetryPolicy
}

type Callback<R> = (Error, R) => void;

interface CallOptions {
  timeout?: number,
  span?: Span,
  retryPolicy?: RetryPolicy,
}

interface CircuitOptions {
  forceClosed: boolean;
  maxConcurrentRequests: number;
  requestVolumeThreshold: number;
  sleepWindow: number;
  errorPercentThreshold: number;
}

interface GenericOptions {
  timeout: number;
  keepalive: boolean;
  retryPolicy: RetryPolicy;
  logger: Logger;
}

interface DiscoveryOptions {
  discovery: true;
  address?: undefined;
}

interface AddressOptions {
  discovery?: false;
  address: string;
}

type swagger-testOptions = (DiscoveryOptions | AddressOptions) & GenericOptions; 


declare class swagger-test {
  constructor(options: swagger-testOptions);

  
}

declare namespace swagger-test {
  Errors: interface {
    BadRequest: Error,
    InternalError: Error,
    NotFound: Error,
  }
}

export = swagger-test;

