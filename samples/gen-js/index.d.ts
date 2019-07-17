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

type SwaggerTestOptions = (DiscoveryOptions | AddressOptions) & GenericOptions; 


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
  state?: string;
  published?: string;
  snakeCase?: string;
  completed?: string;
  maxPages?: number;
  minPages?: number;
  pagesToTime?: number;
  authorization?: string;
  startingAfter?: number;
};

type OmitEmpty = {
  arrayFieldNotOmitted?: string[];
  arrayFieldOmitted?: string[];
};

type Unathorized = {
  message?: string;
};

declare class SwaggerTest {
  constructor(options: SwaggerTestOptions);

  
  getAuthors(params: GetAuthorsParams, options?: RequestOptions, cb?: Callback<AuthorsResponse>): Promise<AuthorsResponse>
  getAuthorsIter(params: GetAuthorsParams, options: RequestOptions): IterResult<ArrayInner<AuthorsResponse["authorSet"]["results"]>>
  
  getAuthorsWithPut(params: GetAuthorsWithPutParams, options?: RequestOptions, cb?: Callback<AuthorsResponse>): Promise<AuthorsResponse>
  getAuthorsWithPutIter(params: GetAuthorsWithPutParams, options: RequestOptions): IterResult<ArrayInner<AuthorsResponse["authorSet"]["results"]>>
  
  getBooks(params: GetBooksParams, options?: RequestOptions, cb?: Callback<Book[]>): Promise<Book[]>
  getBooksIter(params: GetBooksParams, options: RequestOptions): IterResult<ArrayInner<Book[]>>
  
  createBook(newBook: Book, options?: RequestOptions, cb?: Callback<Book>): Promise<Book>
  
  putBook(newBook?: Book, options?: RequestOptions, cb?: Callback<Book>): Promise<Book>
  
  getBookByID(params: GetBookByIDParams, options?: RequestOptions, cb?: Callback<Book>): Promise<Book>
  
  getBookByID2(id: string, options?: RequestOptions, cb?: Callback<Book>): Promise<Book>
  
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
		
		class Unathorized {
  message?: string;
}
		
		class Error {
  code?: number;
  message?: string;
}
		
  }
}

export = SwaggerTest;
