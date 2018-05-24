interface SwaggerTestClass {
  new(options: SwaggerTestOptions): SwaggerTest;
  Errors: any;
}

interface SwaggerTest {
  health(...args: any[]): any;
  
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