interface SwaggerTestClass {
    new(options: SwaggerTestOptions): SwaggerTest;
}

interface SwaggerTest {

  getAuthors(...args: any[]): any;

  getAuthorsWithPut(...args: any[]): any;

  getBooks(...args: any[]): any;

  createBook(...args: any[]): any;

  putBook(...args: any[]): any;

  getBookByID(...args: any[]): any;

  getBookByID2(...args: any[]): any;

  healthCheck(...args: any[]): any;

}

type SwaggerTestOptions = {
  timeout?: number;
  retryPolicy?: any;
  logger?: any;
  circuit?: CircuitBreakerOptions;
} & (
{
  discovery: true;
}
|
{
  discovery?: false;
  address: string;
}
)

type CircuitBreakerOptions = {
  forceClosed:            boolean;
  requestVolumeThreshold: number;
  maxConcurrentRequests:  number;
  sleepWindow:            number;
  errorPercentThreshold:  number;
  logIntervalMs:          number;
}

type RequestOptions = {
  timeout?: number;
  span?: any; // opentracing span
  retryPolicy?: RetryPolicy;
}

interface RetryPolicy {
  backoffs(): void;
  retry(requestOptions: any, err: any, res: any): any;
}

declare var ExportedClass: SwaggerTestClass;
export = ExportedClass;