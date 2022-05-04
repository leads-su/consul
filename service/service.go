package service

import (
	"fmt"
	"os"
	"runtime"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/leads-su/consul/client"
	"github.com/leads-su/consul/http"
	"github.com/leads-su/consul/state"
	"github.com/leads-su/logger"
	"github.com/leads-su/version"
)

// Service represents structure of service configuration
type Service struct {
	client     *client.Client
	name       string
	extraName  string
	extraMeta  map[string]string
	tags       []string
	scheme     string
	host       string
	port       uint
	httpServer bool

	deregisterChannel chan bool
	deregisterAfter   time.Duration
	interval          time.Duration
	timeout           time.Duration
}

// Options represents structure of service options
type Options struct {
	Client          *client.Client
	Name            string
	ExtraName       string
	ExtraMeta       map[string]string
	Scheme          string
	Host            string
	Port            uint
	HttpServer      bool
	Tags            []string
	DeregisterAfter time.Duration
	Interval        time.Duration
	Timeout         time.Duration
}

// NewService creates new instance of Consul service
func NewService(options Options) *Service {
	options.Client.Broker().Publish(state.ConsulCreatingService)
	http.NewServer(options.Port, options.HttpServer)

	service := &Service{
		name:            options.Name,
		extraName:       options.ExtraName,
		extraMeta:       options.ExtraMeta,
		tags:            options.Tags,
		scheme:          options.Scheme,
		host:            options.Host,
		port:            options.Port,
		client:          options.Client,
		httpServer:      options.HttpServer,
		deregisterAfter: options.DeregisterAfter,
		interval:        options.Interval,
		timeout:         options.Timeout,
	}

	if service.scheme == "" {
		service.scheme = "http"
	}

	if service.host == "" {
		service.host = "127.0.0.1"
	}

	if service.deregisterAfter == 0 {
		service.deregisterAfter = time.Duration(1) * time.Minute
	}

	if service.interval == 0 {
		service.interval = time.Duration(10) * time.Second
	}

	if service.timeout == 0 {
		service.timeout = time.Duration(30) * time.Second
	}

	options.Client.Broker().Publish(state.ConsulServiceCreated)
	return service
}

// HostPort returns service host:port string
func (service *Service) HostPort() string {
	return fmt.Sprintf("%s:%d", service.host, service.port)
}

// FullPath returns service full path
func (service *Service) FullPath() string {
	return fmt.Sprintf("%s://%s", service.scheme, service.HostPort())
}

// Register registers service in Consul
func (service *Service) Register() error {
	configuration, err := service.buildServiceConfiguration()
	if err != nil {
		return err
	}
	service.deregisterChannel = service.register(configuration)
	service.client.Broker().Publish(state.ConsulServiceRegistered)
	return nil
}

// Deregister deregisters service from Consul
func (service *Service) Deregister() error {
	if service.deregisterChannel == nil {
		logger.Warnf("consul:service", "this service is not registered in consul")
		return nil
	}
	service.deregisterChannel <- true
	<-service.deregisterChannel
	service.client.Broker().Publish(state.ConsulServiceDeregistered)
	return nil
}

