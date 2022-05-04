package client

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

// ConnectionInformation represents structure of connection information object
type ConnectionInformation struct {
	scheme      string
	host        string
	port        uint
	dataCenter  string
	accessToken string
}

// Connection represents structure of connection object
type Connection struct {
	Scheme      string
	Host        string
	Port        uint
	DataCenter  string
	AccessToken string
}

// NewConnection returns instance of new connection
func NewConnection(connection *Connection) *ConnectionInformation {
	information := &ConnectionInformation{}

	if connection.Scheme == "" {
		information.scheme = "http"
	} else {
		information.scheme = connection.Scheme
	}

	if connection.Host == "" {
		information.host = "localhost"
	} else {
		information.host = connection.Host
	}

	if connection.Port == 0 {
		information.port = 8500
	} else {
		information.port = connection.Port
	}

	if connection.DataCenter == "" {
		information.dataCenter = "dc0"
	} else {
		information.dataCenter = connection.DataCenter
	}

	if connection.AccessToken == "" {
		information.accessToken = ""
	} else {
		information.accessToken = connection.AccessToken
	}

	return information
}

// NewEmptyConnection returns instance of new empty connection
func NewEmptyConnection() *ConnectionInformation {
	return &ConnectionInformation{
		scheme:      "http",
		host:        "localhost",
		port:        8500,
		dataCenter:  "dc0",
		accessToken: "",
	}
}

// Scheme returns scheme for connection
func (information *ConnectionInformation) Scheme() string {
	return information.scheme
}

// SetScheme sets connection scheme
func (information *ConnectionInformation) SetScheme(scheme string) *ConnectionInformation {
	information.scheme = scheme
	return information
}

// Host returns host for connection
func (information *ConnectionInformation) Host() string {
	return information.host
}

// SetHost sets connection host
func (information *ConnectionInformation) SetHost(host string) *ConnectionInformation {
	information.host = host
	return information
}

// Port returns port for connection
func (information *ConnectionInformation) Port() uint {
	return information.port
}

// SetPort sets connection port
func (information *ConnectionInformation) SetPort(port uint) *ConnectionInformation {
	information.port = port
	return information
}

// DataCenter returns datacenter for connection
func (information *ConnectionInformation) DataCenter() string {
	return information.dataCenter
}

// SetDataCenter sets connection data center
func (information *ConnectionInformation) SetDataCenter(dataCenter string) *ConnectionInformation {
	information.dataCenter = dataCenter
	return information
}

// AccessToken returns access token for connection
func (information *ConnectionInformation) AccessToken() string {
	return information.accessToken
}

// SetAccessToken sets connection access token
func (information *ConnectionInformation) SetAccessToken(accessToken string) *ConnectionInformation {
	information.accessToken = accessToken
	return information
}

// UsesAccessToken indicates whether access token is being used for the connection
func (information *ConnectionInformation) UsesAccessToken() bool {
	return information.AccessToken() != ""
}

// HostPort returns host:port string for connection
func (information *ConnectionInformation) HostPort() string {
	return fmt.Sprintf("%s:%d", information.Host(), information.Port())
}

// FullPath returns full path for connection
func (information *ConnectionInformation) FullPath() string {
	return fmt.Sprintf("%s://%s", information.Scheme(), information.HostPort())
}

// IsAvailable checks if specified connection is available for use
func (information *ConnectionInformation) IsAvailable() bool {
	timeout := time.Second * 3
	connector, err := net.DialTimeout("tcp", information.HostPort(), timeout)
	if err != nil {
		return false
	}

	if connector != nil {
		defer connector.Close()
		return true
	}

	return false
}

// RoundTrip returns "ping" value between client and server
func (information *ConnectionInformation) RoundTrip() int64 {
	request, _ := http.NewRequest("GET", information.FullPath(), nil)
	var connectStart, dnsStart time.Time
	var connectEnd, dnsEnd time.Duration

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsEnd = time.Since(dnsStart)
		},
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			connectEnd = time.Since(connectStart)
		},
	}

	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(request); err != nil {
		return -1
	}
	return (connectEnd + dnsEnd).Milliseconds()
}

// IsAvailableWithRoundTrip checks whether server is available and returns round trip time
func (information *ConnectionInformation) IsAvailableWithRoundTrip() (bool, int64) {
	roundTrip := information.RoundTrip()
	return roundTrip != -1, roundTrip
}
