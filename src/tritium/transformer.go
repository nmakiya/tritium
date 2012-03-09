package tritium

import (
	tp "athena/src/athena/proto"
)

type Engine interface {
	Run(transform *tp.Transform, input []byte, vars map[string]string) (output []byte, exports [][]string, logs []string)
}
