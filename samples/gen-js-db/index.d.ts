interface StaticClass {
    new(options: swagger-testOptions): swagger-test;
}

interface swagger-test {
Whatever(): IlTokenServiceOptions;
}

type swagger-testOptions = {
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

declare var ExportedClass: StaticClass;
export = ExportedClass;