# Consul Package for Go Lang
This package provides ability to integrate Consul into the package.

As of now, this package provides the following functionality:
- **Consul Connection** - connects to single/multiple Consul instance
- **Service Registration** - allows to register application in Consul as a service
- **Key Value Watcher** - allows to watch for changes in Consul KV

## Initializing connection with Consul
There are two ways to connect application to Consul.  
You can connect to a single Consul instance, you can also specify multiple servers to connect to.

### Defining Connection Object
There are two ways to define a connection object:
- Create connection from supplied configuration
- Create empty connection and define it's values through exposed methods

**Creating connection from supplied configuration**
```go
connection := client.NewConnection(&client.Connection{
  Scheme:       "http",       // Defaults to: http
  Host:         "localhost",  // Defaults to: localhost
  Port:         8500,         // Defaults to: 8500
  DataCenter:   "dc0",        // Defaults to: dc0
  AccessToken:  "",           // Defaults to: empty string
})
```
Every single field in configuration is optional

**Creating empty connection and defining it's values through exposed methods**
```go
connection := client.NewEmptyConnection().
  SetScheme("https").
  SetHost("localhost3").
  SetPort(8502).
  SetDataCenter("dc1").
  SetAccessToken("access-token-string")
```

### Single Instance Configuration
```go
consulClient := client.SingleServer(client.NewEmptyConnection().SetHost("localhost1"))
```

### Multi Instance Configuration
```go
consulClient := client.MultipleServers([]*client.ConnectionInformation{
	client.NewEmptyConnection().SetHost("localhost1"),
	client.NewEmptyConnection().SetHost("localhost2"),
	client.NewEmptyConnection().SetHost("localhost3"),
})
```

### Global Overrides For Configured Client
Two extra methods are available if you don't want to define properties on each connection:
- **WithDataCenter** - this methods will update all connections and set their datacenter to the specified one
- **WithAccessToken** - this methods will update all connection and set their access token to the specified one

```go
consulClient.WithDataCenter("dc21").WithAccessToken("access-token-string")
```

### Establishing Connection To Consul
After client is configured, you can finally establish connection with Consul server(s).  
In order to do so, you need to call `Connect` methods on Consul client:
```go
consulClient = consulClient.Connect()
```
After connection is established, you can start using Consul.  
Right now, `consulClient` still exposed methods from this library, to get the `ConsulAPI` client, you need to access it through the `APIClient` method:
```go
apiClient := consulClient.APIClient()
```
This client exposes everything offered by the official Consul library by Hashicorp.

## Registering application as Consul service
To register application as a Consul service, you first of all need to obtain instance of Consul connection.  
After it is done, the process of client registration is pretty straight forward:
```go
consulService := service.NewService(service.Options{
  Client:       consulClient,
  Name:         "ServiceName",
})
```

This is the list of all available options:
```go
type Options struct {
	Client          *client.Client        // Consul client instance (not the API client)
	Name            string                // Name of the service without spaces
	ExtraName       string                // Additional name (in case you run two instances of the same service on the same host)
	ExtraMeta       map[string]string     // Extra meta data to be added to the service metadata defined by default
	Scheme          string                // HTTP Service: Scheme used to access this service via HTTP
	Host            string                // HTTP Service: Host used to access this service via HTTP
	Port            uint                  // HTTP Service: Port used to access this service via HTTP
	HttpServer      bool                  // HTTP Service: Indicates whether internal HTTP service is enabled
	Tags            []string              // Tags for the service
	DeregisterAfter time.Duration         // Service deregistration time (in case of critical failure) | Defaults to 1 minute
	Interval        time.Duration         // Service health check interval | Defaults to 10 seconds
	Timeout         time.Duration         // Service health check timeout | Defaults to 30 seconds
}
```

### Packaged HTTP Server
This package provides internal HTTP server which is used to provide ***Health Check over HTTP*** functionality.
If you would like to use bundeled HTTP server, simply pass `HttpServer: true` in options.  
If you are planning to use your own implementation, then you can ignore this option.

However, package will still register health check route available at `scheme://host:port/health`, so it will be picked up by your implementation.

## Watching for changes in Consul KV
Watcher provides ability to watch for changes in the the Consul KV Storage.  
In order to instantiate it, you will need two channels, `errorChannel` and `updateChannel`:
- **errorChannel** - returns any errors which occurred during the processing
- **updateChannel** - returns results of any detected change

### Creating Watcher Instance
```go
updateChannel := make(chan consulAPI.KVPairs)
errorChannel := make(chan error)

consulWatcher := &watcher.Watcher{
  Client:            consulClient.APIClient(),    // Instance of API Client
  Prefix:            "test",                      // Prefix we want to monitor changes in
  UpdateChannel:     updateChannel,               // Channel used to receive changes
  ErrorChannel:      errorChannel,                // Channel used to receive errors
}
```

### Running Watcher
After watcher has been configured, you need to start it and now you are ready to receive updates from Consul KV Storage
```go
go consulWatcher.Start()
defer consulWatcher.Stop()

for {
  select {
    case values := <-updateChannel:
      for _, value := range values {
        fmt.Printf("%v\n", value)
      }
    case err := <-errorChannel:
      fmt.Printf("%s\n", err.Error())
  }
}
```