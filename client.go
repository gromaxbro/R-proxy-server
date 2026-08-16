package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", host+":80")
	if err != nil {

	}
	defer conn.Close()
	fmt.Println(conn.RemoteAddr())
	request := method + " " + location + " " + version + "\r\nHost: " + host + "\r\nConnection: close\r\n\r\n"
	conn.Write([]byte(request))
	for {
		byes := make([]byte, 1024)
		data, err := conn.Read(byes)
		if err != nil {
			break
		}
		fmt.Println(data)
		fmt.Println(string(byes))
	}
}
