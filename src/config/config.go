package config

type Config struct {
	Logger      Logger       `yaml:"logger" json:"logger" validate:"required"`
	Connections []Connection `yaml:"connections" json:"connections" validate:"required,dive"`
	Runners     []Runner     `yaml:"runners" json:"runners" validate:"required,dive"`
	Lines       []Line       `yaml:"lines" json:"lines" validate:"required,dive"`
	Inputs      []Input      `yaml:"inputs" json:"inputs" validate:"required,dive"`
	Outputs     []Output     `yaml:"outputs" json:"outputs" validate:"required,dive"`
	Caches      []Cache      `yaml:"caches" json:"caches" validate:"omitempty,dive"`
	Plugins     []Plugin     `yaml:"plugins" json:"plugins" validate:"omitempty,dive"`
}

type Logger struct {
	Level  string `yaml:"level" json:"level" default:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR debug info warn error"`
	Format string `yaml:"format" json:"format" default:"TEXT" validate:"required,oneof=TEXT JSON text json"`
	Color  bool   `yaml:"color" json:"color"`
}

type Connection struct {
	ID       string `yaml:"id" json:"id" validate:"required"`
	Type     string `yaml:"type" json:"type" validate:"required,oneof=nats redis kafka http grpc plugin"`
	PluginID string `yaml:"plugin_id" json:"plugin_id" validate:"omitempty,required_if=Type plugin"`
	Hostname string `yaml:"hostname" json:"hostname" validate:"omitempty,required,excluded_if=Type plugin"`
	Port     int    `yaml:"port" json:"port" validate:"omitempty,required,excluded_if=Type plugin"`
	Token    string `yaml:"token" json:"token"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type Runner struct {
	ID          string `yaml:"id" json:"id" validate:"required"`
	Type        string `yaml:"type" json:"type" validate:"required,oneof=es5 wasm"`
	ProgramPath string `yaml:"program_path" json:"program_path" validate:"omitempty,required_without=ProgramB64,excluded_with=ProgramB64"`
	ProgramB64  string `yaml:"program_b64" json:"program_b64" validate:"omitempty,required_without=ProgramPath,excluded_with=ProgramPath,base64"`
	Buffer      int    `yaml:"buffer" json:"buffer" default:"128" validate:"required"`
	Timeout     string `yaml:"timeout" json:"timeout" default:"5s" validate:"required,duration"`
	MaxStack    int    `yaml:"max_stack" json:"max_stack" default:"1024" validate:"required"`
	Concurrency int    `yaml:"concurrency" json:"concurrency" default:"1" validate:"required,gt=0"`
}

type Line struct {
	ID        string   `yaml:"id" json:"id" validate:"required"`
	InputID   string   `yaml:"input_id" json:"input_id" validate:"required"`
	RunnerID  string   `yaml:"runner_id" json:"runner_id" validate:"required"`
	OutputID  string   `yaml:"output_id" json:"output_id" validate:"required"`
	CacheID   string   `yaml:"cache_id" json:"cache_id"`
	PluginIDs []string `yaml:"plugin_ids" json:"plugin_ids"`
}

type Input struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Buffer       int    `yaml:"buffer" json:"buffer" default:"128" validate:"required"`
	Topic        string `yaml:"topic" json:"topic" validate:"required"`
	Method       string `yaml:"method" json:"method" validate:"omitempty,oneof=POST PUT PATCH"` // HTTP
	Stream       string `yaml:"stream" json:"stream"`                                           // NATS
	Client       string `yaml:"client" json:"client"`                                           // NATS
}

type Output struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Topic        string `yaml:"topic" json:"topic" validate:"required"`
	Method       string `yaml:"method" json:"method" validate:"omitempty,oneof=POST PUT PATCH"` // HTTP
	Hostname     string `yaml:"hostname" json:"hostname" validate:"omitempty"`                  // gRPC
	Port         int    `yaml:"port" json:"port" validate:"omitempty,gte=0,lte=65535"`          // gRPC
	Stream       string `yaml:"stream" json:"stream" validate:"omitempty"`                      // NATS
	Client       string `yaml:"client" json:"client" validate:"omitempty"`                      // NATS
}

type Cache struct {
	ID           string `yaml:"id" json:"id" validate:"required"`
	ConnectionID string `yaml:"connection_id" json:"connection_id" validate:"required"`
	Bucket       string `yaml:"bucket" json:"bucket" validate:"required"`
	Ttl          string `yaml:"ttl" json:"ttl" validate:"omitempty,duration"`
	Marshal      string `yaml:"marshal" json:"marshal" default:"msgpack" validate:"required,oneof=json msgpack gob"`
}

type Plugin struct {
	ID      string   `yaml:"id" json:"id" validate:"required"`
	Exec    string   `yaml:"exec" json:"exec" validate:"required"`
	Args    []string `yaml:"args" json:"args" validate:"omitempty"`
	Env     []string `yaml:"env" json:"env" validate:"omitempty"`
	Delay   string   `yaml:"delay" json:"delay" default:"1s" validate:"omitempty,duration"`
	Retry   int      `yaml:"retry" json:"retry" default:"3" validate:"omitempty,gt=0"`
	Marshal string   `yaml:"marshal" json:"marshal" default:"msgpack" validate:"required,oneof=json msgpack gob"`
}
