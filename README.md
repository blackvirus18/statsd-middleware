# statsd-middleware

This is a Go library which can be added as a middleware in any Go project for statsd metrics. If your project uses hsytrix this library will make the hystrix metrics also flow into statsd.

# How to use the middleware

## Add the middleware in the following way during server start up

```
n.Use(instrumentation.StatsDMiddlewareLogger())

```
## Enable Hystrix metrics to flow into statsd

```
hystrixStreamHandler := hystrix.NewStreamHandler()
hystrixStreamHandler.Start()
StartHystrixStats(config.StatsD(), "hsytrix")

```
