# octest

A testbed for basic opencensus setup for a Go application, reporting metrics
via Prometheus and distributed traces via Jaeger.

## Get and run

```
$ git clone https://github.com/voutasaurus/octest
$ cd octest
$ ./set/test/run.sh
```

Browse to http://localhost:8080 for Go App

Browse to http://localhost:9090 for Prometheus dashboard

Browse to http://localhost:16686 for Jaeger dashboard
