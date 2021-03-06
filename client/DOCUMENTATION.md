<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# client

```go
import "github.com/leads-su/consul/client"
```

## Index

- [type Client](<#type-client>)
  - [func MultipleServers(connections []*ConnectionInformation) *Client](<#func-multipleservers>)
  - [func SingleServer(connection *ConnectionInformation) *Client](<#func-singleserver>)
  - [func WithCustomBroker(brk *broker.Broker, channel chan interface{}) *Client](<#func-withcustombroker>)
  - [func (client *Client) APIClient() *consulAPI.Client](<#func-client-apiclient>)
  - [func (client *Client) Broker() *broker.Broker](<#func-client-broker>)
  - [func (client *Client) Channel() chan interface{}](<#func-client-channel>)
  - [func (client *Client) Connect() *Client](<#func-client-connect>)
  - [func (client *Client) Disconnect() *Client](<#func-client-disconnect>)
  - [func (client *Client) IsSingleServer() bool](<#func-client-issingleserver>)
  - [func (client *Client) MultipleServers(connections []*ConnectionInformation) *Client](<#func-client-multipleservers>)
  - [func (client *Client) SelectBestServer() *Client](<#func-client-selectbestserver>)
  - [func (client *Client) SingleServer(connection *ConnectionInformation) *Client](<#func-client-singleserver>)
  - [func (client *Client) WithAccessToken(accessToken string) *Client](<#func-client-withaccesstoken>)
  - [func (client *Client) WithDataCenter(dataCenter string) *Client](<#func-client-withdatacenter>)
- [type Connection](<#type-connection>)
- [type ConnectionInformation](<#type-connectioninformation>)
  - [func NewConnection(connection *Connection) *ConnectionInformation](<#func-newconnection>)
  - [func NewEmptyConnection() *ConnectionInformation](<#func-newemptyconnection>)
  - [func (information *ConnectionInformation) AccessToken() string](<#func-connectioninformation-accesstoken>)
  - [func (information *ConnectionInformation) DataCenter() string](<#func-connectioninformation-datacenter>)
  - [func (information *ConnectionInformation) FullPath() string](<#func-connectioninformation-fullpath>)
  - [func (information *ConnectionInformation) Host() string](<#func-connectioninformation-host>)
  - [func (information *ConnectionInformation) HostPort() string](<#func-connectioninformation-hostport>)
  - [func (information *ConnectionInformation) IsAvailable() bool](<#func-connectioninformation-isavailable>)
  - [func (information *ConnectionInformation) IsAvailableWithRoundTrip() (bool, int64)](<#func-connectioninformation-isavailablewithroundtrip>)
  - [func (information *ConnectionInformation) Port() uint](<#func-connectioninformation-port>)
  - [func (information *ConnectionInformation) RoundTrip() int64](<#func-connectioninformation-roundtrip>)
  - [func (information *ConnectionInformation) Scheme() string](<#func-connectioninformation-scheme>)
  - [func (information *ConnectionInformation) SetAccessToken(accessToken string) *ConnectionInformation](<#func-connectioninformation-setaccesstoken>)
  - [func (information *ConnectionInformation) SetDataCenter(dataCenter string) *ConnectionInformation](<#func-connectioninformation-setdatacenter>)
  - [func (information *ConnectionInformation) SetHost(host string) *ConnectionInformation](<#func-connectioninformation-sethost>)
  - [func (information *ConnectionInformation) SetPort(port uint) *ConnectionInformation](<#func-connectioninformation-setport>)
  - [func (information *ConnectionInformation) SetScheme(scheme string) *ConnectionInformation](<#func-connectioninformation-setscheme>)
  - [func (information *ConnectionInformation) UsesAccessToken() bool](<#func-connectioninformation-usesaccesstoken>)


## type Client

Client represents structure of client

```go
type Client struct {
    // contains filtered or unexported fields
}
```

### func MultipleServers

```go
func MultipleServers(connections []*ConnectionInformation) *Client
```

MultipleServers defines multiple Consul serves to connect to

### func SingleServer

```go
func SingleServer(connection *ConnectionInformation) *Client
```

SingleServer defines single Consul server to connect to

### func WithCustomBroker

```go
func WithCustomBroker(brk *broker.Broker, channel chan interface{}) *Client
```

WithCustomBroker initialize client with custom broker

### func \(\*Client\) APIClient

```go
func (client *Client) APIClient() *consulAPI.Client
```

APIClient returns Consul API client

### func \(\*Client\) Broker

```go
func (client *Client) Broker() *broker.Broker
```

Broker returns instance of Broker

### func \(\*Client\) Channel

```go
func (client *Client) Channel() chan interface{}
```

Channel returns Broker channel for package

### func \(\*Client\) Connect

```go
func (client *Client) Connect() *Client
```

Connect connect to best \(available\) Consul server

### func \(\*Client\) Disconnect

```go
func (client *Client) Disconnect() *Client
```

Disconnect disconnects client from Consul server

### func \(\*Client\) IsSingleServer

```go
func (client *Client) IsSingleServer() bool
```

IsSingleServer returns true if there is only one server specified

### func \(\*Client\) MultipleServers

```go
func (client *Client) MultipleServers(connections []*ConnectionInformation) *Client
```

MultipleServers defines multiple Consul serves to connect to \(when custom broker is specified\)

### func \(\*Client\) SelectBestServer

```go
func (client *Client) SelectBestServer() *Client
```

SelectBestServer selects best server to connect to \(simple and dumb\, the first one available\)

### func \(\*Client\) SingleServer

```go
func (client *Client) SingleServer(connection *ConnectionInformation) *Client
```

SingleServer defines single Consul server to connect to \(when custom broker is specified\)

### func \(\*Client\) WithAccessToken

```go
func (client *Client) WithAccessToken(accessToken string) *Client
```

WithAccessToken sets access token for all servers

### func \(\*Client\) WithDataCenter

```go
func (client *Client) WithDataCenter(dataCenter string) *Client
```

WithDataCenter sets datacenter for all servers

## type Connection

Connection represents structure of connection object

```go
type Connection struct {
    Scheme      string
    Host        string
    Port        uint
    DataCenter  string
    AccessToken string
}
```

## type ConnectionInformation

ConnectionInformation represents structure of connection information object

```go
type ConnectionInformation struct {
    // contains filtered or unexported fields
}
```

### func NewConnection

```go
func NewConnection(connection *Connection) *ConnectionInformation
```

NewConnection returns instance of new connection

### func NewEmptyConnection

```go
func NewEmptyConnection() *ConnectionInformation
```

NewEmptyConnection returns instance of new empty connection

### func \(\*ConnectionInformation\) AccessToken

```go
func (information *ConnectionInformation) AccessToken() string
```

AccessToken returns access token for connection

### func \(\*ConnectionInformation\) DataCenter

```go
func (information *ConnectionInformation) DataCenter() string
```

DataCenter returns datacenter for connection

### func \(\*ConnectionInformation\) FullPath

```go
func (information *ConnectionInformation) FullPath() string
```

FullPath returns full path for connection

### func \(\*ConnectionInformation\) Host

```go
func (information *ConnectionInformation) Host() string
```

Host returns host for connection

### func \(\*ConnectionInformation\) HostPort

```go
func (information *ConnectionInformation) HostPort() string
```

HostPort returns host:port string for connection

### func \(\*ConnectionInformation\) IsAvailable

```go
func (information *ConnectionInformation) IsAvailable() bool
```

IsAvailable checks if specified connection is available for use

### func \(\*ConnectionInformation\) IsAvailableWithRoundTrip

```go
func (information *ConnectionInformation) IsAvailableWithRoundTrip() (bool, int64)
```

IsAvailableWithRoundTrip checks whether server is available and returns round trip time

### func \(\*ConnectionInformation\) Port

```go
func (information *ConnectionInformation) Port() uint
```

Port returns port for connection

### func \(\*ConnectionInformation\) RoundTrip

```go
func (information *ConnectionInformation) RoundTrip() int64
```

RoundTrip returns "ping" value between client and server

### func \(\*ConnectionInformation\) Scheme

```go
func (information *ConnectionInformation) Scheme() string
```

Scheme returns scheme for connection

### func \(\*ConnectionInformation\) SetAccessToken

```go
func (information *ConnectionInformation) SetAccessToken(accessToken string) *ConnectionInformation
```

SetAccessToken sets connection access token

### func \(\*ConnectionInformation\) SetDataCenter

```go
func (information *ConnectionInformation) SetDataCenter(dataCenter string) *ConnectionInformation
```

SetDataCenter sets connection data center

### func \(\*ConnectionInformation\) SetHost

```go
func (information *ConnectionInformation) SetHost(host string) *ConnectionInformation
```

SetHost sets connection host

### func \(\*ConnectionInformation\) SetPort

```go
func (information *ConnectionInformation) SetPort(port uint) *ConnectionInformation
```

SetPort sets connection port

### func \(\*ConnectionInformation\) SetScheme

```go
func (information *ConnectionInformation) SetScheme(scheme string) *ConnectionInformation
```

SetScheme sets connection scheme

### func \(\*ConnectionInformation\) UsesAccessToken

```go
func (information *ConnectionInformation) UsesAccessToken() bool
```

UsesAccessToken indicates whether access token is being used for the connection



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
