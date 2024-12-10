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
  
  getAuthors(params: models.GetAuthorsParams, options?: RequestOptions, cb?: Callback<models.AuthorsResponse>): Promise<models.AuthorsResponse>
  getAuthorsIter(params: models.GetAuthorsParams, options?: RequestOptions): IterResult<ArrayInner<models.AuthorsResponse["authorSet"]["results"]>>
  
  getAuthorsWithPut(params: models.GetAuthorsWithPutParams, options?: RequestOptions, cb?: Callback<models.AuthorsResponse>): Promise<models.AuthorsResponse>
  getAuthorsWithPutIter(params: models.GetAuthorsWithPutParams, options?: RequestOptions): IterResult<ArrayInner<models.AuthorsResponse["authorSet"]["results"]>>
  
  getBooks(params: models.GetBooksParams, options?: RequestOptions, cb?: Callback<models.Book[]>): Promise<models.Book[]>
  getBooksIter(params: models.GetBooksParams, options?: RequestOptions): IterResult<ArrayInner<models.Book[]>>
  
  createBook(newBook: models.Book, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  putBook(newBook?: models.Book, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  getBookByID(params: models.GetBookByIDParams, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  getBookByID2(id: string, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  healthCheck(options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
  lowercaseModelsTest(params: models.LowercaseModelsTestParams, options?: RequestOptions, cb?: Callback<void>): Promise<void>
  
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
    
    class Unathorized {
  message?: string;

  constructor(body: ErrorBody);
}
    
    class Error {
  code?: number;
  message?: string;

  constructor(body: ErrorBody);
}
    
  }

  namespace Models {
    
    type Animal = {
  age?: number;
  species?: string;
};
    
    type Author = {
  id?: string;
  name?: string;
};
    
    type AuthorArray = Author[];
    
    type AuthorSet = {
  randomProp?: number;
  results?: AuthorArray;
};
    
    type AuthorsResponse = {
  authorSet?: AuthorSet;
  metadata?: AuthorsResponseMetadata;
};
    
    type AuthorsResponseMetadata = {
  count?: number;
};
    
    type Book = {
  author?: string;
  genre?: ("scifi" | "mystery" | "horror");
  id?: number;
  name?: string;
  other?: { [key: string]: string };
  otherArray?: { [key: string]: string[] };
};
    
    type Dog = Pet & Identifiable & {
  breed?: string;
};
    
    type Error = {
  code?: number;
  message?: string;
};
    
    type GetAuthorsParams = {
  name?: string;
  startingAfter?: string;
};
    
    type GetAuthorsWithPutParams = {
  name?: string;
  startingAfter?: string;
  favoriteBooks?: Book;
};
    
    type GetBookByIDParams = {
  bookID: number;
  authorID?: string;
  authorization?: string;
  XDontRateLimitMeBro?: string;
  randomBytes?: string;
};
    
    type GetBooksParams = {
  authors?: string[];
  available?: boolean;
  state?: ("finished" | "inprogress");
  published?: string;
  snakeCase?: string;
  completed?: string;
  maxPages?: number;
  minPages?: number;
  pagesToTime?: number;
  authorization?: string;
  startingAfter?: number;
};
    
    type Identifiable = {
  id?: string;
};
    
    type LowercaseModelsTestParams = {
  lowercase: lowercase;
  pathParam: string;
};
    
    type OmitEmpty = {
  arrayFieldNotOmitted?: string[];
  arrayFieldOmitted?: string[];
};
    
    type Pet = Animal & {
  name?: string;
};
    
    type Unathorized = {
  message?: string;
};
    
    type UnknownResponse = {
  body?: string;
  statusCode?: number;
};
    
    type lowercase = {
  
};
    
  }
}

export = SwaggerTest;
