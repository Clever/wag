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


type Branch = ("master" | "DEV_BRANCH" | "test");

type Category = ("a" | "b");

type Deployment = {
  application?: string;
  date?: string;
  environment?: string;
  version?: string;
};

type NoRangeThingWithCompositeAttributes = {
  branch: string;
  date: string;
  name: string;
  version?: number;
};

type Object = {
  bar?: string;
  foo?: string;
};

type SimpleThing = {
  id?: string;
  name?: string;
};

type TeacherSharingRule = {
  app?: string;
  district?: string;
  id?: string;
  school?: string;
  sections?: string[];
  teacher?: string;
};

type Thing = {
  category?: Category;
  createdAt?: string;
  id?: string;
  name?: string;
  nestedObject?: Object;
  version?: number;
};

type ThingWithCompositeAttributes = {
  branch: string;
  date: string;
  name: string;
  version?: number;
};

type ThingWithCompositeEnumAttributes = {
  branchID: Branch;
  date: string;
  name: string;
};

type ThingWithDateRange = {
  date?: string;
  name?: string;
};

type ThingWithDateTimeComposite = {
  created?: string;
  id?: string;
  resource?: string;
  type?: string;
};

type ThingWithMatchingKeys = {
  assocID?: string;
  assocType?: string;
  bear?: string;
  created?: string;
};

type ThingWithRequiredFields = {
  id: string;
  name: string;
};

type ThingWithRequiredFields2 = {
  id: string;
  name: string;
};

type ThingWithUnderscores = {
  idApp?: string;
};

declare class SwaggerTest {
  constructor(options: SwaggerTestOptions);

  
  healthCheck(options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
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
		
		class InternalError {
  message?: string;
}
		
  }
}

export = SwaggerTest;
