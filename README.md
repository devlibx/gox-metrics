### Statsd

To test in local you can run Statsd server using following. ```http://localhost:8081/``` to see the data.
```shell
docker run --rm -it -p 8081:8080 -p 8125:8125/udp -p 8125:8125/tcp  rapidloop/statsd-vis -statsdudp 0.0.0.0:8125 -statsdtcp 0.0.0.0:8125

OR

docker run -d --name graphite --restart=always \
 -p 8081:80 \
 -p 2003-2004:2003-2004 \
 -p 2023-2024:2023-2024 \
 -p 8125:8125/udp \
 -p 8126:8126 \
 graphiteapp/graphite-statsd
 
Launch - http://localhost:8081 
```

#### How to use 
```go
config := metrics.Config{
    Prefix:              "some",
    ReportingIntervalMs: 10,
    Statsd: metrics.StatsdConfig{
            Address:         "127.0.0.1:8125",
            FlushIntervalMs: 10,
            FlushBytes:      1440,
			Properties: map[string]interface{}{"comma_perpetrated_stats_reporter": true},
        },
    }

statsdService, err := NewRootScope(config)
assert.NoError(t, err)

counter := statsdService.Counter("some_counter")
go func() {
    for i := 0; i < 100000; i++ {
        counter.Inc(1)
        time.Sleep(1 * time.Second)
        fmt.Println("submitting statsd counter")
    }
}()
```

#### Build and test
```shell
go test ./... -v -real.statsd=true

Note -  -real.statsd=true will generate stats and need running statsd server
```