package runners

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/sandrolain/event-runner/src/config"
)

func GetProgramContent(c config.Runner) (program []byte, err error) {
	if c.ProgramPath != "" {
		program, err = os.ReadFile(c.ProgramPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read program file %s: %w", c.ProgramPath, err)
		}
		return
	}
	if c.ProgramB64 != "" {
		program, err = base64.StdEncoding.DecodeString(c.ProgramB64)
		if err != nil {
			return nil, fmt.Errorf("unable to decode program: %w", err)
		}
		return
	}
	err = fmt.Errorf("program not set")
	return
}
