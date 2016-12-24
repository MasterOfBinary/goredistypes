package main

import (
	"fmt"
	"net"
	"time"

	"github.com/MasterOfBinary/redistypes"
	"github.com/garyburd/redigo/redis"
)

func main() {
	netConn, err := net.Dial("tcp", "192.168.1.109:6379")
	if err != nil {
		panic(err)
	}

	conn := redis.NewConn(netConn, 10*time.Second, time.Second)
	defer conn.Close()

	list := redistypes.NewRedisList(conn, "newlist")

	size, err := list.RPush("hello", "world")
	if err != nil {
		panic(err)
	}
	fmt.Println(size)
	size, err = list.LPush("goodbye", "world")
	if err != nil {
		panic(err)
	}
	fmt.Println(size)

	r, err := list.LRange(0, -1)
	if err != nil {
		panic(err)
	}

	for _, item := range r {
		fmt.Println(string(item.([]byte)))
	}
}
