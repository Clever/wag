import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  baggage?: Map<string, string | number>;
  retryPolicy?: RetryPolicy;
  headers?: { [key: string]: string };
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
  forEachAsync(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
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
  baggage?: Map<string, string | number>;
  keepalive?: boolean;
  retryPolicy?: RetryPolicy;
  logger?: Logger;
  circuit?: CircuitOptions;
  serviceName?: string;
  asynclocalstore?: object;
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

import models = SwaggerTest.Models

declare class SwaggerTest {
  constructor(options: SwaggerTestOptions);

  close(): void;
  
  getBook(id: number, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
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

    
    class ExtendedError {
  code?: number;
  message?: string;

  constructor(body: ErrorBody);
}
    
    class NotFound {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class InternalError {
  code?: number;
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
    type ExtendedError = {
  code?: number;
  message?: string;
};
    
    type UnknownResponse = {
  body?: string;
  statusCode?: number;
};
    
  }
}

export = SwaggerTest;
