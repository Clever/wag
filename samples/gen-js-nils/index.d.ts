interface StaticClass {
    new(options: nil-testOptions): nil-test;
}

interface nil-test {
Whatever(): IlTokenServiceOptions;
}

type nil-testOptions = {
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