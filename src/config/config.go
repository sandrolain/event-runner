package config

type Config struct {
	Connections []Connection
	Runners     []Runner
	Inputs      []Input
	Outputs     []Output
}

type Line struct {
	Connection string `json:"connection" validate:"required"`
	Input      string `json:"input" validate:"required"`
	Runner     string `json:"runner" validate:"required"`
	Output     string `json:"output" validate:"required"`
}

type Connection struct {
	ID       string `json:"id" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=nats redis http"`
	URL      string `json:"url" validate:"required,url"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Runner struct {
	ID         string `json:"id" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=es5 risor"`
	ScriptPath string `json:"script_path" validate:"required,file"`
	ScriptB64  string `json:"script_b64" validate:"required,base64"`
}

type Input struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Stream string `json:"stream"`
	Client string `json:"client"`
}

type Output struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Stream string `json:"stream"`
	Client string `json:"client"`
}