// register handles de/registration process
func (service *Service) register(registration *consulAPI.AgentServiceRegistration) chan bool {
	registered := func(serviceID string) bool {
		if serviceID == "" {
			return false
		}
		services, err := service.client.APIClient().Agent().Services()
		if err != nil {
			logger.Errorf("consul:service", "cannot retrieve list of services - %s", err.Error())
			service.client.Broker().Publish(state.ConsulRestartRequested)
			return false
		}
		return services[serviceID] != nil
	}

	register := func() string {
		if err := service.client.APIClient().Agent().ServiceRegister(registration); err != nil {
			logger.Errorf("consul:service", "failed to register service `%s` in consul - %s", registration.Name, err.Error())
			service.client.Broker().Publish(state.ConsulRestartRequested)
			return ""
		}

		logger.Tracef("consul:service", "registered `%s` with id `%s` at %s",
			registration.Name,
			registration.ID,
			registration.Address,
		)

		for _, check := range registration.Checks {
			logger.Tracef("consul:service", "registered check `%s` for service `%s`", check.CheckID, registration.ID)
		}
		return registration.ID
	}

	deregister := func(serviceID string) {
		logger.Tracef("consul:service", "de-registering service `%s` from consul", registration.Name)
		err := service.client.APIClient().Agent().ServiceDeregister(serviceID)
		if err != nil {
			logger.Errorf("consul:service", "Failed to deregister service - %s", err.Error())
		}
	}

	passTTL := func(serviceTTLCheckID string) {
		err := service.client.APIClient().Agent().UpdateTTL(serviceTTLCheckID, time.Now().UTC().Format(time.RFC3339), consulAPI.HealthPassing)
		if err != nil {
			logger.Errorf("consul:service", "Unable to pass TTL check for service with ID - %s", serviceTTLCheckID)
		}
	}

	deregisterChannel := make(chan bool)

	go func() {
		var serviceID string
		var serviceTTLCheckID string

		for {
			if !registered(serviceID) {
				serviceID = register()
				serviceTTLCheckID = computeServiceTTLCheckID(serviceID)
				passTTL(serviceTTLCheckID)
			}
			select {
			case <-deregisterChannel:
				deregister(serviceID)
				deregisterChannel <- true
				return
			case <-time.After(service.interval):
				passTTL(serviceTTLCheckID)
			}
		}
	}()
	return deregisterChannel
}

// generateServiceID generates service ID from given data
func (service *Service) generateServiceID() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	serviceID := fmt.Sprintf("%s-%s-%s", service.name, hostname, service.host)
	if len(service.extraName) > 0 {
		serviceID = fmt.Sprintf("%s-%s-%s-%s", service.name, service.extraName, hostname, service.host)
	}
	return serviceID, nil
}

// buildServiceConfiguration builds and returns service configuration information, or error
func (service *Service) buildServiceConfiguration() (*consulAPI.AgentServiceRegistration, error) {
	serviceID, err := service.generateServiceID()
	if err != nil {
		return nil, err
	}

	var serviceHealthChecks []*consulAPI.AgentServiceCheck
	serviceMeta := map[string]string{
		"application_build_date":   version.GetBuildDate(),
		"application_build_commit": version.GetCommit(),
		"application_version":      version.GetVersion(),
		"architecture":             runtime.GOARCH,
		"go_version":               runtime.Version(),
		"operating_system":         runtime.GOOS,
	}

	for key, value := range service.extraMeta {
		serviceMeta[key] = value
	}

	serviceHealthChecks = append(serviceHealthChecks, &consulAPI.AgentServiceCheck{
		CheckID:                        computeServiceTTLCheckID(serviceID),
		TTL:                            (service.interval + time.Duration(5)).String(),
		DeregisterCriticalServiceAfter: service.deregisterAfter.String(),
	})

	if service.httpServer && runtime.GOOS != "windows" {
		checkURL := service.FullPath() + "/health"
		serviceHealthChecks = append(serviceHealthChecks, &consulAPI.AgentServiceCheck{
			CheckID:  computeServiceHttpCheckID(serviceID),
			HTTP:     checkURL,
			Interval: service.interval.String(),
			Timeout:  service.timeout.String(),
		})
	}
	return &consulAPI.AgentServiceRegistration{
		ID:      serviceID,
		Name:    service.name,
		Address: service.host,
		Port:    int(service.port),
		Meta:    serviceMeta,
		Tags:    service.tags,
		Checks:  serviceHealthChecks,
	}, nil
}

// computeServiceTTLCheckID generate service TTL check ID
func computeServiceTTLCheckID(serviceID string) string {
	return serviceID + "-ttl"
}

// computeServiceHttpCheckID generate service HTTP check ID
func computeServiceHttpCheckID(serviceID string) string {
	return serviceID + "-http"
}
