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

type BlogOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = Blog.Models

declare class Blog {
  constructor(options: BlogOptions);

  close(): void;
  
  postGradeFileForStudent(params: models.PostGradeFileForStudentParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  getSectionsForStudent(studentID: string, options?: RequestOptions, cb?: Callback<models.Section[]>): Promise<models.Section[]>
  
  postSectionsForStudent(params: models.PostSectionsForStudentParams, options?: RequestOptions, cb?: Callback<models.Section[]>): Promise<models.Section[]>
  
}

declare namespace Blog {
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
    
    type GradeFile = string;
    
    type PostGradeFileForStudentParams = {
  studentID: string;
  file?: GradeFile;
};
    
    type PostSectionsForStudentParams = {
  studentID: string;
  sections: string;
  userType: ("math" | "science" | "reading");
};
    
    type Section = {
  id?: string;
  name?: string;
  period?: string;
};
    
    type SectionType = ("math" | "science" | "reading");
    
    type UnknownResponse = {
  body?: string;
  statusCode?: number;
};
    
  }
}

export = Blog;
