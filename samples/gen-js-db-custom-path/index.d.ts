import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
  timeout?: number;
  retryPolicy?: RetryPolicy;
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
  keepalive?: boolean;
  retryPolicy?: RetryPolicy;
  logger?: Logger;
  circuit?: CircuitOptions;
  serviceName?: string;
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
    interface ErrorBody {
      message: string;
      [key: string]: any;
    }

    
    class BadRequest {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class InternalError {
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
    type Branch = ("master" | "DEV_BRANCH" | "test");
    
    type Category = ("a" | "b");
    
    type Deployment = {
  application?: string;
  date?: string;
  environment?: string;
  version?: string;
};
    
    type Event = {
  data?: string;
  pk?: string;
  sk?: string;
  ttl?: number;
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
    
    type ThingWithEnumHashKey = {
  branch?: Branch;
  date?: string;
  date2?: string;
};
    
    type ThingWithMatchingKeys = {
  assocID?: string;
  assocType?: string;
  bear?: string;
  created?: string;
};
    
    type ThingWithMultiUseCompositeAttribute = {
  four: string;
  one: string;
  three: string;
  two: string;
};
    
    type ThingWithRequiredCompositePropertiesAndKeysOnly = {
  propertyOne: string;
  propertyThree: string;
  propertyTwo: string;
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
  id_app?: string;
};
    
  }
}

export = SwaggerTest;
