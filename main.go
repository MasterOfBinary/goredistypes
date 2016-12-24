package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/MasterOfBinary/redistypes/hyperloglog"
	"github.com/MasterOfBinary/redistypes/list"
	"github.com/garyburd/redigo/redis"
)

type config struct {
	HostAndPort string `json:"hostAndPort"`
}

var defaultConfig = config{"127.0.0.1:6379"}

func loadConfig() config {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Errorf("Falling back to default config")
		return defaultConfig
	}

	var c config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		fmt.Errorf("Falling back to default config")
		return defaultConfig
	}

	return c
}

func main() {
	conf := loadConfig()
	netConn, err := net.Dial("tcp", conf.HostAndPort)
	if err != nil {
		panic(err)
	}

	conn := redis.NewConn(netConn, time.Second, time.Second)
	defer conn.Close()

	l := list.NewRedisList(conn, "newlist")

	size, err := l.RPush("hello", "world")
	if err != nil {
		panic(err)
	}
	fmt.Println(size)
	size, err = l.LPush("world", "hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(size)

	r, err := redis.Strings(l.LRange(0, -1))
	if err != nil {
		panic(err)
	}

	fmt.Println(r)

	hll := hyperloglog.NewRedisHyperLogLog(conn, "hll")

	modified, _ := hll.Add("hello", "goodbye", "abc", "def")
	size, _ = hll.Count()
	fmt.Println(size)

	hll2 := hyperloglog.NewRedisHyperLogLog(conn, "hll2")

	modified, _ = hll2.Add("yay", "abc", "boo", "hi", "hey")
	fmt.Println(modified)
	size, _ = hll.Count()
	fmt.Println(size)

	hll3, err := hll.Merge("hll3", hll2)
	if err != nil {
		panic(err)
	}

	size, _ = hll3.Count()
	fmt.Println(size)
}
