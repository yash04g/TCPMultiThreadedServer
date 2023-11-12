package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func do(conn net.Conn, mutex *sync.Mutex, cc *int) {
	read := make([]byte, 1024)
	_, err := conn.Read(read)
	if err != nil {
		log.Fatal(err)
	}
	mutex.Lock()
	*cc += 1
	mutex.Unlock()
	log.Println("processing the request!")
	time.Sleep(10 * time.Second) // To check server is multithreaded
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nConnection successfull!\r\n"))
	conn.Close()
	mutex.Lock()
	*cc -= 1
	mutex.Unlock()
}

func main() {
	listener, err := net.Listen("tcp", ":6900")
	if err != nil {
		log.Fatal(err)
	}
	cc := 0
	mutex := &sync.Mutex{}
	for {
		log.Println("waiting for a client to connect!")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("client connected! Active connections: %v", cc)
		go do(conn, mutex, &cc)
	}
}
