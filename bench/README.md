Unload:
```shell
Running 1m test @ http://192.168.56.1:8090/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   183.18ms   99.59ms 826.08ms   71.50%
    Req/Sec   273.77    101.69     1.20k    79.53%
  Latency Distribution
     50%  223.54ms
     75%  241.33ms
     90%  267.15ms
     99%  325.42ms
  325072 requests in 1.00m, 39.06MB read
  Socket errors: connect 0, read 53, write 0, timeout 0
  Non-2xx or 3xx responses: 19
Requests/sec:   5409.40
Transfer/sec:    665.60KB
```

Traefik:
```shell
Running 1m test @ http://test.traefik:8000/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   181.26ms  147.99ms 809.41ms   50.72%
    Req/Sec   291.67     89.69     0.97k    72.41%
  Latency Distribution
     50%  199.24ms
     75%  291.60ms
     90%  370.37ms
     99%  538.59ms
  345813 requests in 1.00m, 33.64MB read
Requests/sec:   5754.70
Transfer/sec:    573.22KB
```

To be updated...