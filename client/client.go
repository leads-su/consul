package client

import (
	"os"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/leads-su/broker"
	"github.com/leads-su/consul/state"
	"github.com/leads-su/logger"
)

// PingedServer holds information about pinged server
type PingedServer struct {
	server *ConnectionInformation
	rtt    int64
}

// Client represents structure of client
type Client struct {
	broker    *broker.Broker
	channel   chan interface{}
	servers   []*ConnectionInformation
	server    *ConnectionInformation
	apiClient *consulAPI.Client
	apiConfig *consulAPI.Config
}

// WithCustomBroker initialize client with custom broker
func WithCustomBroker(brk *broker.Broker, channel chan interface{}) *Client {
	return &Client{
		broker:  brk,
		channel: channel,
	}
}

// SingleServer defines single Consul server to connect to (when custom broker is specified)
func (client *Client) SingleServer(connection *ConnectionInformation) *Client {
	client.servers = []*ConnectionInformation{
		connection,
	}
	client.Broker().Publish(state.ConsulConfigurationPending)
	return client
}

// MultipleServers defines multiple Consul serves to connect to (when custom broker is specified)
func (client *Client) MultipleServers(connections []*ConnectionInformation) *Client {
	client.servers = connections
	client.Broker().Publish(state.ConsulConfigurationPending)
	return client
}

// SingleServer defines single Consul server to connect to
func SingleServer(connection *ConnectionInformation) *Client {
	client := &Client{
		servers: []*ConnectionInformation{
			connection,
		},
		server: connection,
	}
	client.broker, client.channel = client.initializeBroker()
	client.Broker().Publish(state.ConsulConfigurationPending)
	return client
}

// MultipleServers defines multiple Consul serves to connect to
func MultipleServers(connections []*ConnectionInformation) *Client {
	client := &Client{
		servers: connections,
	}

	client.broker, client.channel = client.initializeBroker()
	client.Broker().Publish(state.ConsulConfigurationPending)
	return client
}

// WithAccessToken sets access token for all servers
func (client *Client) WithAccessToken(accessToken string) *Client {
	for _, server := range client.servers {
		server.accessToken = accessToken
	}
	client.server.accessToken = accessToken
	return client
}

// WithDataCenter sets datacenter for all servers
func (client *Client) WithDataCenter(dataCenter string) *Client {
	for _, server := range client.servers {
		server.dataCenter = dataCenter
	}
	client.server.dataCenter = dataCenter
	return client
}

// Connect connect to best (available) Consul server
func (client *Client) Connect() *Client {
	client.Broker().Publish(state.ConsulStarting)
	client.configureAPIClient()
	apiClient, err := consulAPI.NewClient(client.apiConfig)
	if err != nil {
		logger.Fatalf(
			"consul:client",
			"failed to initialize connection to %s (datacenter: %s)",
			client.server.HostPort(),
			client.server.DataCenter(),
		)
		return nil
	}
	logger.Infof(
		"consul:client",
		"connecting to %s (datacenter: %s)",
		client.server.HostPort(),
		client.server.DataCenter(),
	)
	client.apiClient = apiClient
	client.Broker().Publish(state.ConsulStarted)
	return client
}

// Disconnect disconnects client from Consul server
func (client *Client) Disconnect() *Client {
	client.Broker().Publish(state.ConsulShuttingDown)
	client.server = nil
	client.apiClient = nil
	client.apiConfig = nil
	return client
}

// APIClient returns Consul API client
func (client *Client) APIClient() *consulAPI.Client {
	return client.apiClient
}

// Broker returns instance of Broker
func (client *Client) Broker() *broker.Broker {
	return client.broker
}

// Channel returns Broker channel for package
func (client *Client) Channel() chan interface{} {
	return client.channel
}

// IsSingleServer returns true if there is only one server specified
func (client *Client) IsSingleServer() bool {
	return len(client.servers) == 1
}

// SelectBestServer selects best server to connect to (simple and dumb, the first one available)
func (client *Client) SelectBestServer() *Client {
	var pingedServers []PingedServer
	var bestServer *ConnectionInformation = nil
	var previousPing int64 = 999

	for _, server := range client.servers {
		available, rtt := server.IsAvailableWithRoundTrip()
		if available {
			pingedServers = append(pingedServers, PingedServer{
				server: server,
				rtt:    rtt,
			})
		} else {
			logger.Warnf("consul:client", "server %s is not available for connection", server.HostPort())
		}
	}

	serversCount := len(pingedServers)

	if serversCount == 0 {
		logger.Fatalf("consul:client", "there are no alive consul servers available to connect to")
		os.Exit(1)
	}

	for _, entity := range pingedServers {
		if entity.rtt < previousPing {
			bestServer = entity.server
			previousPing = entity.rtt
		}
	}

	logger.Infof("consul:client", "selecting %s as a target server with ping of %dms", bestServer.HostPort(), previousPing)

	client.server = bestServer
	client.Broker().Publish(state.ConsulConfigured)
	return client
}

// configureAPIClient configures Consul API client
func (client *Client) configureAPIClient() {
	clientConfiguration := consulAPI.DefaultConfig()
	clientConfiguration.Scheme = client.server.Scheme()
	clientConfiguration.Address = client.server.HostPort()
	clientConfiguration.Datacenter = client.server.DataCenter()

	if client.server.UsesAccessToken() {
		clientConfiguration.Token = client.server.AccessToken()
	}

	client.apiConfig = clientConfiguration
}

// initializeBroker creates new instance of messages broker
func (client *Client) initializeBroker() (*broker.Broker, chan interface{}) {
	instance := broker.NewBroker()
	go instance.Start()
	channel := instance.Subscribe()
	return instance, channel
}
