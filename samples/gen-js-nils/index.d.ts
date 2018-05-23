interface NilTestClass {
    new(options: NilTestOptions): NilTest;
}

interface NilTest {

  nilCheck(...args: any[]): any;

}

type NilTestOptions = {
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

declare var ExportedClass: NilTestClass;
export = ExportedClass;