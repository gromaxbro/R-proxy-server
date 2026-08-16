package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func forward(host string, request string) []byte {
	conn, err := net.Dial("tcp", host+":80")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return nil
	}
	defer conn.Close()
	fmt.Println(conn.RemoteAddr())
	conn.Write([]byte(request))
	var requestbuf []byte

	for {
		byes := make([]byte, 1024)

		data, err := conn.Read(byes)
		if err != nil {
			break
		}
		//fmt.Println(data)
		//fmt.Println(string(byes[:data]))
		requestbuf = append(requestbuf, byes[:data]...)
	}
	fmt.Println(string(requestbuf))
	return requestbuf
}

func main() {
	fmt.Println("Hello World!")
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("unable to listen on port 8080")
		return
	}
	defer listner.Close()
	fmt.Println("listner started in : ", listner.Addr())

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("unable to accept connection")
			return
		}
		fmt.Println(conn.RemoteAddr().String())
		reader := bufio.NewReader(conn)
		methodfound := true

		//host := true
		for {

			line, err := reader.ReadString('\n') // read till ReadString ('\n')

			if err != nil {
				return
			}
			//fmt.Printf("%q\n", line)
			if methodfound {
				header := strings.Fields(line)
				method := header[0]
				location := header[1]
				urll, err := url.Parse(location)
				if err != nil {
				}
				fmt.Println(urll.Path)
				host := urll.Host
				path := urll.Path
				version := header[2]
				fmt.Println("method:", method, "location:", location, "version:", version, host, path)
				methodfound = false

				request := method + " " + path + " " + version + "\r\nHost: " + host + "\r\nConnection: close\r\n\r\n"

				requestbuf := forward(host, request)
				conn.Write(requestbuf)

			}
			//if strings.Contains(strings.ToLower(line), "host: ") {
			//	host := strings.Split(strings.ToLower(line), "host: ")
			//	fmt.Println(host[1])
			//}

		}

	}

}
