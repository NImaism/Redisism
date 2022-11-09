package Redisism

import (
	"fmt"
	"sync"
)

type Cache struct{}

type Request struct {
	Request string
	Key     string
	Val     string
	RTL     int
}

type Response struct {
	Val string
	Msg string
}

type ServerConnection struct {
	snd     chan Request
	rcv     chan Response
	storage map[string]string
	owner   string
	Mutex   sync.Mutex
}

type ClientConnection struct {
	snd chan Request
	rcv chan Response
}

func InitCache() *Cache {
	return &Cache{}
}

func (cache *Cache) Connect(name string) *ClientConnection {
	request := make(chan Request)
	response := make(chan Response)
	storage := make(map[string]string)

	serverConn := &ServerConnection{
		snd:     request,
		rcv:     response,
		owner:   name,
		storage: storage,
	}

	clientConn := &ClientConnection{
		snd: request,
		rcv: response,
	}

	fmt.Println("[Redisism] New Connection established \n", "Name : "+name)

	go func(serverConn *ServerConnection) {
		for {
			Data := <-serverConn.snd
			if Data.Request == "Set" {

				serverConn.Mutex.Lock()

				serverConn.storage[Data.Key] = Data.Val
				serverConn.rcv <- Response{
					Data.Val,
					"Done",
				}
				serverConn.Mutex.Unlock()

			} else if Data.Request == "Get" {
				serverConn.Mutex.Lock()

				D, HaveData := serverConn.storage[Data.Key]
				if !HaveData {
					serverConn.rcv <- Response{
						"",
						"No Data",
					}
				} else {
					serverConn.rcv <- Response{
						D,
						"Done",
					}
				}

				serverConn.Mutex.Unlock()
			}
		}
	}(serverConn)

	return clientConn
}

func (conn *ClientConnection) Set(key string, val string) {
	conn.snd <- Request{
		Request: "Set",
		Key:     key,
		Val:     val,
	}
	<-conn.rcv
}

func (conn *ClientConnection) Get(key string) string {
	conn.snd <- Request{
		Request: "Get",
		Key:     key,
	}
	res := <-conn.rcv

	return res.Val
}
