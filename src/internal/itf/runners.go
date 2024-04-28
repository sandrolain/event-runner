package itf

type RunnerManager interface {
	New(EventCache, EventPlugins) (Runner, error)
	StopAll() error
}

type Runner interface {
	Ingest(chan EventMessage) (chan RunnerResult, error)
	Stop() error
}

type RunnerResult interface {
	// Setters
	SetData(any)
	AddMetadata(string, string)
	SetMetadata(string, string)
	SetConfig(string, string)
	// Getters
	HasResult() bool
	Message() EventMessage
	Destination() (string, error)
	Metadata() (map[string][]string, error)
	Config() (map[string]string, error)
	Data() (any, error)
	// Ack and Nak
	Ack() error
	Nak() error
}
