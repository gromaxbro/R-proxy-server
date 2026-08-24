# Go Forward Proxy

A  HTTP/HTTPS forward proxy built from scratch in Go. The proxy accepts client connections, authenticates users, forwards HTTP requests, and establishes TCP tunnels for HTTPS using the HTTP `CONNECT` method.

blog :- https://blog.brocue.online/post/proxy-server

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
```

## HTTP Proxy

For normal HTTP requests, the client sends the proxy an absolute URL:
```
GET http://example.com/ HTTP/1.1
Host: example.com
```
The proxy extracts the destination and forwards the request to the target server. 

## HTTPS Proxying

HTTPS uses the CONNECT method:
```
CONNECT example.com:443 HTTP/1.1
Host: example.com:443
```

## Authentication

The proxy uses the Proxy-Authorization header.

Example:
```
Proxy-Authorization: Basic <base64 credentials>
```

## Configuration

Configuration is stored in config.json.

Example:
```
{
  "username": "admin",
  "password": "secret123",
  "blacklist": [
    "youtube.com",
    "twitch.tv",
    "example.com"
  ]
}
```

## Logging

Important proxy events can be written to proxy.log.

Example:
```
2026-08-25 20:31:42 | admin | 192.168.1.5 | twitch.tv:443 | HTTPS | ALLOWED
2026-08-25 20:31:45 | admin | 192.168.1.5 | youtube.com:443 | HTTPS | BLOCKED
2026-08-25 20:32:01 | unknown | 192.168.1.5 | - | AUTH_FAILED
```
