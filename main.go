package main

import (
	"flag"
	"fmt"
	// "io"
	"bufio"
	"net"
	"os"
	"strings"
)

var (
	port        = flag.String("port", "6532", "服务端口号")
	datadir     = flag.String("datadir", "data", "数据文件夹, 文件名为{datadir}/{day}.blf, day=1970/01/01到今天的天数.")
	item_number = flag.Uint("number", 1<<31, "布隆过滤器数量, 通常是总条数的20倍.")
	days        = flag.Int("days", 0, "去重最多存储天数")
	uniq        *Uniq
)

func main() {
	flag.Parse()
	var l net.Listener
	var err error

	if *days < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	uniq, err = Open(*datadir, *item_number, *days)
	if err != nil {
		fmt.Println("open error", err)
	}

	l, err = net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on :" + *port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic info is: ", err)
		}
	}()
	defer conn.Close()

	ioc := bufio.NewReader(conn)
	for {
		line, _, err := ioc.ReadLine()
		if err != nil {
			return
		} else {
			cmd := strings.Split(string(line), " ")
			if cmd[0] == "get" {
				if len(cmd) > 1 && uniq.TestAndAdd([]byte(cmd[1]), 0) {
					conn.Write([]byte("VALUE " + cmd[1] + " 0 4\r\nTRUE\r\nEND\r\n"))
				} else {
					conn.Write([]byte("NOT_FOUND\r\n"))
				}
			}
		}
	}
}
