import { Span, Tracer } from "opentracing";
import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
}

interface RequestOptions {
	/** The timeout to use for all client requests, in milliseconds. */
	timeout?: number;
	/** An OpenTracing span - For example from the parent request */
	span?: Span;
	/** The logic to determine which requests to retry, as well as how many times to retry. */
  retryPolicy?: RetryPolicy;
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
  forEachAsync(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CircuitOptions {
	/** When set to true the circuit will always be closed. Default: true. */
	forceClosed?: boolean;
	/** The maximum number of concurrent requests the client can make at the same time. Default: 100. */
	maxConcurrentRequests?: number;
	/** The minimum number of requests needed before a circuit can be tripped due to health. Default: 20. */
	requestVolumeThreshold?: number;
	/** How long, in milliseconds, to wait after a circuit opens before testing for recovery. Default: 5000. */
	sleepWindow?: number;
	/** The threshold to place on the rolling error rate. Once the error rate exceeds this percentage, the circuit opens. Default: 90. */
  errorPercentThreshold?: number;
}

interface GenericOptions {
	/** The timeout to use for all client requests, in milliseconds. This can be overridden on a per-request basis. Default is 5000ms. */
	timeout?: number;
	/** Set keepalive to true for client requests. This sets the forever: true attribute in request. Defaults to true. */
	keepalive?: boolean;
	/** The logic to determine which requests to retry, as well as how many times to retry. */
	retryPolicy?: RetryPolicy;
	/** The Kayvee logger to use in the client. */
	logger?: Logger;
	/** The OpenTracing Tracer to use. Defaults to the OpenTracing globalTracer. */
	tracer?: Tracer;
	/** Options for constructing the client's circuit breaker. */
	circuit?: CircuitOptions;
	/** Overrides the default service name. This is necessary if the same client is used multiple times, but with different settings, such as with sso- clients. */
  serviceName?: string;
}

interface DiscoveryOptions {
	/** Use clever-discovery to locate the server. Must provide this or the address argument. */
	discovery: true;
  address?: undefined;
}

interface AddressOptions {
	discovery?: false;
	/** URL where the server is located. Must provide this or the discovery argument. */
  address: string;
}

type SwaggerTestOptions = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = SwaggerTest.Models

declare class SwaggerTest {
	/**
	* Create a new client object.
	* @param options - Options for constructing a client object.
	*/
  constructor(options: SwaggerTestOptions);

  
  /**
	  Gets authors
@throws BadRequest
@throws InternalError
	*/
	getAuthors(params: models.GetAuthorsParams, options?: RequestOptions, cb?: Callback<models.AuthorsResponse>): Promise<models.AuthorsResponse>
  /**
	  Gets authors
@throws BadRequest
@throws InternalError
	*/
	getAuthorsIter(params: models.GetAuthorsParams, options?: RequestOptions): IterResult<ArrayInner<models.AuthorsResponse["authorSet"]["results"]>>
  
  /**
	  Gets authors, but needs to use the body so it's a PUT
@throws BadRequest
@throws InternalError
	*/
	getAuthorsWithPut(params: models.GetAuthorsWithPutParams, options?: RequestOptions, cb?: Callback<models.AuthorsResponse>): Promise<models.AuthorsResponse>
  /**
	  Gets authors, but needs to use the body so it's a PUT
@throws BadRequest
@throws InternalError
	*/
	getAuthorsWithPutIter(params: models.GetAuthorsWithPutParams, options?: RequestOptions): IterResult<ArrayInner<models.AuthorsResponse["authorSet"]["results"]>>
  
  /**
	  For a given district:app pair, provides boolean on whether the user exists. This endpoint is preferred over an alternative getUser endpoint if existence confirmation is all that's needed, as this provides significant performance improvements.

		
		Returns a list of books
@throws BadRequest
@throws InternalError
	*/
	getBooks(params: models.GetBooksParams, options?: RequestOptions, cb?: Callback<models.Book[]>): Promise<models.Book[]>
  /**
	  For a given district:app pair, provides boolean on whether the user exists. This endpoint is preferred over an alternative getUser endpoint if existence confirmation is all that's needed, as this provides significant performance improvements.

		
		Returns a list of books
@throws BadRequest
@throws InternalError
	*/
	getBooksIter(params: models.GetBooksParams, options?: RequestOptions): IterResult<ArrayInner<models.Book[]>>
  
  /**
	  Creates a book
@throws BadRequest
@throws InternalError
	*/
	createBook(newBook: models.Book, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  /**
	  Puts a book
@throws BadRequest
@throws InternalError
	*/
	putBook(newBook?: models.Book, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  /**
	  Returns a book
@throws BadRequest
@throws Unathorized
@throws Error
@throws InternalError
	*/
	getBookByID(params: models.GetBookByIDParams, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  /**
	  Retrieve a book
@throws BadRequest
@throws Error
@throws InternalError
	*/
	getBookByID2(id: string, options?: RequestOptions, cb?: Callback<models.Book>): Promise<models.Book>
  
  /**
	  
@throws BadRequest
@throws InternalError
	*/
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
  /** A list of authors. Must specify at least one and at most two */
authors?: string[];
  available?: boolean;
  /** The state of the thing */
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
    
    type OmitEmpty = {
  arrayFieldNotOmitted?: string[];
  arrayFieldOmitted?: string[];
};
    
    type Unathorized = {
  message?: string;
};
    
  }
}

export = SwaggerTest;
