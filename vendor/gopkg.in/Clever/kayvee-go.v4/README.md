# kayvee
--
    import "gopkg.in/Clever/kayvee-go.v4"

Package kayvee provides methods to output human and machine parseable strings,
with a "json" format.

## [Logger API Documentation](./logger)

* [gopkg.in/Clever/kayvee-go.v4/logger](https://godoc.org/gopkg.in/Clever/kayvee-go.v4/logger)
* [gopkg.in/Clever/kayvee-go.v4/middleware](https://godoc.org/gopkg.in/Clever/kayvee-go.v4/middleware)

## Example

```go
    package main

    import(
        "fmt"
        "time"

        "gopkg.in/Clever/kayvee-go.v4/logger"
    )

    func main() {
        myLogger := logger.New("myApp")

        // Simple debugging
        myLogger.Debug("Service has started")

        // Make a query and log its length
        query_start := time.Now()
        myLogger.GaugeFloat("QueryTime", time.Since(query_start).Seconds())

        // Output structured data
        myLogger.InfoD("DataResults", logger.M{"key": "value"})

        // You can use the M alias for your key value pairs
        myLogger.InfoD("DataResults", logger.M{"shorter": "line"})
    }
```


## Testing

Run `make test` to execute the tests

## Change log

- v4.0
  - Added methods to read and write the `Logger` object from a a `context.Context` object.
  - Middleware now injects the logger into the request context.
  - Updated to require Go 1.7.
- v4.0 - Removed sentry-go dependency
- v2.4 - Add kayvee-go/validator for asserting that raw log lines are in a valid kayvee format.
- v2.3 - Expose logger.M.
- v2.2 - Remove godeps.
- v2.1 - Add kayvee-go/logger with log level, counters, and gauge support
- v0.1 - Initial release.

## Backward Compatibility

The kayvee 1.x interface still exist but is considered deprecated. You can find documentation on using it in the [compatibility guide](./compatibility.md)
