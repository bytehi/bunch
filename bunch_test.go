package bunch

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestBunch(t *testing.T) {
	b := New()

	uniq := int32(0)
	limit := int32(100000)
	doneCnt := int32(0)
	mustUnique := func() {
		if atomic.AddInt32(&uniq, 1) != 1 {
			panic("must unique")
		}
		atomic.AddInt32(&uniq, -1)
	}

	go func() {
		for i := int32(0); i < limit; i++ {
			b.NewCalls().Call(func() (interface{}, error) {
				time.Sleep(time.Millisecond * 5)
				return nil, nil
			}).After(func(i interface{}, e error) {
				mustUnique()
			}).Call(func() (interface{}, error) {
				time.Sleep(time.Millisecond * 5)
				return nil, nil
			}).After(func(i interface{}, e error) {
				mustUnique()
				doneCnt++
				if doneCnt%10000 == 0 {
					fmt.Println(doneCnt)
				}
			}).Commit()
		}
	}()

	for {
		select {
		case f := <-b.AfterQ():
			f()
			if doneCnt == limit {
				return
			}
		}
	}
}

func TestBunchRedisBenchmark(t *testing.T) {
	// 监听本地的 6379 端口
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server started, listening on localhost:6379")

	b := New()

	go func() {
		for {
			// 接受客户端连接
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err.Error())
				continue
			}

			// 处理客户端请求
			go handleRequest(conn, b)
		}
	}()

	for {
		select {
		case f := <-b.AfterQ():
			f()
		}
	}
}

func handleRequest(conn net.Conn, b *Bunch) {
	defer conn.Close()

	// 创建一个读取器，用于读取客户端发送的数据
	reader := bufio.NewReader(conn)

	for {
		// 读取客户端发送的命令
		cmd, err := reader.ReadString('\n')
		if err != nil {
			//fmt.Println("Error reading command:", err.Error())
			return
		}

		// 如果命令是 "ping"，则回复 "PONG"
		if cmd == "PING\r\n" {
			b.NewCalls().Call(func() (interface{}, error) {
				//_, err := rdb.Ping(context.TODO()).Result()
				//if err != nil {
				//	panic(err)
				//}
				return nil, nil
			}).After(func(i interface{}, e error) {
				conn.Write([]byte("+PONG\r\n"))
			}).Commit()
		}
	}
}
