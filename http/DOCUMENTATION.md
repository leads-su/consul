<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# http

```go
import "github.com/leads-su/consul/http"
```

## Index

- [type Server](<#type-server>)
  - [func NewServer(port uint, enabled bool) *Server](<#func-newserver>)


## type Server

Server represents structure of HTTP server

```go
type Server struct {
    Enabled bool
    Port    uint
}
```

### func NewServer

```go
func NewServer(port uint, enabled bool) *Server
```

NewServer creates new instance of HTTP server



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
