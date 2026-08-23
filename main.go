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

func forward(client_conn net.Conn, host string, port int, reader *bufio.Reader, bodycount int, header string) error {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return nil
	}
	defer conn.Close()
	fmt.Println(conn.RemoteAddr())
	conn.Write([]byte(header))

	if bodycount > 0 {
		io.CopyN(conn, reader, int64(bodycount))
	}

	//response_reader := bufio.NewReader(conn)
	io.Copy(client_conn, conn)
	return nil
}
func https_coonection(host string, port int, conn net.Conn, reader *bufio.Reader) {
	target_conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer target_conn.Close()

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go io.Copy(target_conn, reader)
	io.Copy(conn, target_conn)
}
func read_client(conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	host := ""
	port := 80
	methodfound := true
	havebody := 0
	var headersBuilder strings.Builder
	is_http := true
	//host := true
	for {

		line, err := reader.ReadString('\n') // read till ReadString ('\n')
		trimmedLine := strings.TrimRight(line, "\r\n")

		// An empty line signals the end of the HTTP headers section

		if strings.Contains(strings.ToLower(line), "content-length: ") {
			linetemp := strings.TrimRight(line, "\r\n")

			// 2. Perform the split
			content := strings.Split(strings.ToLower(linetemp), "content-length: ")
			fmt.Println(strconv.Atoi(content[1]))
			havebody, err = strconv.Atoi(content[1])
		}

		if err != nil {
			break
		}
		fmt.Printf("%q\n", line)
		if methodfound {
			header := strings.Fields(line)
			method := header[0]
			location := header[1]
			version := header[2]

			methodfound = false

			if method == "CONNECT" {
				host = location
				port = 443 // default for HTTPS
				is_http = false
				if strings.Contains(host, ":") {
					h, pStr, err := net.SplitHostPort(host)
					if err == nil {
						host = h
						if parsedPort, err := strconv.Atoi(pStr); err == nil {
							port = parsedPort
						}
					}
				}
			} else { // <-- ADD THIS ELSE BLOCK
				urll, err := url.Parse(location)
				if err == nil && urll != nil {
					if urll.Host != "" {
						host = urll.Host
					}
					path := urll.Path
					if path == "" {
						path = "/"
					}
					headersBuilder.WriteString(fmt.Sprintf("%s %s %s\r\n", method, path, version))
				} else {
					headersBuilder.WriteString(fmt.Sprintf("%s %s %s\r\n", method, location, version))
				}

				if strings.Contains(host, ":") {
					h, pStr, err := net.SplitHostPort(host)
					if err == nil {
						host = h
						if parsedPort, err := strconv.Atoi(pStr); err == nil {
							port = parsedPort
						}
					}
				}
			}

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
	defer conn.Close() // <-- ADD THIS to prevent memory/socket leaks

	if is_http {
		forward(conn, host, port, reader, havebody, headersBuilder.String())
	} else {
		https_coonection(host, port, conn, reader) // <-- CALL IT HERE
	}
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
