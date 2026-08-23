package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func read_client(conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	host := ""
	port := 443
	methodfound := true
	var headersBuilder strings.Builder

	//host := true
	for {

		line, err := reader.ReadString('\n') // read till ReadString ('\n')
		trimmedLine := strings.TrimRight(line, "\r\n")

		// An empty line signals the end of the HTTP headers section

		if err != nil {
			break
		}
		fmt.Printf("%q\n", line)
		if methodfound {
			header := strings.Fields(line)
			method := header[0]
			location := header[1]
			urll, err := url.Parse(location)
			if err != nil {
			}
			fmt.Println(urll.Path)
			host = urll.Host
			path := urll.Path
			if path == "" {
				path = "/"
			}
			version := header[2]
			fmt.Println("method:", method, "location:", location, "version:", version, host, path, port)
			methodfound = false
			//port := 80

			//request = method + " " + path + " " + version + "\r\nHost: " + host + "\r\nConnection: close\r\n\r\n"
			headersBuilder.WriteString(fmt.Sprintf("%s %s %s\r\n", method, path, version))
			continue
		}

		if trimmedLine == "" {
			fmt.Println("GOT END HEADER*******")
			break
		}

		headersBuilder.WriteString(line)

	}
	headersBuilder.WriteString("Connection: close\r\n\r\n")
	fmt.Println("**********************")
	fmt.Println(headersBuilder.String())
	fmt.Println("**********************")

	target_conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("tunnel establised")

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	go io.Copy(target_conn, reader)
	io.Copy(conn, target_conn)
	//requestbuf := forward(conn, host, port, reader, havebody, headersBuilder.String())
	//if requestbuf != nil {
	//	fmt.Println("no error")
	//}
	//for i := 0; i < havebody; i++ {
	//	line, err := reader.ReadString('\n') // read till ReadString ('\n')
	//	fmt.Println("bodyyyyyyy*****8")
	//	fmt.Println(line)
	//	if err != nil {
	//		break
	//	}
	//}
}

func main() {
	fmt.Println("Hello World!")
	listner, err := net.Listen("tcp", ":4241")
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
		go read_client(conn)

	}

}
