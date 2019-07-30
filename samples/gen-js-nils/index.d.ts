import { Span, Tracer } from "opentracing";
import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  span?: Span;
  retryPolicy?: RetryPolicy;
}

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

type NilTestOptions = (DiscoveryOptions | AddressOptions) & GenericOptions; 


type NilCheckParams = {
  id: string;
  query?: string;
  header?: string;
  array?: string[];
  body?: NilFields;
};

type NilFields = {
  id?: string;
  optional?: string;
};

declare class NilTest {
  constructor(options: NilTestOptions);

  
  nilCheck(params: NilCheckParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
}

declare namespace NilTest {
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
    
    class InternalError {
  message?: string;
}
    
  }
}

export = NilTest;
