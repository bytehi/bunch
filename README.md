# bunch

# redis-benchmark -h localhost -p 6379 -t ping -c 1000 -n 10000
```
====== PING_INLINE ======
  10000 requests completed in 0.12 seconds
  1000 parallel clients
  3 bytes payload
  keep alive: 1

0.01% <= 2 milliseconds
0.26% <= 3 milliseconds
20.66% <= 4 milliseconds
40.79% <= 5 milliseconds
60.83% <= 6 milliseconds
81.35% <= 7 milliseconds
95.25% <= 8 milliseconds
96.33% <= 9 milliseconds
97.45% <= 10 milliseconds
98.64% <= 11 milliseconds
99.75% <= 12 milliseconds
100.00% <= 12 milliseconds
84033.61 requests per second

====== PING_BULK ======
  10000 requests completed in 0.12 seconds
  1000 parallel clients
  3 bytes payload
  keep alive: 1

0.01% <= 2 milliseconds
0.89% <= 3 milliseconds
18.05% <= 4 milliseconds
39.48% <= 5 milliseconds
60.12% <= 6 milliseconds
82.94% <= 7 milliseconds
95.37% <= 8 milliseconds
96.54% <= 9 milliseconds
97.66% <= 10 milliseconds
98.80% <= 11 milliseconds
99.91% <= 12 milliseconds
100.00% <= 12 milliseconds
84033.61 requests per second
```

# 原生redis-server: redis-benchmark -h localhost -p 30302 -t ping -c 1000 -n 10000
```
====== PING_INLINE ======
  10000 requests completed in 0.12 seconds
  1000 parallel clients
  3 bytes payload
  keep alive: 1

0.01% <= 4 milliseconds
32.97% <= 5 milliseconds
90.31% <= 6 milliseconds
93.21% <= 7 milliseconds
94.65% <= 8 milliseconds
96.92% <= 9 milliseconds
98.82% <= 10 milliseconds
100.00% <= 10 milliseconds
86206.90 requests per second

====== PING_BULK ======
  10000 requests completed in 0.12 seconds
  1000 parallel clients
  3 bytes payload
  keep alive: 1

0.01% <= 4 milliseconds
32.74% <= 5 milliseconds
85.77% <= 6 milliseconds
93.13% <= 7 milliseconds
95.30% <= 8 milliseconds
97.30% <= 9 milliseconds
98.78% <= 10 milliseconds
100.00% <= 10 milliseconds
86206.90 requests per second
```
