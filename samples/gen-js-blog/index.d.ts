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
  tracer?: Tracer;
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

type BlogOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = Blog.Models

declare class Blog {
  constructor(options: BlogOptions);

  
  postGradeFileForStudent(params: models.PostGradeFileForStudentParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  getSectionsForStudent(student_id: string, options?: RequestOptions, cb?: Callback<models.Section[]>): Promise<models.Section[]>
  
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
    
  }
}

export = Blog;
