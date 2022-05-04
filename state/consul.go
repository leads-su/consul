package state

const (
	ConsulConfigurationPending = iota + consulIotaValue
	ConsulConfigured

	ConsulCreatingService
	ConsulServiceCreated
	ConsulServiceRegistered
	ConsulServiceDeregistered

	ConsulStarting
	ConsulStarted
	ConsulShuttingDown
	ConsulRestartRequested
)
