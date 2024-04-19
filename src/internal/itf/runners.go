package itf

type RunnerManager interface {
	New(EventCache) (Runner, error)
	StopAll() error
}

type Runner interface {
	Ingest(chan EventMessage) (chan RunnerResult, error)
	Stop() error
}

type RunnerResult interface {
	Message() EventMessage
	Destination() (string, error)
	Metadata() (map[string][]string, error)
	Config() (map[string]string, error)
	Data() (any, error)
	Ack() error
	Nak() error
}
