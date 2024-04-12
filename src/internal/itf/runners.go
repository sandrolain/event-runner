package itf

type RunnerManager interface {
	New() (Runner, error)
	StopAll() error
}

type Runner interface {
	Ingest(chan EventMessage, int) (chan RunnerResult, error)
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
