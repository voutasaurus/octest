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

Browse to http://localhost:TODO for Jaeger dashboard

## TODO

-[ ] fix docker-compose networking, localhost from within a container doesn't work.
-[ ] figure out which port Jaeger uses for a dashboard, if any.
