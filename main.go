package main

import (
	"fmt"
	"net"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/MasterOfBinary/redistypes"
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
