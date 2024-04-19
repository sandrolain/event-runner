package config

type Config struct {
	Logger      Logger       `yaml:"logger" json:"logger" validate:"required"`
	Connections []Connection `yaml:"connections" json:"connections" validate:"required,dive"`
	Runners     []Runner     `yaml:"runners" json:"runners" validate:"required,dive"`
	Lines       []Line       `yaml:"lines" json:"lines" validate:"required,dive"`
	Inputs      []Input      `yaml:"inputs" json:"inputs" validate:"required,dive"`
	Outputs     []Output     `yaml:"outputs" json:"outputs" validate:"required,dive"`
	Caches      []Cache      `yaml:"caches" json:"caches" validate:"dive"`
}

type Logger struct {
	Level  string `yaml:"level" json:"level" default:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR debug info warn error"`
	Format string `yaml:"format" json:"format" default:"TEXT" validate:"required,oneof=TEXT JSON text json"`
	Color  bool   `yaml:"color" json:"color"`
}

type Connection struct {
	ID       string `yaml:"id" json:"id" validate:"required"`
	Type     string `yaml:"type" json:"type" validate:"required,oneof=nats redis http"`
	Hostname string `yaml:"hostname" json:"hostname" validate:"required"`
	Port     int    `yaml:"port" json:"port" validate:"required"`
	Token    string `yaml:"token" json:"token"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type Runner struct {
	ID          string `yaml:"id" json:"id" validate:"required"`
	Type        string `yaml:"type" json:"type" validate:"required,oneof=es5 risor"`
	ProgramPath string `yaml:"program_path" json:"program_path" validate:"omitempty,required_without=ProgramB64,excluded_with=ProgramB64"`
	ProgramB64  string `yaml:"program_b64" json:"program_b64" validate:"omitempty,required_without=ProgramPath,excluded_with=ProgramPath,base64"`
	Buffer      int    `yaml:"buffer" json:"buffer"`
}

type Line struct {
	ID       string `yaml:"id" json:"id" validate:"required"`
	InputID  string `yaml:"input_id" json:"input_id" validate:"required"`
	RunnerID string `yaml:"runner_id" json:"runner_id" validate:"required"`
	OutputID string `yaml:"output_id" json:"output_id" validate:"required"`
	CacheID  string `yaml:"cache_id" json:"cache_id"`
}

type Input struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Topic        string `yaml:"topic" json:"topic" validate:"required"`
	Method       string `yaml:"method" json:"method" validate:"omitempty,oneof=POST PUT PATCH"`
	Stream       string `yaml:"stream" json:"stream"`
	Client       string `yaml:"client" json:"client"`
	Buffer       int    `yaml:"buffer" json:"buffer"`
}

type Output struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Topic        string `yaml:"topic" json:"topic" validate:"required"`
	Method       string `yaml:"method" json:"method" validate:"omitempty,oneof=POST PUT PATCH"`
	Stream       string `yaml:"stream" json:"stream"`
	Client       string `yaml:"client" json:"client"`
}

type Cache struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Bucket       string `yaml:"bucket" json:"bucket" validate:"required"`
	Ttl          string `yaml:"ttl" json:"ttl" validate:"required"`
	Marshal      string `yaml:"marshal" json:"marshal" default:"msgpack" validate:"required,oneof=json msgpack gob"`
}
