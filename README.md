# 同步非阻塞的golang libray: bunch
- Call: 访问其它服务(IO密集型最佳)，不阻塞主逻辑
- After: 主逻辑单线程,同步写代码,减少异步callback
- 设计为library,而非framework,方便被集成

# 游戏服务器使用礼包码举例
```
	//主逻辑tick中集成如下代码
	select {
    case f := <-b.AfterQ():
      f() 
  } 

	//处理玩家使用礼包码的请求
  b.NewCalls().Call(func() (interface{}, error) {
    //阻塞访问礼包码,服务，校验是否可以领取礼包
    //return 礼包信息,错误码
  }).After(func(i interface{}, e error) {
    // 主逻辑操作玩家数据
    // i: 礼包信息
    // 发放礼包
  }).Call(func() (interface{}, error) {
    //阻塞通知数据统计平台
  }).After(func(i interface{}, e error) {
    //回归主逻辑
  }).Commit()
```

# 测试异步切换开销,so easy
### redis-benchmark -h localhost -p 6379 -t ping -c 1000 -n 10000
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

### 原生redis-server: redis-benchmark -h localhost -p 30302 -t ping -c 1000 -n 10000
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
