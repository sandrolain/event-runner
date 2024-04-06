package itf

type RunnerManager interface {
	New() (Runner, error)
	StopAll() error
}

type Runner interface {
	Ingest(chan EventMessage, int) (chan RunnerResult, chan error, error)
	Stop() error
}

type RunnerResult interface {
	Message() EventMessage
	Destination() (string, error)
	Metadata(string) (map[string][]string, error)
	Data() (any, error)
	Ack() error
	Nak() error
}
