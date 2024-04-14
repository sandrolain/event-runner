package config

type Config struct {
	Logger      Logger
	Connections []Connection
	Runners     []Runner
	Lines       []Line
	Inputs      []Input
	Outputs     []Output
}

type Logger struct {
	Level  string `json:"level" default:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	Format string `json:"format" default:"TEXT" validate:"required,oneof=TEXT JSON"`
	Color  bool   `json:"color"`
}

type Connection struct {
	ID       string `json:"id" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=nats redis http"`
	Hostname string `json:"hostname" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Runner struct {
	ID          string `json:"id" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=es5 risor"`
	ProgramPath string `json:"program_path" validate:"required,file"`
	ProgramB64  string `json:"program_b64" validate:"required,base64"`
	Buffer      int    `json:"buffer"`
}

type Line struct {
	ID         string `json:"id" validate:"required"`
	Connection string `json:"connection" validate:"required"`
	Input      string `json:"input" validate:"required"`
	Runner     string `json:"runner" validate:"required"`
	Output     string `json:"output" validate:"required"`
}

type Input struct {
	ID     string `json:"id" validate:"required"`
	Topic  string `json:"topic" validate:"required"`
	Stream string `json:"stream"`
	Client string `json:"client"`
	Buffer int    `json:"buffer"`
}

type Output struct {
	ID     string `json:"id" validate:"required"`
	Topic  string `json:"topic" validate:"required"`
	Method string `json:"method" validate:"oneof=POST PUT PATCH"`
	Stream string `json:"stream"`
	Client string `json:"client"`
}
