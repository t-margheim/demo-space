## Need

The Chronicle server code was originally written with a strong preference for explicit dependencies, requiring code authors to pass loggers, metrics providers, and tracing providers as arguments to functions and/or properties on structs. While the idea behind this was to improve clarity of code and provide easy dependency injection for unit testing, over time the team has moved away from this approach with logging specifically. `rslog` provides simple functions for us to use on an as-needed basis, like `rslog.Errorf`, meaning that we don't have to modify a chain of functions between `main` and some internal service operation in order to add a log.

Since we've gone that direction with logging, I think it makes sense to try to keep our other observability pals (metrics and tracing) consistent with that behavior where possible. 

## Approach

We would create a new `rsmetrics` package. Eventually this could grow into something that would make sense to live in `platform-lib`, but I'd suggest building and using it in just Chronicle server until we feel like the API is right.

The initial API would consist of methods like these:

- `func Initialize(m metric.Meter) error`: The calling code is responsible for creating a `go.opentelemetry.io/otel/metric.Meter`, which this package then uses to create new metrics as needed. This also means that the calling code is not limited by the capabilities of this package, as it will still be able to interact directly with OpenTelemetry SDK if needed. This function returns an error only in the case that `nil` is passed as the meter.
- `func Count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue)`: Count provides a way to increment a counter by a specified value.
- `func Timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue)`: Timing provides a way to record timing data to histograms using OTEL.

These are the only two metrics collecting functions we are currently using in the Chronicle server, but eventually it would be good to support all of the different types of metrics Opentelemetry supports. The signatures here are based on the `statsd` package from Datadog, but that's definitely just a personal preference. We could tie them more closely to the OpenTelemetry API, but that's a bit of a moving target.

## Benefits

This approach would provide the ease-of-use described in the Need above and still maintain a simple unit testing capability. As I was working on this proposal, I discovered an additional benefit. The [OpenTelemetry SDK](https://opentelemetry.io/docs/instrumentation/go/) for metrics that we are using is very much a beta product. There has been a pretty dramatic change to the `metric.Meter` interface since we implemented metrics in the Chronicle server. By isolating our exposure to this volatile API to a single package, we will be able to stay up to date with the SDK without needing to change code in many different places.  

Here's the (breaking) change introduced recently in `metric.Meter`. We are currently importing v0.34.0, and would have a number of places that we need to update in order to upgrade to v0.37.0

v0.34.0 `metric.Meter` interface:

``` go
type Meter interface {
	AsyncInt64() asyncint64.InstrumentProvider
	AsyncFloat64() asyncfloat64.InstrumentProvider
	RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error
	SyncInt64() syncint64.InstrumentProvider
	SyncFloat64() syncfloat64.InstrumentProvider
}
```

v0.37.0 `metric.Meter` interface:

``` go
type Meter interface {
	Int64Counter(name string, options ...instrument.Int64Option) (instrument.Int64Counter, error)
	Int64UpDownCounter(name string, options ...instrument.Int64Option) (instrument.Int64UpDownCounter, error)
	Int64Histogram(name string, options ...instrument.Int64Option) (instrument.Int64Histogram, error)
	Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableCounter, error)
	Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableUpDownCounter, error)
	Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableGauge, error)
	Float64Counter(name string, options ...instrument.Float64Option) (instrument.Float64Counter, error)
	Float64UpDownCounter(name string, options ...instrument.Float64Option) (instrument.Float64UpDownCounter, error)
	Float64Histogram(name string, options ...instrument.Float64Option) (instrument.Float64Histogram, error)
	Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableCounter, error)
	Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableUpDownCounter, error)
	Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableGauge, error)
	RegisterCallback(f Callback, instruments ...instrument.Asynchronous) (Registration, error)
}
```

## Costs

This would add an additional layer that we would be responsible for maintaining and documenting.