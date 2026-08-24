package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type config_stk struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var config config_stk

func forward(client_conn net.Conn, host string, port int, reader *bufio.Reader, bodycount int, header string) error {
	dialer := net.Dialer{
		Timeout: 10 * time.Second,
	}
	conn, err := dialer.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		writeProxyError(client_conn, "502 Bad Gateway")
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
	dialer := net.Dialer{
		Timeout: 10 * time.Second,
	}

	target_conn, err := dialer.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		writeProxyError(conn, "502 Bad Gateway")
		return
	}
	defer target_conn.Close()

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go io.Copy(target_conn, reader)
	io.Copy(conn, target_conn)
}

func authentication(line string, auth *bool) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return
	}
	if !strings.EqualFold(parts[1], "Basic") {
		return
	}
	encoded := parts[2]
	fmt.Println(encoded)

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("Error base64 decoding")
	}
	fmt.Println(string(data))
	bigga := strings.Split(string(data), ":")
	username := bigga[0]
	password := bigga[1]
	fmt.Println(username, password)
	if username == config.Username && password == config.Password {
		fmt.Println("Authentication successful")
		*auth = true
	}
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
	auth := false
	//host := true
	for {

		line, err := reader.ReadString('\n') // read till ReadString ('\n')
		trimmedLine := strings.TrimRight(line, "\r\n")

		// An empty line signals the end of the HTTP headers section
		if strings.Contains(strings.ToLower(line), "proxy-authorization: ") {
			fmt.Println("Proxy authorization header found!!!!")
			authentication(line, &auth)
			continue
		}
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
			if len(header) != 3 {
				writeProxyError(conn, "400 Bad Request")
				return
			}

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
					path := urll.RequestURI()
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
	if !auth {
		conn.Write([]byte(
			"HTTP/1.1 407 Proxy Authentication Required\r\n" +
				"Proxy-Authenticate: Basic realm=\"MyProxy\"\r\n" +
				"Content-Length: 0\r\n" +
				"Connection: close\r\n" +
				"\r\n",
		))
		fmt.Println("Invalid protocol")
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

func writeProxyError(conn net.Conn, status string) {
	msg := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Length: 0\r\nConnection: close\r\n\r\n",
		status,
	)
	conn.Write([]byte(msg))
}

func main() {
	fmt.Println("Hello World!")

	data, err := os.ReadFile("config.json")

	json.Unmarshal(data, &config)
	//fmt.Println(config.Password, config.Username)

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
