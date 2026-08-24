# Go Forward Proxy

A basic HTTP/HTTPS forward proxy built from scratch in Go.

The proxy accepts client connections, authenticates users, forwards HTTP requests, and establishes TCP tunnels for HTTPS using the HTTP `CONNECT` method.

## Features

- HTTP forward proxy
- HTTPS tunneling using `CONNECT`
- Bidirectional TCP forwarding
- Proxy authentication using `Proxy-Authorization: Basic`
- JSON configuration
- Domain blacklisting
- Connection statistics
- Authentication failure tracking
- Rejected request tracking
- Request/connection logging
- Concurrent client handling using goroutines
- Mutex-protected shared statistics and logs

## Architecture

```text
Client
  |
  | HTTP / CONNECT
  v
+----------------+
|   Go Proxy     |
|    :4241       |
+----------------+
  |
  | Authentication
  v
+----------------+
| Blacklist      |
| Check          |
+----------------+
  |
  | Allowed
  v
+----------------+
| Target Server  |
+----------------+
